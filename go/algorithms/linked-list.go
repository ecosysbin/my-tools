package main

import "fmt"

type Node struct {
	Val  int
	Next *Node
}

type ListNode struct {
	Val  int
	Next *ListNode
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

func getIntersectionNode(headA, headB *Node) *Node {
	if headA == nil || headB == nil {
		return nil
	}
	// 计算两个链表的长度，找出长链表
	lengthA, lengthB := 0, 0
	tmpHeadA, tmpHeadB := headA, headB
	for tmpHeadA != nil {
		lengthA++
		tmpHeadA = tmpHeadA.Next
	}

	for tmpHeadB != nil {
		lengthB++
		tmpHeadB = tmpHeadB.Next
	}

	if lengthA > lengthB {
		step := lengthA - lengthB
		for ; step > 0; step-- {
			headA = headA.Next
		}
	} else {
		step := lengthB - lengthA
		for ; step > 0; step-- {
			headB = headB.Next
		}
	}

	// 一样长了
	for headA != nil && headB != nil {
		if headA == headB {
			return headA
		}
		headA = headA.Next
		headB = headB.Next
	}
	return nil
}

// 回文链表
func isPalindrome(head *Node) bool {
	if head == nil {
		return false
	}
	// 1. 先找到中间节点(快慢指针找中间节点更快一些，时间是n/2)
	// 1.1 计算链表长度
	length := 0
	tmpHead := head
	for tmpHead != nil {
		length++
		tmpHead = tmpHead.Next
	}
	steps := length / 2
	tmpHead1 := head
	for steps > 0 {
		tmpHead1 = tmpHead1.Next
		steps--
	}
	// 2. 从中间节点开始对后文的链表进行翻转
	newHead := reverseNode(tmpHead1)
	// 3. 从前文，后文链表开始遍历，看看节点是否相同
	for head != nil && newHead != nil {
		if head.Val != newHead.Val {
			return false
		}
		head = head.Next
		newHead = newHead.Next
	}
	return true
}

// 判断链表是否有环
func hasCycle(head *ListNode) bool {
	if head == nil {
		return false
	}
	slow, fast := head, head.Next
	for slow != nil && fast != nil && fast.Next != nil {
		if slow == fast {
			return true
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	return false
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	var preNewHead = &ListNode{
		Val:  -1,
		Next: nil,
	}
	tmpPreNewHead := preNewHead
	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			tmpPreNewHead.Next = list1
			list1 = list1.Next
		} else {
			tmpPreNewHead.Next = list2
			list2 = list2.Next
		}
		tmpPreNewHead = tmpPreNewHead.Next
	}

	// 拼接剩余的链表
	if list1 != nil {
		tmpPreNewHead.Next = list1
	} else {
		tmpPreNewHead.Next = list2
	}
	return preNewHead.Next
}
