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
}

type Person struct {
	param map[string]string
}

func (p *Person) GetParam() map[string]string {
	return p.param
}
