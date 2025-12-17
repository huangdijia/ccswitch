package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestListCommand(t *testing.T) {
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
	rootCmd.AddCommand(listCmd)

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
	rootCmd.AddCommand(listCmd)

	t.Run("list with no profiles", func(t *testing.T) {
		rootCmd.SetArgs([]string{"profiles", "-p", profilesPath})
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("profiles command failed: %v", err)
		}
		// Just verify command executes without error
	})
}

// Test that command name is now 'list'
func TestListCommandName(t *testing.T) {
	if listCmd.Use != "list" {
		t.Errorf("Expected command name to be 'list', got '%s'", listCmd.Use)
	}
}

func TestProfilesCommandAliases(t *testing.T) {
	// Test that 'ls' alias works
	if len(listCmd.Aliases) == 0 {
		t.Error("list command should have aliases")
	}

	hasLsAlias := false
	hasProfilesAlias := false
	for _, alias := range listCmd.Aliases {
		if alias == "ls" {
			hasLsAlias = true
		}
		if alias == "profiles" {
			hasProfilesAlias = true
		}
	}

	if !hasLsAlias {
		t.Error("list command should have 'ls' alias")
	}
	if !hasProfilesAlias {
		t.Error("list command should have 'profiles' alias")
	}
}
