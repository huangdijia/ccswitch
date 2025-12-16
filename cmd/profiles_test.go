package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestProfilesCommand(t *testing.T) {
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "profiles.json")

	// Create test profiles config with multiple profiles
	profilesConfig := map[string]interface{}{
		"settingsPath": "~/.claude/settings.json",
		"default":      "default",
		"profiles": map[string]interface{}{
			"default": map[string]string{
				"ANTHROPIC_BASE_URL": "https://api.anthropic.com",
				"ANTHROPIC_MODEL":    "opus",
			},
			"custom": map[string]string{
				"ANTHROPIC_BASE_URL": "https://api.custom.com",
				"ANTHROPIC_MODEL":    "custom-model",
			},
		},
		"descriptions": map[string]string{
			"default": "Default profile",
			"custom":  "Custom profile",
		},
	}

	profilesData, _ := json.MarshalIndent(profilesConfig, "", "    ")
	if err := os.WriteFile(profilesPath, profilesData, 0644); err != nil {
		t.Fatalf("Failed to write profiles config: %v", err)
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("profiles", "p", profilesPath, "profiles path")
	rootCmd.AddCommand(profilesCmd)

	t.Run("list all profiles", func(t *testing.T) {
		rootCmd.SetArgs([]string{"profiles", "-p", profilesPath})
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("profiles command failed: %v", err)
		}
		// Just verify command executes without error
		// Output is printed directly, which is expected behavior
	})
}

func TestProfilesCommandWithNoProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "profiles.json")

	// Create config with no profiles
	profilesConfig := map[string]interface{}{
		"settingsPath": "~/.claude/settings.json",
		"profiles":     map[string]interface{}{},
	}

	profilesData, _ := json.MarshalIndent(profilesConfig, "", "    ")
	if err := os.WriteFile(profilesPath, profilesData, 0644); err != nil {
		t.Fatalf("Failed to write profiles config: %v", err)
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("profiles", "p", profilesPath, "profiles path")
	rootCmd.AddCommand(profilesCmd)

	t.Run("list with no profiles", func(t *testing.T) {
		rootCmd.SetArgs([]string{"profiles", "-p", profilesPath})
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("profiles command failed: %v", err)
		}
		// Just verify command executes without error
	})
}

func TestProfilesCommandAliases(t *testing.T) {
	// Test that 'ls' alias works
	if len(profilesCmd.Aliases) == 0 {
		t.Error("profiles command should have aliases")
	}

	hasLsAlias := false
	for _, alias := range profilesCmd.Aliases {
		if alias == "ls" {
			hasLsAlias = true
			break
		}
	}

	if !hasLsAlias {
		t.Error("profiles command should have 'ls' alias")
	}
}
