package golimiter

import (
	"context"
	"sync"
	"time"
)

type GoLimiter struct {
	rate       int        // 请求评率限制
	burst      int        // 桶大小（允许突发的请求数量）
	mtx        sync.Mutex // 互斥锁
	currTokens int        // 当前可用的令牌数量
	lastUpdate time.Time  // 上次更新令牌数量的时间
	ctx        context.Context
	cancel     context.CancelFunc
}

// rate：请求评率限制，单位为每秒请求数
// burst：桶大小
// 返回一个新的限流器实例指针
func NewGoLimiter(rate, burst int) *GoLimiter {
	l := &GoLimiter{
		rate:       rate,
		burst:      burst,
		currTokens: burst,
		lastUpdate: time.Now(),
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
			l.refill()
		}
	}
}

func (l *GoLimiter) refill() {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	elapsed := time.Since(l.lastUpdate)
	tokensToAdd := int(float64(elapsed.Nanoseconds()) / float64(time.Second.Nanoseconds()) * float64(l.rate))

	if l.currTokens+tokensToAdd > l.burst {
		l.currTokens = l.burst
	} else {
		l.currTokens += tokensToAdd
	}

	l.lastUpdate = time.Now()
}

func (l *GoLimiter) Allow(curr int) bool {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.currTokens >= curr {
		l.currTokens -= curr
		return true
	}

	return false
}

func (l *GoLimiter) StopLimiter() {
	l.cancel()
}
