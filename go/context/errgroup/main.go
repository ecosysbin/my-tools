package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
	}
	for i := range urls {
		url := urls[i]
		g.Go(func() error {
			fmt.Println("Fetching", url)
			// resp, err := http.Get(url)
			// if err == nil {
			// 	resp.Body.Close()
			// }
			time.Sleep(3 * time.Second)
			return fmt.Errorf("failed to fetch %s", url)
		})
	}
	// 好处是能帮忙把go协程的err记录下来，只会记录第一个err，其他的协程不会阻塞，发生err也不会记录
	err := g.Wait()
	if err == nil {
		fmt.Println("Successfully fetched all URLs.")
	} else {
		fmt.Println("Failed to fetch some URLs:", err)
	}
}
