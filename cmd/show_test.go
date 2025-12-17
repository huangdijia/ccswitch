package cmd

import (
	"testing"

	"github.com/huangdijia/ccswitch/internal/output"
)

func TestMaskSensitiveValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "(not set)",
		},
		{
			name:     "sk- prefix only",
			input:    "sk-",
			expected: "(not set)",
		},
		{
			name:     "ms- prefix only",
			input:    "ms-",
			expected: "(not set)",
		},
		{
			name:     "sk-kimi- prefix only",
			input:    "sk-kimi-",
			expected: "(not set)",
		},
		{
			name:     "short value",
			input:    "short",
			expected: "*****",
		},
		{
			name:     "long token",
			input:    "sk-1234567890abcdef",
			expected: "sk-1***********cdef",
		},
		{
			name:     "very long token",
			input:    "sk-ant-api03-1234567890abcdefghijklmnopqrstuvwxyz",
			expected: "sk-a*****************************************wxyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := output.MaskSensitiveValue(tt.input)
			if result != tt.expected {
				t.Errorf("output.MaskSensitiveValue(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
