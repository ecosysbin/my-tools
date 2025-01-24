package map_test

import (
	"fmt"
)

func PrintPerson() {
	var p = &Person{
		param: map[string]string{
			"name": "张三",
			"age":  "25",
		},
	}
	p.GetParam()["name"] = "李四"
	fmt.Println(p.GetParam())
	tt := map[int]int{
		1:   1,
		2:   2,
		5:   5,
		10:  10,
		11:  11,
		100: 100,
	}
	for k, v := range tt {
		fmt.Printf("k %d v %d \n", k, v)
	}
	var bb map[string]string
	var cc = map[string]string{}
	fmt.Printf("len bb %d, len cc %d", len(bb), len(cc))
}

type Person struct {
	param map[string]string
}

func (p *Person) GetParam() map[string]string {
	return p.param
}
