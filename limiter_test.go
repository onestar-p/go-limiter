package golimiter_test

import (
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	golimiter "github.com/onestar-p/go-limiter"
)

var Limiter *golimiter.GoLimiter

func TestLimiter(t *testing.T) {

	var wg sync.WaitGroup
	var iList []int

	var mx sync.Mutex
	for i := 0; i < 112; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				if !Limiter.Allow(1) {
					continue
				}

				fmt.Printf("append %d\n", i)

				mx.Lock()
				iList = append(iList, i)
				mx.Unlock()
				break
			}
		}(i)

	}
	wg.Wait()

	sort.Slice(iList, func(i, j int) bool {
		return iList[i] < iList[j]
	})
	fmt.Println(iList)

	fmt.Println(len(iList))

}

func TestLimiter2(t *testing.T) {
	// 启动补充令牌循环
	Limiter.StartRefillLoop()
	// 模拟一些请求
	for i := 1; i <= 100; i++ {
		// 假设每个请求需要一定时间处理
		time.Sleep(time.Millisecond * 200)

		// 检查是否允许执行请求
		if Limiter.Allow(1) {
			// 执行请求
			fmt.Println("执行请求", i)
		} else {
			// 请求被限流
			fmt.Println("请求被限流", i)
		}
	}

	// 停止限流器
	Limiter.StopLimiter()
}

func TestMain(m *testing.M) {
	Limiter = golimiter.NewGoLimiter(5, 5)

	// 启动补充令牌循环
	Limiter.StartRefillLoop()

	m.Run()
}
