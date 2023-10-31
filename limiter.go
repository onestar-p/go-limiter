package golimiter

import (
	"context"
	"sync"
	"time"
)

type GoLimiter struct {
	rate        int                // 请求评率限制
	burst       int                // 桶大小（允许突发的请求数量）
	mtx         sync.Mutex         // 互斥锁
	currTokens  int                // 当前可用的令牌数量
	lastUpdated time.Time          // 上次更新令牌数量的时间
	ctx         context.Context    // 上下文
	cancel      context.CancelFunc // 取消函数
	refillCh    chan struct{}      // 补充令牌的通道
}

func NewGoLimiter(rate, burst int) *GoLimiter {
	l := &GoLimiter{
		rate:        rate,
		burst:       burst,
		currTokens:  burst,
		lastUpdated: time.Now(),
		refillCh:    make(chan struct{}),
	}

	l.ctx, l.cancel = context.WithCancel(context.Background())
	go l.run()

	return l
}

func (l *GoLimiter) run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			select {
			case l.refillCh <- struct{}{}:
			default:
			}
		}
	}
}

func (l *GoLimiter) refill() {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	elapsed := time.Since(l.lastUpdated)
	tokensToAdd := int(float64(elapsed.Nanoseconds()) / float64(time.Second.Nanoseconds()) * float64(l.rate))

	if l.currTokens+tokensToAdd > l.burst {
		l.currTokens = l.burst
	} else {
		l.currTokens += tokensToAdd
	}

	l.lastUpdated = time.Now()
}

func (l *GoLimiter) Allow(curr int) bool {
	if curr <= 0 {
		return false
	}

	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.currTokens >= curr {
		l.currTokens -= curr
		return true
	}

	return false
}

func (l *GoLimiter) StartRefillLoop() {
	go l.refillLoop()
}

func (l *GoLimiter) refillLoop() {
	for {
		select {
		case <-l.refillCh:
			l.refill()
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *GoLimiter) StopLimiter() {
	l.cancel()
	close(l.refillCh)
}
