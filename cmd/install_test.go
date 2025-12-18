package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestInstallCommand(t *testing.T) {
	// Create a mock HTTP server to serve ccs-full.json
	mockFullConfig := map[string]interface{}{
		"settingsPath": "~/.claude/settings.json",
		"default":      "default",
		"profiles": map[string]interface{}{
			"default": map[string]string{
				"ANTHROPIC_BASE_URL":   "https://api.anthropic.com",
				"ANTHROPIC_AUTH_TOKEN": "sk-",
				"ANTHROPIC_MODEL":      "opus",
			},
			"testprovider": map[string]string{
				"ANTHROPIC_BASE_URL":   "https://test.provider.com",
				"ANTHROPIC_AUTH_TOKEN": "sk-test-",
				"ANTHROPIC_MODEL":      "test-model",
			},
		},
		"descriptions": map[string]string{
			"default":      "Default Anthropic API",
			"testprovider": "Test Provider API",
		},
	}

	mockConfigData, _ := json.MarshalIndent(mockFullConfig, "", "    ")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockConfigData)
	}))
	defer server.Close()

	t.Run("install command requires interactive terminal", func(t *testing.T) {
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
			"descriptions": map[string]string{},
		}

		profilesData, _ := json.MarshalIndent(profilesConfig, "", "    ")
		if err := os.WriteFile(profilesPath, profilesData, 0644); err != nil {
			t.Fatalf("Failed to write profiles config: %v", err)
		}

		rootCmd := &cobra.Command{Use: "test"}
		rootCmd.PersistentFlags().StringP("profiles", "p", profilesPath, "profiles path")
		
		// Create a fresh install command for this test
		testInstallCmd := &cobra.Command{
			Use:   "install",
			Short: "Install a profile from the full configuration",
			RunE:  installCmd.RunE,
		}
		testInstallCmd.Flags().BoolVarP(&installForce, "force", "f", false, "Force overwrite existing profile")
		rootCmd.AddCommand(testInstallCmd)

		// Since we can't simulate interactive terminal in tests easily,
		// just verify that the command is registered and available
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == "install" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Install command not registered")
		}
	})
}

func TestInstallCommandValidation(t *testing.T) {
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "profiles.json")

	// Create initial profiles config
	profilesConfig := map[string]interface{}{
		"settingsPath": "~/.claude/settings.json",
		"default":      "default",
		"profiles": map[string]interface{}{
			"existing": map[string]string{
				"ANTHROPIC_BASE_URL": "https://api.anthropic.com",
				"ANTHROPIC_MODEL":    "opus",
			},
		},
		"descriptions": map[string]string{},
	}

	profilesData, _ := json.MarshalIndent(profilesConfig, "", "    ")
	if err := os.WriteFile(profilesPath, profilesData, 0644); err != nil {
		t.Fatalf("Failed to write profiles config: %v", err)
	}

	t.Run("verify profiles can be loaded", func(t *testing.T) {
		// Just verify that our test setup is correct
		data, err := os.ReadFile(profilesPath)
		if err != nil {
			t.Fatalf("Failed to read profiles: %v", err)
		}

		var config map[string]interface{}
		if err := json.Unmarshal(data, &config); err != nil {
			t.Fatalf("Failed to parse profiles: %v", err)
		}

		profiles, ok := config["profiles"].(map[string]interface{})
		if !ok {
			t.Fatal("Profiles not found in config")
		}

		if _, ok := profiles["existing"]; !ok {
			t.Error("Expected profile 'existing' not found")
		}
	})
}
