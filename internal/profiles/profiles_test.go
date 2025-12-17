package profiles

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func createTestProfilesFile(t *testing.T, tmpDir string, config *Config) string {
	t.Helper()
	profilesPath := filepath.Join(tmpDir, "profiles.json")
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if err := os.WriteFile(profilesPath, data, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return profilesPath
}

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	config := &Config{
		SettingsPath: "~/.claude/settings.json",
		Default:      "default",
		Profiles: map[string]map[string]string{
			"default": {
				"ANTHROPIC_BASE_URL": "https://api.anthropic.com",
				"ANTHROPIC_MODEL":    "opus",
			},
		},
		Descriptions: map[string]string{
			"default": "Default profile",
		},
	}

	profilesPath := createTestProfilesFile(t, tmpDir, config)

	// Test creating new Profiles instance
	profiles, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if profiles.Path != profilesPath {
		t.Errorf("profiles.Path = %v, want %v", profiles.Path, profilesPath)
	}

	if profiles.Data.Default != "default" {
		t.Errorf("profiles.Data.Default = %v, want %v", profiles.Data.Default, "default")
	}
}

func TestNewFileNotFound(t *testing.T) {
	_, err := New("/non/existent/file.json")
	if err == nil {
		t.Error("New() expected error for non-existent file, got nil")
	}
}

func TestGetSettingsPath(t *testing.T) {
	tmpDir := t.TempDir()
	expectedPath := "~/.claude/settings.json"
	config := &Config{
		SettingsPath: expectedPath,
		Profiles:     make(map[string]map[string]string),
	}

	profilesPath := createTestProfilesFile(t, tmpDir, config)
	profiles, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if profiles.GetSettingsPath() != expectedPath {
		t.Errorf("GetSettingsPath() = %v, want %v", profiles.GetSettingsPath(), expectedPath)
	}
}

