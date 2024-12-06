package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	println("Hello, world!")
	// context.Background和context.TODO都是返回一个空的上下文（都是包装的emptyCtx，emptyCtx实现context接口返回的都是nil）
	// context.Background 是上下文的默认值，所有其他的上下文都应该从它衍生出来；
	// context.TODO 应该仅在不确定应该使用哪种上下文时使用；
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go handle(ctx, 5*time.Second)
	// 1. ctx超时后会触发ctx.Done()返回的channel，因此各协程里通过ctx.Done()都会接收到通知（注意主协程打印退出后，子协程则会还来不及打印，主协程阻塞时，都会有打印）。
	// 协程要是先走了time.After()的分支则go协程会退出，不会再有ctx.Done()分支的通知
	// ctx.Done()接收的信号后，可以通过ctx.Err()获取原因
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}
