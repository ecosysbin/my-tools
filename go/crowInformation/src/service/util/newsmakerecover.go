package main

import (
	"fmt"
)

func main() {

	fmt.Println(panictest(3))
	fmt.Println(panictest(0))
	fmt.Println("finish test")
	newtest()
}

type NT struct {
	name string
	age string
}

func panictest(i int) int {
	// 可以扑获
	defer func() {
		recover()
	}()
	// 不能扑获
	// defer recover()
	a := 1/i
	return a
}

func newtest() {
	a := make(chan int, 1)
	fmt.Println(a)
	n := new(NT)
	// t := new(NT{"name": "BILL", "age": "12"})  wrong
	// t := NT{name: "BILL", age: "12"} ok
	// new返回的是lingzhi ,不会做初始化，如下：
	t := new(chan int)
	i := new([]int)
	fmt.Println(n)
	fmt.Println(t)
	fmt.Println(i)
}

//1、与C语言不同，T{}分配的局部变量是可以返回的，且返回后该空间不会释放，例如
// 这样既能反映go是面向对象，而c是面向过程
//type T struct {
//	i, j int
//}
//func a(i, j int) T {
//	i := T { i, j}
//	return i
//}
//func b {
//	t = a(1, 2)
//	fmt.Println(t)
//}

//new 可以创建任意类型，并返回对象指针，但不会做初始化
//make 只能创建slice、map、channel，并做初始化，返回的是对象的引用
//结构体对象的创建还可以使用字面量T{}, 可在{}中初始化结构体
