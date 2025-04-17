package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertGiga(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "int input",
			input:    5,
			expected: "5Gi",
		},
		{
			name:     "string input",
			input:    "10",
			expected: "10Gi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				result := ConvertGiga[int](v)
				assert.Equal(t, tt.expected, result)
			case string:
				result := ConvertGiga[string](v)
				assert.Equal(t, tt.expected, result)
			default:
				t.Errorf("Unsupported type in test case")
			}
		})
	}
}
