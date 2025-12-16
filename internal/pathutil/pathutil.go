package pathutil

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandHome expands the tilde (~) in a path to the user's home directory
func ExpandHome(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		return home, nil
	}

	return filepath.Join(home, path[1:]), nil
}

// EnsureDir creates a directory and all necessary parent directories if they don't exist
func EnsureDir(path string, perm os.FileMode) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, perm)
	}
	return nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DefaultSettingsPath returns the default Claude settings path
func DefaultSettingsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.claude/settings.json"
	}
	return filepath.Join(home, ".claude", "settings.json")
}

// DefaultProfilesPath returns the default profiles configuration path
func DefaultProfilesPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.ccswitch/ccs.json"
	}
	return filepath.Join(home, ".ccswitch", "ccs.json")
}
