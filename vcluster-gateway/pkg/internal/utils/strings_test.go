package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToString(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected string
	}{
		{
			name:     "Empty map",
			input:    map[string]string{},
			expected: "",
		},
		{
			name: "Single key-value pair",
			input: map[string]string{
				"key1": "value1",
			},
			expected: "key1=value1",
		},
		{
			name: "Multiple key-value pairs",
			input: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			expected: "key1=value1,key2=value2,key3=value3",
		},
		{
			name: "Key with empty value",
			input: map[string]string{
				"key1": "value1",
				"key2": "",
			},
			expected: "key1=value1,key2=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := MapToString(tt.input)
			assert.Equal(t, tt.expected, output)
		})
	}
}
