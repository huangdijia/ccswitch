package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	// Test creating new settings file
	settings, err := New(settingsPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if settings.Path != settingsPath {
		t.Errorf("settings.Path = %v, want %v", settings.Path, settingsPath)
	}

	// Verify file was created
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		t.Errorf("settings file was not created at %v", settingsPath)
	}
}

func TestNewWithTildeExpansion(t *testing.T) {
	tmpDir := t.TempDir()

	// Set HOME to temp directory for this test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	settingsPath := "~/test-settings.json"
	settings, err := New(settingsPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "test-settings.json")
	if settings.Path != expectedPath {
		t.Errorf("settings.Path = %v, want %v", settings.Path, expectedPath)
	}
}

func TestReadWrite(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	// Create settings
	settings := &ClaudeSettings{
		Path:  settingsPath,
		Model: "test-model",
		Env: map[string]interface{}{
			"KEY1": "value1",
			"KEY2": "value2",
		},
	}

	// Write settings
	err := settings.Write()
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Read settings back
	newSettings := &ClaudeSettings{Path: settingsPath}
	err = newSettings.Read()
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	// Verify model
	if newSettings.Model != settings.Model {
		t.Errorf("Model = %v, want %v", newSettings.Model, settings.Model)
	}

	// Verify env
	if len(newSettings.Env) != len(settings.Env) {
		t.Errorf("Env length = %v, want %v", len(newSettings.Env), len(settings.Env))
	}

	for key, value := range settings.Env {
		if newSettings.Env[key] != value {
			t.Errorf("Env[%v] = %v, want %v", key, newSettings.Env[key], value)
		}
	}
}

func TestWriteEmptyEnv(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	// Create settings with empty env
	settings := &ClaudeSettings{
		Path: settingsPath,
		Env:  make(map[string]interface{}),
	}

	// Write settings
	err := settings.Write()
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Read file and verify it's empty JSON object
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty JSON object, got %v", result)
	}
}

func TestWriteWithModel(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	// Create settings with model
	settings := &ClaudeSettings{
		Path:  settingsPath,
		Model: "claude-3-opus",
		Env:   make(map[string]interface{}),
	}

	// Write settings
	err := settings.Write()
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Read file and verify structure
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if result["model"] != "claude-3-opus" {
		t.Errorf("model = %v, want %v", result["model"], "claude-3-opus")
	}
}

func TestReadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "invalid.json")

	// Write invalid JSON
	if err := os.WriteFile(settingsPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Try to read
	settings := &ClaudeSettings{Path: settingsPath}
	err := settings.Read()
	if err == nil {
		t.Error("Read() expected error for invalid JSON, got nil")
	}
}
