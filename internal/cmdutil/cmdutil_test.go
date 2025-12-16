package cmdutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSettingsPath(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		settingsPath string
		profilesPath string
		want         string
	}{
		{
			name:         "explicit settings path",
			settingsPath: "/explicit/path",
			profilesPath: "",
			want:         "/explicit/path",
		},
		{
			name:         "default when no profiles",
			settingsPath: "",
			profilesPath: filepath.Join(tmpDir, "nonexistent.json"),
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveSettingsPath(tt.settingsPath, tt.profilesPath)
			if tt.want != "" && got != tt.want {
				t.Errorf("ResolveSettingsPath() = %v, want %v", got, tt.want)
			}
			// If want is empty, we just check that got is not empty (it should be default)
			if tt.want == "" && got == "" {
				t.Errorf("ResolveSettingsPath() returned empty string")
			}
		})
	}
}

func TestLoadProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "profiles.json")

	// Create a valid profiles file
	profilesData := `{
		"default": "test",
		"profiles": {
			"test": {
				"ANTHROPIC_MODEL": "test-model"
			}
		}
	}`
	if err := os.WriteFile(profilesPath, []byte(profilesData), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("load valid profiles", func(t *testing.T) {
		profs, err := LoadProfiles(profilesPath)
		if err != nil {
			t.Errorf("LoadProfiles() error = %v", err)
			return
		}
		if profs == nil {
			t.Error("LoadProfiles() returned nil")
		}
	})

	t.Run("load nonexistent profiles", func(t *testing.T) {
		_, err := LoadProfiles(filepath.Join(tmpDir, "nonexistent.json"))
		if err == nil {
			t.Error("LoadProfiles() expected error for nonexistent file")
		}
	})
}

func TestLoadSettings(t *testing.T) {
	tmpDir := t.TempDir()
	settingsPath := filepath.Join(tmpDir, "settings.json")

	t.Run("load or create settings", func(t *testing.T) {
		settings, err := LoadSettings(settingsPath)
		if err != nil {
			t.Errorf("LoadSettings() error = %v", err)
			return
		}
		if settings == nil {
			t.Error("LoadSettings() returned nil")
		}
	})
}

func TestValidateProfile(t *testing.T) {
	tmpDir := t.TempDir()
	profilesPath := filepath.Join(tmpDir, "profiles.json")

	// Create a valid profiles file
	profilesData := `{
		"default": "test",
		"profiles": {
			"test": {
				"ANTHROPIC_MODEL": "test-model"
			},
			"another": {
				"ANTHROPIC_MODEL": "another-model"
			}
		}
	}`
	if err := os.WriteFile(profilesPath, []byte(profilesData), 0644); err != nil {
		t.Fatal(err)
	}

	profs, err := LoadProfiles(profilesPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("valid profile", func(t *testing.T) {
		err := ValidateProfile(profs, "test")
		if err != nil {
			t.Errorf("ValidateProfile() error = %v, expected no error", err)
		}
	})

	t.Run("invalid profile", func(t *testing.T) {
		err := ValidateProfile(profs, "nonexistent")
		if err == nil {
			t.Error("ValidateProfile() expected error for nonexistent profile")
		}
	})
}
