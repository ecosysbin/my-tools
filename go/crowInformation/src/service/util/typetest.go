// crowInformation project main.go
package main

import (
	"fmt"
	"reflect"
	"os"
)

func main() {

	var s = Do(add, 2,3)
	var b = Do(sub, 2,3)
	fmt.Println(s)
	fmt.Println(b)
	os.Exit(0)
}
func testByteString() {
	var b = []byte{104, 101, 108}
	// b = append(b, 104, 101, 108)
	fmt.Println(getStringfrombyte(b))
	// fmt.Println(getbyteFromstring("hello"))
}

func getRandInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100)
}

func getStringfrombyte(b []byte) string {
	return string(b)
}

func getbyteFromstring(str string) []byte {
	return []byte(str)
}

// 定义一个函数类型
type Op func(int64, int64) int64

func Do(f Op, a, b int64) int64{
	return f(a, b)
}

func add(a, b int64) int64{
	return a + b
}

func sub(a, b int64) int64{
	return a - b
}

func slimapTest() {
	// 1. 字面量方式创建map
	ma := map[string]string{}
	ma["name"] = "wangbin"
	ma["address"] = "bj"
	// 2. make 函数方式创建map
	mm := make(map[string]string)
	mm["happy"] = "food"
	mm["love"] = "friend\""
	fmt.Println(ma)
	fmt.Println(mm)
}

func stringTest() {
	st := "abcdf"
	// str[i] 是byte 类型 // 字符串是底层是一个二元数据类型，一个指向字节数组的起点，另一个是长度
	fmt.Println(st[1])
	// 转化为string
	fmt.Println(string(st[1]))
	fmt.Println(st[1:3])
	fmt.Println(st)
	ru := []rune(st)
	bt := []byte(st)
	fmt.Println(ru)
	fmt.Println(bt)
}

func indexTest() {
	// index
	in := [10]int64{}
	for i := 0; i < 12; i++ {
		// append方法会报错
		in[i] = int64(i)
	}
	// slice
	sl := make([]int,0)
	for i := 0; i< 12; i++{
		// sl[i] = i 会角标越界
		sl = append(sl, i)
	}
	isl := []int64{}
	isl = append(isl, 1)

}

func typeTest() {
	in := [10]int64{}
	sl := make([]int,0)
	fmt.Println("in", reflect.TypeOf(in))
	fmt.Println("sl", reflect.TypeOf(sl))
}

