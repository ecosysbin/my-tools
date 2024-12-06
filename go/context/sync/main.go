package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

var status int64

func main() {
	// sync.Once 和sync.WaitGroup已经很熟悉了，这里就不再赘述了。
	// 这里主要说明下sync.NewCond
	fmt.Println("hello")
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		go listen(c, i)
	}
	time.Sleep(3 * time.Second)
	go broadcast(c)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func listen(c *sync.Cond, index int) {
	c.L.Lock()
	fmt.Println("listen lock")
	// 看着atomic.LoadInt64和直接使用status变量也没啥区别，但是使用atomic.LoadInt64可以保证线程安全
	for atomic.LoadInt64(&status) != 1 {
		// for status != 1 {
		// 会卡主，等待broadcast通知唤醒，这时候则释放了线程资源，相比于for轮转占用cpu资源，这种方式资源利用率更高
		c.Wait()
		fmt.Printf("listen wait, %d \n", index)
	}
	fmt.Printf("listen, %d \n", index)
	c.L.Unlock()
}

func broadcast(c *sync.Cond) {
	c.L.Lock()
	atomic.StoreInt64(&status, 1)
	// status = 1
	c.Broadcast()
	c.L.Unlock()
}
