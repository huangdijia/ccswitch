package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestAddCommand(t *testing.T) {
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "profiles.json")

	// Create initial profiles config
	profilesConfig := map[string]interface{}{
		"settingsPath": "~/.claude/settings.json",
		"default":      "default",
		"profiles": map[string]interface{}{
			"default": map[string]string{
				"ANTHROPIC_BASE_URL": "https://api.anthropic.com",
				"ANTHROPIC_MODEL":    "opus",
			},
		},
		"descriptions": map[string]string{
			"default": "Default profile",
		},
	}

	profilesData, _ := json.MarshalIndent(profilesConfig, "", "    ")
	if err := os.WriteFile(profilesPath, profilesData, 0644); err != nil {
		t.Fatalf("Failed to write profiles config: %v", err)
	}

	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().StringP("profiles", "p", profilesPath, "profiles path")
	rootCmd.AddCommand(addCmd)

	t.Run("add profile with all flags", func(t *testing.T) {
		// Reset flags
		addAPIKey = ""
		addBaseURL = ""
		addModel = ""
		addDescription = ""
		addForce = false

		rootCmd.SetArgs([]string{
			"add", "testprofile",
			"-p", profilesPath,
			"--api-key", "sk-test-key",
			"--base-url", "https://api.test.com",
			"--model", "test-model",
			"--description", "Test profile",
		})

		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("add command failed: %v", err)
		}

		// Verify the profile was added
		data, err := os.ReadFile(profilesPath)
		if err != nil {
			t.Fatalf("Failed to read profiles: %v", err)
		}

		var config map[string]interface{}
		if err := json.Unmarshal(data, &config); err != nil {
			t.Fatalf("Failed to unmarshal profiles: %v", err)
		}

		profiles := config["profiles"].(map[string]interface{})
		if _, ok := profiles["testprofile"]; !ok {
			t.Error("Profile 'testprofile' was not added")
		}

		testProfile := profiles["testprofile"].(map[string]interface{})
		if testProfile["ANTHROPIC_API_KEY"] != "sk-test-key" {
			t.Errorf("ANTHROPIC_API_KEY = %v, want %v", testProfile["ANTHROPIC_API_KEY"], "sk-test-key")
		}
		if testProfile["ANTHROPIC_BASE_URL"] != "https://api.test.com" {
			t.Errorf("ANTHROPIC_BASE_URL = %v, want %v", testProfile["ANTHROPIC_BASE_URL"], "https://api.test.com")
		}
		if testProfile["ANTHROPIC_MODEL"] != "test-model" {
			t.Errorf("ANTHROPIC_MODEL = %v, want %v", testProfile["ANTHROPIC_MODEL"], "test-model")
		}

		descriptions := config["descriptions"].(map[string]interface{})
		if descriptions["testprofile"] != "Test profile" {
			t.Errorf("Description = %v, want %v", descriptions["testprofile"], "Test profile")
		}
	})

	t.Run("add profile without force flag should fail for duplicate", func(t *testing.T) {
		// Reset flags
		addAPIKey = ""
		addBaseURL = ""
		addModel = ""
		addDescription = ""
		addForce = false

		rootCmd.SetArgs([]string{
			"add", "testprofile",
			"-p", profilesPath,
			"--api-key", "sk-test-key-2",
			"--base-url", "https://api.test2.com",
			"--model", "test-model-2",
		})

		err := rootCmd.Execute()
		if err == nil {
			t.Error("add command should fail for duplicate profile without --force flag")
		}
	})

	t.Run("add profile with force flag should overwrite", func(t *testing.T) {
		// Reset flags
		addAPIKey = ""
		addBaseURL = ""
		addModel = ""
		addDescription = ""
		addForce = false

		rootCmd.SetArgs([]string{
			"add", "testprofile",
			"-p", profilesPath,
			"--api-key", "sk-test-key-updated",
			"--base-url", "https://api.updated.com",
			"--model", "updated-model",
			"--description", "Updated profile",
			"--force",
		})

		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("add command with --force failed: %v", err)
		}

		// Verify the profile was updated
		data, err := os.ReadFile(profilesPath)
		if err != nil {
			t.Fatalf("Failed to read profiles: %v", err)
		}

		var config map[string]interface{}
		if err := json.Unmarshal(data, &config); err != nil {
			t.Fatalf("Failed to unmarshal profiles: %v", err)
		}

		profiles := config["profiles"].(map[string]interface{})
		testProfile := profiles["testprofile"].(map[string]interface{})
		if testProfile["ANTHROPIC_API_KEY"] != "sk-test-key-updated" {
			t.Errorf("ANTHROPIC_API_KEY = %v, want %v", testProfile["ANTHROPIC_API_KEY"], "sk-test-key-updated")
		}
		if testProfile["ANTHROPIC_BASE_URL"] != "https://api.updated.com" {
			t.Errorf("ANTHROPIC_BASE_URL = %v, want %v", testProfile["ANTHROPIC_BASE_URL"], "https://api.updated.com")
		}

		descriptions := config["descriptions"].(map[string]interface{})
		if descriptions["testprofile"] != "Updated profile" {
			t.Errorf("Description = %v, want %v", descriptions["testprofile"], "Updated profile")
		}
	})

	t.Run("add profile with minimal flags", func(t *testing.T) {
		// Reset flags
		addAPIKey = ""
		addBaseURL = ""
		addModel = ""
		addDescription = ""
		addForce = false

		rootCmd.SetArgs([]string{
			"add", "minimal",
			"-p", profilesPath,
			"--base-url", "https://api.minimal.com",
		})

		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("add command failed: %v", err)
		}

		// Verify the profile was added with defaults
		data, err := os.ReadFile(profilesPath)
		if err != nil {
			t.Fatalf("Failed to read profiles: %v", err)
		}

		var config map[string]interface{}
		if err := json.Unmarshal(data, &config); err != nil {
			t.Fatalf("Failed to unmarshal profiles: %v", err)
		}

		profiles := config["profiles"].(map[string]interface{})
		minimalProfile := profiles["minimal"].(map[string]interface{})
		
		if minimalProfile["ANTHROPIC_BASE_URL"] != "https://api.minimal.com" {
			t.Errorf("ANTHROPIC_BASE_URL = %v, want %v", minimalProfile["ANTHROPIC_BASE_URL"], "https://api.minimal.com")
		}

		// Check for default model keys
		if minimalProfile["ANTHROPIC_DEFAULT_HAIKU_MODEL"] != "haiku" {
			t.Errorf("ANTHROPIC_DEFAULT_HAIKU_MODEL = %v, want %v", minimalProfile["ANTHROPIC_DEFAULT_HAIKU_MODEL"], "haiku")
		}
	})
}

func TestAddCommandName(t *testing.T) {
	if addCmd.Use != "add [profile-name]" {
		t.Errorf("Expected command name to be 'add [profile-name]', got '%s'", addCmd.Use)
	}
}
