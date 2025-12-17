package cmd

import (
	"github.com/huangdijia/ccswitch/internal/cmdutil"
	"github.com/huangdijia/ccswitch/internal/output"
	"github.com/huangdijia/ccswitch/internal/pathutil"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Claude settings to default state",
	Long:  "This command resets your Claude settings to their default state",
	RunE: func(cmd *cobra.Command, args []string) error {
		settingsPath := cmd.Flag("settings").Value.String()
		profilesPath := cmd.Flag("profiles").Value.String()

		// If no settings path provided, try to resolve it
		if settingsPath == "" {
			settingsPath = cmdutil.ResolveSettingsPath(settingsPath, profilesPath)
		}

		if settingsPath == "" {
			settingsPath = pathutil.DefaultSettingsPath()
		}

		currentSettings, err := cmdutil.LoadSettings(settingsPath)
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

		output.Success("Settings have been reset to default")

		return nil
	},
}
