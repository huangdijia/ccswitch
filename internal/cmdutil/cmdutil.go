package cmdutil

import (
	"fmt"

	"github.com/huangdijia/ccswitch/internal/output"
	"github.com/huangdijia/ccswitch/internal/pathutil"
	"github.com/huangdijia/ccswitch/internal/profiles"
	"github.com/huangdijia/ccswitch/internal/settings"
)

// LoadProfiles loads profiles from the given path with error handling
func LoadProfiles(profilesPath string) (*profiles.Profiles, error) {
	profs, err := profiles.New(profilesPath)
	if err != nil {
		return nil, err
	}
	return profs, nil
}

// LoadSettings loads settings from the given path with error handling
func LoadSettings(settingsPath string) (*settings.ClaudeSettings, error) {
	currentSettings, err := settings.New(settingsPath)
	if err != nil {
		return nil, err
	}
	return currentSettings, nil
}

// ResolveSettingsPath resolves the settings path from profiles or uses default
func ResolveSettingsPath(settingsPath, profilesPath string) string {
	if settingsPath != "" {
		return settingsPath
	}

	// Try to get from profiles
	if pathutil.FileExists(profilesPath) {
		if profs, err := profiles.New(profilesPath); err == nil {
			if path := profs.GetSettingsPath(); path != "" {
				return path
			}
		}
	}

	// Return default
	return pathutil.DefaultSettingsPath()
}

// ValidateProfile validates that a profile exists and returns error with suggestions if not
func ValidateProfile(profs *profiles.Profiles, profileName string) error {
	if !profs.Has(profileName) {
		output.Error("Profile '%s' not found.", profileName)
		fmt.Println("Available profiles:")
		for _, name := range profs.GetAll() {
			marker := "  "
			if name == profs.Default() {
				marker = " *"
			}
			fmt.Printf("%s%s\n", marker, name)
		}
		return fmt.Errorf("profile not found")
	}
	return nil
}
