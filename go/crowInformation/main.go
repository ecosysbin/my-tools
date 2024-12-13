// crowInformation project main.go
package main

import (
	_ "service/routers"
	_ "service/model"
	_ "service/dao"

	"github.com/astaxie/beego"
	"fmt"
)

func main() {
	fmt.Println("Hello crow news!")
	beego.Run()
}
