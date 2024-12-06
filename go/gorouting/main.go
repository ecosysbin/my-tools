package main

import (
	"fmt"
	"sync"
)

func main() {
	// 1. channel管道必须有客户端开始消费才能生产 2. 客户端退出（即不在消费，如go协程退出）channel管道则不能再生产，1.  2. 都会触发死锁报错
	var wg sync.WaitGroup
	chOdd := make(chan struct{})
	chEven := make(chan struct{})
	// 从debug看 channel对象确实有很多属性包括消费者，生产者，以及锁等
	// qcount — Channel 中的元素个数；
	// dataqsiz — Channel 中的循环队列的长度；
	// buf — Channel 的缓冲区数据指针；
	// sendx — Channel 的发送操作处理到的位置；
	// recvx — Channel 的接收操作处理到的位置；
	wg.Add(2)
	
	go printOdd(&wg, chOdd, chEven)
	go printEven(&wg, chOdd, chEven)
	// 必须有消费者了，才能开始生产
	chOdd <- struct{}{}
	wg.Wait()
}

func printOdd(wg *sync.WaitGroup, chOdd, chEven chan struct{}) {
	defer wg.Done()
	for i := 1; i <= 100; i += 2 {
		<-chOdd
		fmt.Println(i)
		chEven <- struct{}{}
	}
}

func printEven(wg *sync.WaitGroup, chOdd, chEven chan struct{}) {
	defer wg.Done()
	for i := 2; i <= 100; i += 2 {
		<-chEven
		fmt.Println(i)
		if i < 100 {
			chOdd <- struct{}{}
		}
	}
}