func TestHas(t *testing.T) {
	tmpDir := t.TempDir()
	config := &Config{
		Profiles: map[string]map[string]string{
			"profile1": {"key": "value"},
			"profile2": {"key": "value"},
		},
	}

	profilesPath := createTestProfilesFile(t, tmpDir, config)
	profiles, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name     string
		profile  string
		expected bool
	}{
		{"existing profile1", "profile1", true},
		{"existing profile2", "profile2", true},
		{"non-existing", "profile3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := profiles.Has(tt.profile); got != tt.expected {
				t.Errorf("Has(%v) = %v, want %v", tt.profile, got, tt.expected)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name: "with default set",
			config: &Config{
				Default:  "custom",
				Profiles: make(map[string]map[string]string),
			},
			expected: "custom",
		},
		{
			name: "without default set",
			config: &Config{
				Default:  "",
				Profiles: make(map[string]map[string]string),
			},
			expected: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			profilesPath := createTestProfilesFile(t, tmpDir, tt.config)
			profiles, err := New(profilesPath)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if got := profiles.Default(); got != tt.expected {
				t.Errorf("Default() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tmpDir := t.TempDir()
	config := &Config{
		Profiles: map[string]map[string]string{
			"test": {
				"ANTHROPIC_BASE_URL": "https://api.test.com",
				"ANTHROPIC_MODEL":    "test-model",
			},
		},
	}

	profilesPath := createTestProfilesFile(t, tmpDir, config)
	profiles, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result := profiles.Get("test")
	if result["ANTHROPIC_BASE_URL"] != "https://api.test.com" {
		t.Errorf("Get() ANTHROPIC_BASE_URL = %v, want %v", result["ANTHROPIC_BASE_URL"], "https://api.test.com")
	}

	if result["ANTHROPIC_MODEL"] != "test-model" {
		t.Errorf("Get() ANTHROPIC_MODEL = %v, want %v", result["ANTHROPIC_MODEL"], "test-model")
	}
}

func TestGetWithModelFallback(t *testing.T) {
	tmpDir := t.TempDir()
	config := &Config{
		Profiles: map[string]map[string]string{
			"test": {
				"ANTHROPIC_MODEL": "main-model",
			},
		},
	}

	profilesPath := createTestProfilesFile(t, tmpDir, config)
	profiles, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result := profiles.Get("test")

	// Check that missing model keys are filled with ANTHROPIC_MODEL
	for _, key := range defaultModelKeys {
		if result[key] != "main-model" {
			t.Errorf("Get() %v = %v, want %v", key, result[key], "main-model")
		}
	}
}

func TestGetNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	config := &Config{
		Profiles: map[string]map[string]string{},
	}

	profilesPath := createTestProfilesFile(t, tmpDir, config)
	profiles, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result := profiles.Get("non-existent")
	if len(result) != 0 {
		t.Errorf("Get() for non-existent profile returned %v, want empty map", result)
	}
}

func TestGetAll(t *testing.T) {
	tmpDir := t.TempDir()
	config := &Config{
		Profiles: map[string]map[string]string{
			"profile1": {"key": "value1"},
			"profile2": {"key": "value2"},
			"profile3": {"key": "value3"},
		},
	}

	profilesPath := createTestProfilesFile(t, tmpDir, config)
	profiles, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result := profiles.GetAll()
	if len(result) != 3 {
		t.Errorf("GetAll() length = %v, want %v", len(result), 3)
	}

	// Check that all profile names are present
	nameMap := make(map[string]bool)
	for _, name := range result {
		nameMap[name] = true
	}

	for _, expected := range []string{"profile1", "profile2", "profile3"} {
		if !nameMap[expected] {
			t.Errorf("GetAll() missing profile %v", expected)
		}
	}
}

func TestAdd(t *testing.T) {
	tmpDir := t.TempDir()
	config := &Config{
		Profiles: map[string]map[string]string{
			"existing": {"key": "value"},
		},
		Descriptions: map[string]string{
			"existing": "Existing profile",
		},
	}

	profilesPath := createTestProfilesFile(t, tmpDir, config)
	profiles, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test adding a new profile
	newEnv := map[string]string{
		"ANTHROPIC_BASE_URL": "https://api.test.com",
		"ANTHROPIC_MODEL":    "test-model",
	}
	err = profiles.Add("newprofile", newEnv, "New test profile")
	if err != nil {
		t.Errorf("Add() error = %v", err)
	}

	// Verify the profile was added
	if !profiles.Has("newprofile") {
		t.Error("Add() did not add the profile")
	}

	// Verify the environment variables
	result := profiles.Get("newprofile")
	if result["ANTHROPIC_BASE_URL"] != "https://api.test.com" {
		t.Errorf("Add() ANTHROPIC_BASE_URL = %v, want %v", result["ANTHROPIC_BASE_URL"], "https://api.test.com")
	}

	// Verify the description
	if profiles.Data.Descriptions["newprofile"] != "New test profile" {
		t.Errorf("Add() description = %v, want %v", profiles.Data.Descriptions["newprofile"], "New test profile")
	}

	// Test adding a duplicate profile
	err = profiles.Add("existing", newEnv, "Duplicate")
	if err == nil {
		t.Error("Add() should return error for duplicate profile")
	}
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	config := &Config{
		SettingsPath: "~/.claude/settings.json",
		Default:      "default",
		Profiles: map[string]map[string]string{
			"default": {
				"ANTHROPIC_BASE_URL": "https://api.anthropic.com",
				"ANTHROPIC_MODEL":    "opus",
			},
		},
		Descriptions: map[string]string{
			"default": "Default profile",
		},
	}

	profilesPath := createTestProfilesFile(t, tmpDir, config)
	profiles, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Add a new profile
	newEnv := map[string]string{
		"ANTHROPIC_BASE_URL": "https://api.new.com",
		"ANTHROPIC_MODEL":    "new-model",
	}
	err = profiles.Add("newprofile", newEnv, "New profile")
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// Save the profiles
	err = profiles.Save()
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// Load the profiles again to verify they were saved
	profiles2, err := New(profilesPath)
	if err != nil {
		t.Fatalf("New() after Save() error = %v", err)
	}

	// Verify the new profile exists
	if !profiles2.Has("newprofile") {
		t.Error("Save() did not persist the new profile")
	}

	// Verify the environment variables
	result := profiles2.Get("newprofile")
	if result["ANTHROPIC_BASE_URL"] != "https://api.new.com" {
		t.Errorf("Save() ANTHROPIC_BASE_URL = %v, want %v", result["ANTHROPIC_BASE_URL"], "https://api.new.com")
	}

	// Verify the description
	if profiles2.Data.Descriptions["newprofile"] != "New profile" {
		t.Errorf("Save() description = %v, want %v", profiles2.Data.Descriptions["newprofile"], "New profile")
	}
}
