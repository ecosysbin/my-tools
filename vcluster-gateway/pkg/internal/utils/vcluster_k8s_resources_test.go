package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVClusterNamespaceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty vclusterId",
			input:    "",
			expected: "vcluster-",
		},
		{
			name:     "Normal vclusterId",
			input:    "12345",
			expected: "vcluster-12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := GetVClusterNamespaceName(tt.input)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestGetVClusterSecretName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty vclusterId",
			input:    "",
			expected: "vc-",
		},
		{
			name:     "Normal vclusterId",
			input:    "12345",
			expected: "vc-12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := GetVClusterSecretName(tt.input)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestGetVClusterServiceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty vclusterId",
			input:    "",
			expected: "",
		},
		{
			name:     "Normal vclusterId",
			input:    "12345",
			expected: "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := GetVClusterServiceName(tt.input)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestGetVClusterResourceQuotaName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty vclusterId",
			input:    "",
			expected: "-quota",
		},
		{
			name:     "Normal vclusterId",
			input:    "12345",
			expected: "12345-quota",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := GetVClusterResourceQuotaName(tt.input)
			assert.Equal(t, tt.expected, output)
		})
	}
}
