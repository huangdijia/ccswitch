package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestInitCommand(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "test-config", "ccs.json")

	// Create a mock config file to copy from
	configDir := filepath.Join(tmpDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	mockConfig := map[string]interface{}{
		"settingsPath": "~/.claude/settings.json",
		"default":      "test",
		"profiles": map[string]interface{}{
			"test": map[string]string{
				"ANTHROPIC_MODEL": "test-model",
			},
		},
	}
	configData, _ := json.MarshalIndent(mockConfig, "", "    ")
	mockConfigPath := filepath.Join(configDir, "ccs.json")
	if err := os.WriteFile(mockConfigPath, configData, 0644); err != nil {
		t.Fatalf("Failed to write mock config: %v", err)
	}

	// Change to temp directory so config/ccs.json is found
	originalWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalWd)

	// Create a new root command with init command for testing
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("profiles", "p", profilesPath, "profiles path")
	rootCmd.AddCommand(initCmd)

	// Test 1: Initialize config successfully
	t.Run("initialize config successfully", func(t *testing.T) {
		rootCmd.SetArgs([]string{"init", "-p", profilesPath})
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("init command failed: %v", err)
		}

		// Verify config file was created
		if _, err := os.Stat(profilesPath); os.IsNotExist(err) {
			t.Error("Config file was not created")
		}
	})

	// Test 2: Error when config already exists without force flag
	t.Run("error without force flag", func(t *testing.T) {
		rootCmd.SetArgs([]string{"init", "-p", profilesPath})
		err := rootCmd.Execute()
		if err == nil {
			t.Error("Expected error when config already exists, got nil")
		}
	})

	// Test 3: Overwrite with force flag
	t.Run("overwrite with force flag", func(t *testing.T) {
		rootCmd.SetArgs([]string{"init", "-p", profilesPath, "--force"})
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("init command with force flag failed: %v", err)
		}
	})
}

func TestInitCommandFull(t *testing.T) {
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "ccs.json")

	// Create mock full config
	configDir := filepath.Join(tmpDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	mockConfig := map[string]interface{}{
		"settingsPath": "~/.claude/settings.json",
		"profiles": map[string]interface{}{
			"profile1": map[string]string{"key": "value1"},
			"profile2": map[string]string{"key": "value2"},
		},
	}
	configData, _ := json.MarshalIndent(mockConfig, "", "    ")
	mockConfigPath := filepath.Join(configDir, "ccs-full.json")
	if err := os.WriteFile(mockConfigPath, configData, 0644); err != nil {
		t.Fatalf("Failed to write mock config: %v", err)
	}

	originalWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalWd)

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("profiles", "p", profilesPath, "profiles path")
	rootCmd.AddCommand(initCmd)

	// Test with --full flag
	t.Run("initialize with full config", func(t *testing.T) {
		rootCmd.SetArgs([]string{"init", "-p", profilesPath, "--full"})
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("init command with --full flag failed: %v", err)
		}

		// Verify config file was created and has multiple profiles
		data, err := os.ReadFile(profilesPath)
		if err != nil {
			t.Fatalf("Failed to read created config: %v", err)
		}

		var config map[string]interface{}
		if err := json.Unmarshal(data, &config); err != nil {
			t.Fatalf("Failed to parse config: %v", err)
		}

		profiles, ok := config["profiles"].(map[string]interface{})
		if !ok || len(profiles) < 2 {
			t.Error("Full config should have multiple profiles")
		}
	})
}
