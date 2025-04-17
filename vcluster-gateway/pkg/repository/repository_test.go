package repository

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateIdWithPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		charset      string
		length       int
		expectPrefix bool
	}{
		{
			name:         "Valid ID with prefix",
			prefix:       "vc",
			charset:      "abcdefghijklmnopqrstuvwxyz0123456789",
			length:       12,
			expectPrefix: true,
		},
		{
			name:         "Empty prefix",
			prefix:       "",
			charset:      "abcdefghijklmnopqrstuvwxyz0123456789",
			length:       12,
			expectPrefix: true,
		},
		{
			name:         "Short ID length",
			prefix:       "vc",
			charset:      "abcdefghijklmnopqrstuvwxyz0123456789",
			length:       5,
			expectPrefix: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := generateIdWithPrefix(tt.prefix, tt.charset, tt.length-len(tt.prefix))

			t.Logf("prefix: '%s', id: %s, len: %d", tt.prefix, id, len(id))

			assert.NoError(t, err, "did not expect error but got one")

			assert.True(t, strings.HasPrefix(id, tt.prefix), "expected ID to have prefix %s, got %s", tt.prefix, id)

			assert.Equal(t, tt.length, len(id), "expected ID length to be %d, got %d", tt.length, len(id))

			for _, char := range id[len(tt.prefix):] {
				assert.Contains(t, tt.charset, string(char), "unexpected character %c in ID %s", char, id)
			}
		})
	}
}
