package main

import "fmt"

type Node struct {
	Val  int
	Next *Node
}

func print(head *Node) {
	cur := head
	for cur != nil {
		fmt.Printf("%d\n", cur.Val)
		cur = cur.Next
	}
}

func reverse(head *Node) *Node {
	var pre, cur, next *Node
	pre = nil
	cur = head
	for cur != nil {
		next = cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}
	return pre
}

func combine(l1, l2 *Node) *Node {
	// 哨兵
	newHead := &Node{Val: 0}
	cur := newHead
	for l1 != nil && l2 != nil {
		if l1.Val <= l2.Val {
			cur.Next = l1
			l1 = l1.Next
		} else {
			cur.Next = l2
			l2 = l2.Next
		}
		cur = cur.Next
	}
	if l1 != nil {
		cur.Next = l1
	}
	if l2 != nil {
		cur.Next = l2
	}
	return newHead.Next
}

func middleNode(head *Node) *Node {
	var slow, fast *Node
	slow, fast = head, head
	for slow != nil && fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}
