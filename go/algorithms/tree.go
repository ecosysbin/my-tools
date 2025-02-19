package main

import (
	"fmt"
	"math"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// 树的按层遍历
func cengPrint(head *TreeNode) {
	printNode := []*TreeNode{}
	queue := []*TreeNode{head}

	for len(queue) > 0 {
		tmpQueue := []*TreeNode{}
		for _, node := range queue {
			printNode = append(printNode, node)
			if node.Left != nil {
				tmpQueue = append(tmpQueue, node.Left)
			}
			if node.Right != nil {
				tmpQueue = append(tmpQueue, node.Right)
			}
		}
		queue = tmpQueue
	}

	for _, node := range printNode {
		fmt.Printf("%d\n", node.Val)
	}
}

func maxDepth(root *TreeNode) int {
	return 1
}

// 前序：先根，再左，后右
// 中序：先左，再根，后右
// 后序：先左，再右，后根

//      A
//    B    C
//   D E     F

// 前序：A B D E C F
// 中序：D B E A C F
// 后序：D E B F C A

func qianXuBianli(root *TreeNode) []int {
	var result []int
	var qianxu func(node *TreeNode)
	qianxu = func(node *TreeNode) {
		if node == nil {
			return
		}
		result = append(result, node.Val)
		qianxu(node.Left)
		qianxu(node.Right)
	}
	qianxu(root)
	return result
}

func zhongxunBianli(root *TreeNode) []int {
	var result []int
	var zhongxu func(node *TreeNode)
	zhongxu = func(node *TreeNode) {
		if node == nil {
			return
		}
		zhongxu(node.Left)
		result = append(result, node.Val)
		zhongxu(node.Right)
	}
	zhongxu(root)
	return result
}

func houxuBianli(root *TreeNode) []int {
	var result []int
	var houxu func(node *TreeNode)
	houxu = func(node *TreeNode) {
		if node == nil {
			return
		}
		houxu(node.Left)
		houxu(node.Right)
		result = append(result, node.Val)
	}
	houxu(root)
	return result
}

// 翻转二叉树（递归）
func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	invertLeft := invertTree(root.Left)
	invertRight := invertTree(root.Right)
	root.Left = invertRight
	root.Right = invertLeft
	return root
}

// 判断二叉树是否对称
func isSymmetric(root *TreeNode) bool {
	if root == nil {
		return true
	}
	return isDuicheng(root.Left, root.Right)
}

func isDuicheng(leftNode, rightNode *TreeNode) bool {
	if leftNode == nil && rightNode == nil {
		return true
	}
	if leftNode == nil || rightNode == nil || leftNode.Val != rightNode.Val {
		return false
	}
	return isDuicheng(leftNode.Left, rightNode.Right) && isDuicheng(leftNode.Right, rightNode.Left)
}

// 求最大直径
var max int

func qianxu(root *TreeNode) {
	if root == nil {
	}
	if (depth(root.Left) + depth(root.Right)) > max {
		max = depth(root.Left) + depth(root.Right)
	}
	if root.Left != nil {
		qianxu(root.Left)
	}
	if root.Right != nil {
		qianxu(root.Right)
	}
}

func depth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return int(math.Max(float64(depth(root.Left)), float64(depth(root.Right)))) + 1
}

// 验证二叉树搜索树
func isValidBST(root *TreeNode) bool {
	// 二叉树搜索树，左，中间，右依次变大。每一个叶子节点也是二叉树搜索树
	// 先中序遍历，然后是不是一个递增序列
	result := zhongxunBianli(root)
	for i := 0; i < len(result)-1; i++ {
		if result[i] >= result[i+1] {
			return false
		}
	}
	return true
}
