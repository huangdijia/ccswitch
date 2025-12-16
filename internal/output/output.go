package output

import (
	"fmt"
	"strings"
)

// Token prefixes that should not be masked when empty or just the prefix
var emptyTokenPrefixes = []string{"sk-", "ms-", "sk-kimi-"}

// MaskSensitiveValue masks sensitive information in a string value
// Used for API keys, tokens, and other sensitive data
func MaskSensitiveValue(value string) string {
	// Check if value is empty or just a prefix
	if value == "" {
		return "(not set)"
	}
	for _, prefix := range emptyTokenPrefixes {
		if value == prefix {
			return "(not set)"
		}
	}

	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}

	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}

// IsSensitiveKey checks if a key name indicates sensitive data
func IsSensitiveKey(key string) bool {
	lowerKey := strings.ToLower(key)
	return strings.Contains(lowerKey, "token") || 
	       strings.Contains(lowerKey, "key") ||
	       strings.Contains(lowerKey, "secret") ||
	       strings.Contains(lowerKey, "password")
}

// Success prints a success message with a checkmark
func Success(format string, args ...interface{}) {
	fmt.Printf("✓ "+format+"\n", args...)
}

// Error prints an error message
func Error(format string, args ...interface{}) {
	fmt.Printf("Error: "+format+"\n", args...)
}

// PrintProfileDetails prints profile details in a formatted way
func PrintProfileDetails(env map[string]string) {
	if len(env) == 0 {
		return
	}

	fmt.Println("\nProfile details:")
	if url, ok := env["ANTHROPIC_BASE_URL"]; ok {
		fmt.Printf("  URL: %s\n", url)
	}
	if model, ok := env["ANTHROPIC_MODEL"]; ok {
		fmt.Printf("  Model: %s\n", model)
	}
	if fastModel, ok := env["ANTHROPIC_SMALL_FAST_MODEL"]; ok {
		fmt.Printf("  Fast Model: %s\n", fastModel)
	}
}

// PrintEnvVariables prints environment variables with sensitive data masked
func PrintEnvVariables(env map[string]interface{}) {
	if len(env) == 0 {
		return
	}

	fmt.Println("\nEnvironment Variables:")
	for key, value := range env {
		valStr := fmt.Sprintf("%v", value)
		if IsSensitiveKey(key) {
			valStr = MaskSensitiveValue(valStr)
		}
		fmt.Printf("  %s: %s\n", key, valStr)
	}
}
