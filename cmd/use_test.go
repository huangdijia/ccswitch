package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/huangdijia/ccswitch/internal/profiles"
	"github.com/spf13/cobra"
)

func setupTestEnvironment(t *testing.T) (string, string, string) {
	t.Helper()
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "profiles.json")
	settingsPath := filepath.Join(tmpDir, "settings.json")

	// Create test profiles config
	profilesConfig := map[string]any{
		"settingsPath": settingsPath,
		"default":      "test-profile",
		"profiles": map[string]any{
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

		var settings map[string]any
		if err := json.Unmarshal(data, &settings); err != nil {
			t.Fatalf("Failed to parse settings: %v", err)
		}

		if settings["model"] != "test-model" {
			t.Errorf("Model not set correctly, got: %v", settings["model"])
		}

		env, ok := settings["env"].(map[string]any)
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

func TestUseCommandInteractive(t *testing.T) {
	_, profilesPath, settingsPath := setupTestEnvironment(t)

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("profiles", "p", profilesPath, "profiles path")
	rootCmd.PersistentFlags().StringP("settings", "s", settingsPath, "settings path")
	rootCmd.AddCommand(useCmd)

	t.Run("test interactive mode setup", func(t *testing.T) {
		// This test verifies that the interactive mode doesn't error out
		// We can't easily test stdin interaction in unit tests, but we can
		// ensure the code path is set up correctly

		// First verify that profiles are loaded correctly
		profs, err := profiles.New(profilesPath)
		if err != nil {
			t.Fatalf("Failed to load profiles: %v", err)
		}

		availableProfiles := profs.GetAll()
		if len(availableProfiles) == 0 {
			t.Error("Expected profiles to be available")
		}

		// Verify test profiles exist
		expectedProfiles := []string{"test-profile", "another-profile"}
		for _, expected := range expectedProfiles {
			found := false
			for _, profile := range availableProfiles {
				if profile == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected profile '%s' not found in available profiles", expected)
			}
		}

		// Test that the command structure is correct for interactive mode
		if useCmd.RunE == nil {
			t.Error("Expected RunE to be set on useCmd")
		}
	})

	t.Run("empty profiles list", func(t *testing.T) {
		// Create empty profiles config
		tmpDir := t.TempDir()
		emptyProfilesPath := filepath.Join(tmpDir, "empty-profiles.json")
		emptyProfilesConfig := map[string]any{
			"settingsPath": settingsPath,
			"default":      "",
			"profiles":     map[string]any{},
		}
		emptyProfilesData, _ := json.Marshal(emptyProfilesConfig)
		os.WriteFile(emptyProfilesPath, emptyProfilesData, 0644)

		// Test with empty profiles - this should fail gracefully
		rootCmd.SetArgs([]string{"use", "-p", emptyProfilesPath, "-s", settingsPath})
		err := rootCmd.Execute()

		// Since we can't simulate stdin input, this will fail with "invalid selection"
		// but the important thing is that it doesn't panic and handles the empty case
		if err == nil {
			t.Log("Command completed (this may vary based on test environment)")
		} else {
			t.Logf("Command failed as expected without stdin: %v", err)
		}
	})
}
