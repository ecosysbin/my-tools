package main

import (
	"math"
)

func twoSum(nums []int, target int) []int {
	result := []int{}
	// hash表实现
	numMap := map[int][]int{}
	for index, n := range nums {
		if numMap[n] == nil {
			numMap[n] = []int{}
		}
		numMap[n] = append(numMap[n], index)
	}

	for mIndex, m := range nums {
		indexs, ok := numMap[target-m]
		if !ok {
			continue
		}
		if m == target-m && len(indexs) >= 2 {
			result = append(result, indexs[0])
			result = append(result, indexs[1])
			return result
		}
		if m != target-m {
			result = append(result, mIndex)
			result = append(result, indexs[0])
			return result
		}
	}
	return result
}

func longestConsecutive(nums []int) int {
	longest := 0
	numsMap := map[int]bool{}
	for _, n := range nums {
		numsMap[n] = true
	}
	for _, n := range nums {
		longestN := 1
		if _, ok := numsMap[n+1]; !ok {
			for i := n - 1; ; i-- {
				if _, ok := numsMap[i]; !ok {
					break
				}
				longestN++
			}
		}
		if longestN > longest {
			longest = longestN
		}
	}
	return longest
}

func moveZeroes(nums []int) {
	for i := 0; i < len(nums); i++ {
		if nums[i] == 0 {
			find := false
			for j := i + 1; j < len(nums); j++ {
				if nums[j] != 0 {
					nums[i], nums[j] = nums[j], nums[i]
					find = true
					break
				}
			}
			if !find {
				break
			}
		}
	}
}

func maxArea(height []int) int {
	max := 0
	if len(height) < 2 {
		return 0
	}
	left := 0
	right := len(height) - 1
	length := len(height) - 1
	high := math.Min(float64(height[left]), float64(height[right]))
	max = length * int(high)
	for left < right {
		high = math.Min(float64(height[left]), float64(height[right]))
		length = right - left
		if length*int(high) > max {
			max = length * int(high)
		}
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}
	return max
}

func lengthOfLongestSubstring(s string) int {
	if len(s) == 0 {
		return 0
	}
	if len(s) == 1 {
		return 1
	}
	length := 1
	for i := 0; i < len(s)-1; i++ {
		tmpMap := map[byte]bool{}
		tmpMap[s[i]] = true
		if s[i+1] != s[i] {
			j := i + 1
			for ; j < len(s); j++ {
				tmpMap[s[j]] = true
				if len(tmpMap) != j-i+1 {
					break
				}
			}
			i = j
			if len(tmpMap) > length {
				length = len(tmpMap)
			}
		}
	}
	return length
}

func subarraySum(nums []int, k int) int {
	if len(nums) == 0 {
		return 0
	}
	times := 0
	for i := 0; i < len(nums); i++ {
		// if nums[i] == k {
		// 	times++
		// }
		sumI := 0
		for j := i; j < len(nums); j++ {
			sumI += nums[j]
			if sumI == k {
				times++
			}
		}
	}
	return times
}

func reverseNode(head *Node) *Node {
	var pre, current, next *Node
	pre = nil
	current = head

	for current != nil {
		next = current.Next
		current.Next = pre
		pre = current
		current = next
	}
	return pre
}

func maxSubArray(nums []int) int {
	if len(nums) == 0 {
		return -1
	}
	if len(nums) == 1 {
		return nums[0]
	}
	// 全小于等于0的情况
	max := math.MinInt
	for i := 0; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}
	if max <= 0 {
		return max
	}
	maxI := 0
	for i := 0; i < len(nums); i++ {
		maxI += nums[i]
		if maxI <= 0 {
			maxI = 0
			continue
		}
		if maxI > max {
			max = maxI
		}
	}
	return max
}

// 轮转数组(关键是原地)
func rotate(nums []int, k int) {
	// k和nums取余，避免多次翻转
	k %= len(nums)
	// 1. 先翻转整个链
	reverseNums(nums)
	// 2. 再翻转前K个
	reverseNums(nums[:k])
	// 3. 再翻转后k个
	reverseNums(nums[k:])
}

// 翻转数组
func reverseNums(nums []int) {
	for i := 0; i < len(nums)/2; i++ {
		nums[i], nums[len(nums)-1-i] = nums[len(nums)-1-i], nums[i]
	}
}

func productExceptSelf(nums []int) []int {
	if len(nums) == 0 {
		return nil
	}
	left, anwser, right := make([]int, len(nums)), make([]int, len(nums)), make([]int, len(nums))
	left[0] = 1
	for i := 1; i < len(nums); i++ {
		left[i] = left[i-1] * nums[i-1]
	}

	right[len(nums)-1] = 1
	for i := len(nums) - 2; i >= 0; i-- {
		right[i] = right[i+1] * nums[i+1]
	}

	for i := 0; i < len(nums); i++ {
		anwser[i] = left[i] * right[i]
	}
	return anwser
}
