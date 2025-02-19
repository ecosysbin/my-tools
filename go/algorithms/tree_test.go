package main

import (
	"fmt"
	"testing"
)

//      1
//    2    3
//   4 5     6

// 前序：A B D E C F（1 2 4 5 3 6）
// 中序：D B E A C F（4 2 5 1 3 6）
// 后序：D E B F C A（4 5 2 6 3 1）

func TestPrintTree(t *testing.T) {
	tree := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 4,
			},
			Right: &TreeNode{
				Val: 5,
			},
		},
		Right: &TreeNode{
			Val: 3,
			Right: &TreeNode{
				Val: 6,
			},
		},
	}
	// cengPrint(tree)
	qianxu := qianXuBianli(tree)
	fmt.Printf("qianxu: %v \n", qianxu)

	zhongxu := zhongxunBianli(tree)
	fmt.Printf("zhongxu: %v \n", zhongxu)

	houxu := houxuBianli(tree)
	fmt.Printf("houxu: %v \n", houxu)
}

func TestMaxSubArray(t *testing.T) {
	nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	max := maxSubArray(nums)
	fmt.Println(max)
}

func TestRotate(t *testing.T) {
	nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	rotate(nums, 3)
	fmt.Println(nums)
}
