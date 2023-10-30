# 令牌桶限流器（GoLimiter）

## 简介
`golimiter` 包提供了一个基于令牌桶算法实现的限流器。

## 使用方法
```go
package main

import (
	"fmt"
	"time"

	"github.com/your/package/golimiter"
)

func main() {
	// 创建一个新的限流器实例，限制请求数为 10 次/秒，桶大小为 100
	limiter := golimiter.NewGoLimiter(10, 100)

	// 模拟一些请求
	for i := 1; i <= 100; i++ {
		// 假设每个请求需要一定时间处理
		time.Sleep(time.Millisecond * 200)

		// 检查是否允许执行请求
		if limiter.Allow(1) {
			// 执行请求
			fmt.Println("执行请求", i)
		} else {
			// 请求被限流
			fmt.Println("请求被限流", i)
		}
	}

	// 停止限流器
	limiter.StopLimiter()
}