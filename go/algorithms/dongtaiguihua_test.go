package main

import (
	"fmt"
	"testing"
)

func TestLengthOfLIS(t *testing.T) {
	maxLength := lengthOfLIS([]int{10, 7, 8, 9})
	fmt.Println(maxLength)
}
