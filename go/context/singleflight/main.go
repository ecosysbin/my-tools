package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {
	fmt.Println("Hello, world!")
	time.Sleep(1 * time.Second)
}

type service struct {
	requestGroup singleflight.Group
}

// requestGroup.Do则可以保证同一个请求执行一次，不会将数据库打爆
// func (s *service) handleRequest(ctx context.Context, request Request) (Response, error) {
// 	v, err, _ := requestGroup.Do(request.Hash(), func() (interface{}, error) {
// 		rows, err := // select * from tables
// 		if err != nil {
// 		    return nil, err
// 		}
// 		return rows, nil
// 		return nil, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return Response{
// 		rows: rows,
// 	}, nil
// }
