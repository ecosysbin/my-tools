package main

// 全排列
func permute(nums []int) [][]int {
	if len(nums) == 0 || len(nums) == 1 {
		return [][]int{nums}
	}

	if len(nums) == 2 {
		return [][]int{[]int{nums[0], nums[1]}, []int{nums[1], nums[0]}}
	}
	r := [][]int{}
	for i := 0; i < len(nums); i++ {
		input := []int{}
		input = append(input, nums[:i]...)
		input = append(input, nums[i+1:]...)
		arr := permute(input)
		for j := range arr {
			arr[j] = append([]int{nums[i]}, arr[j]...)
		}
		r = append(r, arr...)
	}
	return r
}

// 子集
func subsets(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{nums}
	}

	if len(nums) == 1 {
		return [][]int{[]int{}, nums}
	}

	if len(nums) == 2 {
		return [][]int{[]int{}, nums, []int{nums[0]}, []int{nums[1]}}
	}

	r := [][]int{}
	first := nums[0]
	r = append(r, []int{first})
	input := nums[1:]
	arr := subsets(input)
	r = append(r, arr...)
	for j := range arr {
		if len(arr[j]) > 0 {
			tmpJ := append([]int{first}, arr[j]...)
			r = append(r, tmpJ)
		}
	}
	return r
}
