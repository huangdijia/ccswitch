package output

import (
	"testing"
)

func TestMaskSensitiveValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "empty string",
			value: "",
			want:  "(not set)",
		},
		{
			name:  "sk- prefix only",
			value: "sk-",
			want:  "(not set)",
		},
		{
			name:  "ms- prefix only",
			value: "ms-",
			want:  "(not set)",
		},
		{
			name:  "sk-kimi- prefix only",
			value: "sk-kimi-",
			want:  "(not set)",
		},
		{
			name:  "short value",
			value: "abc123",
			want:  "******",
		},
		{
			name:  "long token",
			value: "sk-ant-1234567890abcdef",
			want:  "sk-a***************cdef",
		},
		{
			name:  "very long token",
			value: "sk-ant-api03-1234567890abcdefghijklmnopqrstuvwxyz",
			want:  "sk-a*****************************************wxyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskSensitiveValue(tt.value); got != tt.want {
				t.Errorf("MaskSensitiveValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSensitiveKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{
			name: "api key",
			key:  "ANTHROPIC_API_KEY",
			want: true,
		},
		{
			name: "token",
			key:  "ACCESS_TOKEN",
			want: true,
		},
		{
			name: "secret",
			key:  "CLIENT_SECRET",
			want: true,
		},
		{
			name: "password",
			key:  "USER_PASSWORD",
			want: true,
		},
		{
			name: "non-sensitive",
			key:  "ANTHROPIC_MODEL",
			want: false,
		},
		{
			name: "base url",
			key:  "ANTHROPIC_BASE_URL",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSensitiveKey(tt.key); got != tt.want {
				t.Errorf("IsSensitiveKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
