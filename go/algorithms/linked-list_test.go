package main

import (
	"fmt"
	"testing"
)

func TestReverse(t *testing.T) {
	head := &Node{Val: 1, Next: &Node{Val: 2, Next: &Node{Val: 3, Next: &Node{Val: 4, Next: &Node{Val: 5}}}}}
	fmt.Println("Original List:")
	print(head)
	newHead := reverse(head)
	fmt.Println("Reversed List:")
	print(newHead)
}

func TestCombine(t *testing.T) {
	l1 := &Node{Val: 1, Next: &Node{Val: 3, Next: &Node{Val: 5, Next: &Node{Val: 7, Next: &Node{Val: 9}}}}}
	l2 := &Node{Val: 2, Next: &Node{Val: 4, Next: &Node{Val: 8, Next: &Node{Val: 9}}}}
	newHead := combine(l1, l2)
	print(newHead)
}

func TestMiddleNode(t *testing.T) {
	head := &Node{Val: 1, Next: &Node{Val: 3, Next: &Node{Val: 5, Next: &Node{Val: 7, Next: &Node{Val: 9}}}}}
	middle := middleNode(head)
	fmt.Println("Middle node:", middle.Val)
}
