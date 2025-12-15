package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestResetCommand(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	// Create a settings file with some data
	initialSettings := map[string]interface{}{
		"model": "some-model",
		"env": map[string]string{
			"KEY1": "value1",
			"KEY2": "value2",
		},
	}
	settingsData, _ := json.MarshalIndent(initialSettings, "", "    ")
	if err := os.WriteFile(settingsPath, settingsData, 0644); err != nil {
		t.Fatalf("Failed to write initial settings: %v", err)
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("settings", "s", settingsPath, "settings path")
	rootCmd.AddCommand(resetCmd)

	t.Run("reset settings", func(t *testing.T) {
		rootCmd.SetArgs([]string{"reset", "-s", settingsPath})
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("reset command failed: %v", err)
		}

		// Verify settings were reset
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			t.Fatalf("Failed to read settings: %v", err)
		}

		var settings map[string]interface{}
		if err := json.Unmarshal(data, &settings); err != nil {
			t.Fatalf("Failed to parse settings: %v", err)
		}

		// Settings should be empty after reset
		if len(settings) != 0 {
			t.Errorf("Settings should be empty after reset, got: %v", settings)
		}
	})
}

func TestResetCommandWithProfilesPath(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")
	profilesPath := filepath.Join(tmpDir, "profiles.json")

	// Create a profiles config
	profilesConfig := map[string]interface{}{
		"settingsPath": settingsPath,
		"profiles":     map[string]interface{}{},
	}
	profilesData, _ := json.MarshalIndent(profilesConfig, "", "    ")
	if err := os.WriteFile(profilesPath, profilesData, 0644); err != nil {
		t.Fatalf("Failed to write profiles config: %v", err)
	}

	// Create settings with data
	initialSettings := map[string]interface{}{
		"model": "test-model",
		"env": map[string]string{
			"TEST_KEY": "test_value",
		},
	}
	settingsData, _ := json.MarshalIndent(initialSettings, "", "    ")
	if err := os.WriteFile(settingsPath, settingsData, 0644); err != nil {
		t.Fatalf("Failed to write initial settings: %v", err)
	}

	// Set HOME to tmpDir so profiles are found
	originalHome := os.Getenv("HOME")
	ccswitchDir := filepath.Join(tmpDir, ".ccswitch")
	os.MkdirAll(ccswitchDir, 0755)
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Copy profiles to expected location
	if err := os.WriteFile(filepath.Join(ccswitchDir, "ccs.json"), profilesData, 0644); err != nil {
		t.Fatalf("Failed to write profiles to .ccswitch: %v", err)
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("settings", "s", "", "settings path")
	rootCmd.AddCommand(resetCmd)

	t.Run("reset with profiles path discovery", func(t *testing.T) {
		rootCmd.SetArgs([]string{"reset"})
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("reset command failed: %v", err)
		}

		// Just verify command executes without error
		// The settings file location is determined from profiles config
	})
}
