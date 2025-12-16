package cmd

import (
	"fmt"
	"os"

	"github.com/huangdijia/ccswitch/internal/profiles"
	"github.com/huangdijia/ccswitch/internal/settings"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Claude settings to default state",
	Long:  "This command resets your Claude settings to their default state",
	RunE: func(cmd *cobra.Command, args []string) error {
		settingsPath := cmd.Flag("settings").Value.String()

		// If no settings path provided, try to get it from profiles
		if settingsPath == "" {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				profilesPath := homeDir + "/.ccswitch/ccs.json"
				if _, err := os.Stat(profilesPath); err == nil {
					profs, err := profiles.New(profilesPath)
					if err == nil {
						settingsPath = profs.GetSettingsPath()
					}
				}
			}
		}

		if settingsPath == "" {
			settingsPath = "~/.claude/settings.json"
		}

		currentSettings, err := settings.New(settingsPath)
		if err != nil {
			return err
		}

		// Reset settings to empty state
		currentSettings.Env = make(map[string]interface{})
		currentSettings.Model = ""

		// Write the reset settings
		if err := currentSettings.Write(); err != nil {
			return err
		}

		fmt.Println("✓ Settings have been reset to default")

		return nil
	},
}
