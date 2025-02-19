package main

import (
	"fmt"
	"testing"
)

func TestSubsets(t *testing.T) {
	input := []int{1, 2, 3, 4}
	result := subsets(input)
	fmt.Printf("result: %v", result)
}
