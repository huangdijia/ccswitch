package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func setupTestEnvironment(t *testing.T) (string, string, string) {
	t.Helper()
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "profiles.json")
	settingsPath := filepath.Join(tmpDir, "settings.json")

	// Create test profiles config
	profilesConfig := map[string]interface{}{
		"settingsPath": settingsPath,
		"default":      "test-profile",
		"profiles": map[string]interface{}{
			"test-profile": map[string]string{
				"ANTHROPIC_BASE_URL":         "https://api.test.com",
				"ANTHROPIC_MODEL":            "test-model",
				"ANTHROPIC_SMALL_FAST_MODEL": "fast-model",
			},
			"another-profile": map[string]string{
				"ANTHROPIC_MODEL": "another-model",
			},
		},
	}

	profilesData, _ := json.MarshalIndent(profilesConfig, "", "    ")
	if err := os.WriteFile(profilesPath, profilesData, 0644); err != nil {
		t.Fatalf("Failed to write profiles config: %v", err)
	}

	return tmpDir, profilesPath, settingsPath
}

func TestUseCommand(t *testing.T) {
	_, profilesPath, settingsPath := setupTestEnvironment(t)

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("profiles", "p", profilesPath, "profiles path")
	rootCmd.PersistentFlags().StringP("settings", "s", settingsPath, "settings path")
	rootCmd.AddCommand(useCmd)

	t.Run("switch to valid profile", func(t *testing.T) {
		rootCmd.SetArgs([]string{"use", "test-profile", "-p", profilesPath, "-s", settingsPath})
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("use command failed: %v", err)
		}

		// Verify settings file was updated
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			t.Fatalf("Failed to read settings: %v", err)
		}

		var settings map[string]interface{}
		if err := json.Unmarshal(data, &settings); err != nil {
			t.Fatalf("Failed to parse settings: %v", err)
		}

		if settings["model"] != "test-model" {
			t.Errorf("Model not set correctly, got: %v", settings["model"])
		}

		env, ok := settings["env"].(map[string]interface{})
		if !ok {
			t.Fatal("Env not set in settings")
		}

		if env["ANTHROPIC_BASE_URL"] != "https://api.test.com" {
			t.Errorf("Base URL not set correctly, got: %v", env["ANTHROPIC_BASE_URL"])
		}
	})

	t.Run("error on non-existent profile", func(t *testing.T) {
		rootCmd.SetArgs([]string{"use", "non-existent", "-p", profilesPath, "-s", settingsPath})
		err := rootCmd.Execute()
		if err == nil {
			t.Error("Expected error for non-existent profile, got nil")
		}
	})

	t.Run("error without profile argument", func(t *testing.T) {
		rootCmd.SetArgs([]string{"use", "-p", profilesPath, "-s", settingsPath})
		err := rootCmd.Execute()
		if err == nil {
			t.Error("Expected error without profile argument, got nil")
		}
	})
}

func TestUseCommandProfileNotFound(t *testing.T) {
	_, profilesPath, settingsPath := setupTestEnvironment(t)

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("profiles", "p", profilesPath, "profiles path")
	rootCmd.PersistentFlags().StringP("settings", "s", settingsPath, "settings path")
	rootCmd.AddCommand(useCmd)

	// Test with missing profile
	t.Run("profile not found error message", func(t *testing.T) {
		rootCmd.SetArgs([]string{"use", "missing-profile", "-p", profilesPath, "-s", settingsPath})
		err := rootCmd.Execute()
		if err == nil {
			t.Error("Expected error for missing profile")
		}
	})
}
