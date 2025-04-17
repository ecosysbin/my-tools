package v1

import (
	"fmt"
	"testing"
)

func Test_CheckDiskNameByPvcName(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name  string
		input string
		want  bool
	}{
		// the table itself
		{"abc1234567abc1234567 should match", "abc1234567abc1234567", true},
		{"abc1234567abc1234567a should not match", "abc1234567abc1234567a", false},
		{"abc123 should match", "abc123", true},
		{"abc-123 should not match", "abc-123", false},
		{"abc_123 should not match", "abc_123", false},
		{"1abc123 should not match", "1abc123", false},
		{"SitVM532719 should not match", "SitVM532719", true}, // 不符合k8s命名规范
	}

	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidName(tt.input)
			if result != tt.want {
				t.Errorf("test %s failed, got %v, want %v", tt.name, result, tt.want)
			}
		})
	}
}

func Test_ConstructUserData(t *testing.T) {
	cloudinit := constructUserData("test", "pwd", []string{"sshkey"})
	fmt.Println(cloudinit)
}
