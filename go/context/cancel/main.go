package main

import (
	"context"
	"fmt"
	"time"
)

type myContextKey string

func main() {
	// 下面是context传参，cancel比较简单，不再举例
	ctx := context.Background()
	ctx = context.WithValue(ctx, myContextKey("name"), "tony")
	ctx = context.WithValue(ctx, myContextKey("age"), "21")
	go func(ctx context.Context) {
		fmt.Println(ctx.Value(myContextKey("name")))
		fmt.Println(ctx.Value(myContextKey("age")))
		println("context end")
	}(ctx)
	time.Sleep(1 * time.Second)
	fmt.Println("main end")
}
