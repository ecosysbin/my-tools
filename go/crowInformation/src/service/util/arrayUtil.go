package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello world")
	// ran := [10]int{3,1,5,6,8,9,10,23,43,21}的[]里面有了长度则定义的是数组，这时{}里面只能传10个int数字；ran := []int{3,1,5,6,8,9,10,23,43,21}的[]里面没有长度则定义的是切片（动态数组），这时{}里面可任意个数数字
	// 函数的参数是数组类型只能传数组（数组长度是数组的属性，需要保持一致），函数参数是切片时只能传切片，不能传数组
	ran := []int{3,1,5,6,8,9,10,23,43,21}
	//sortRange(ran)
	//fmt.Println(ran)
	//ran = []int{2,1,5,6,8,9,10,23,43,21}
	//maopaoSort(ran)
	// ran = []int{2,1,5,6,8,9,10,23,43,21}
	quickSort(ran,0,9)
	fmt.Println(ran)
}

// 选择排序
func sortRange(ran []int) []int{
	for i := 0; i < len(ran) - 1; i ++ {
		for j := i + 1;j < len(ran); j++ {
			if ran[j] < ran[i] {
				tmp := ran[i]
				ran[i] = ran[j]
				ran[j] = tmp
			}
		}
	}
	return  ran
}

// 冒泡排序
func maopaoSort(ran []int) []int{
	for i := 0; i < len(ran); i++ {
		for j := 0; j < len(ran) - i - 1; j++ {
			if ran[j] > ran[j +1] {
				tmp := ran[j]
				ran[j] = ran[j+1]
				ran[j+1] = tmp
			}
		}
	}
	return ran
}

// 桶排序
func tongSort(ran []int) []int{
	// 1. 先找到最大的数
	max := ran[0]
	for i := 0; i < len(ran); i++ {
		if max < ran[i] {
			max = ran[i]
		}
	}
	fmt.Println(max)
	tongRange := []int{}
	for i := 0; i < max; i++{
		tongRange[i] = 0
	}

	for i := 0; i < len(ran); i++ {
		tongRange[ran[i]]++
	}

	fmt.Println(tongRange)

	return ran
}

func quickSort(arr []int, left, right int) {
	if left > right {
		return
	}
	fmt.Println(arr)
	i, j := left, right
	tmp := arr[i]
	for ( i != j) {
		// 先找右边小于tmp的数
		for (arr[j] >= tmp && i < j ){
			j--
		}
		// 再从左边找大于tmp的数
		for (arr[i] <= tmp && i < j){
			i++
		}
		// 交换数组中两个元素位置
		arr[i],arr[j]=arr[j],arr[i]
	}
	// 准基数归位
	arr[left] = arr[i]
	arr[i] = tmp
	// 递归左边
	quickSort(arr, left, i-1)
	// 递归右边
	quickSort(arr, i+1, right)
}
