package main

// 爬楼梯
func climbStairs(n int) int {
	if n == 0 || n == 1 || n == 2 {
		return n
	}
	pq := make([]int, n+1)
	pq[1] = 1
	pq[2] = 2
	for i := 3; i <= n; i++ {
		pq[i] = pq[i-1] + pq[i-2]
	}
	return pq[n]
}

// 最长连续递增子序列长度
func lengthOfLIS(nums []int) int {
	if len(nums) == 0 || len(nums) == 1 {
		return len(nums)
	}
	// 初始化 dp 数组，dp[i] 表示以 nums[i] 结尾的最长连续递增子序列的长度
	dp := make([]int, len(nums))
	// 每个元素自身可以构成一个长度为 1 的子序列
	for i := range dp {
		dp[i] = 1
	}
	// 记录最长连续递增子序列的长度
	maxLength := 1
	// 遍历数组
	for i := 1; i < len(nums); i++ {
		if nums[i] > nums[i-1] {
			// 如果当前元素大于前一个元素，则更新 dp[i] 的值
			dp[i] = dp[i-1] + 1
		}
		// 更新最长连续递增子序列的长度
		if dp[i] > maxLength {
			maxLength = dp[i]
		}
	}
	return maxLength
}
