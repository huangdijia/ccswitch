package cmd

import (
	"fmt"

	"github.com/huangdijia/ccswitch/internal/cmdutil"
	"github.com/huangdijia/ccswitch/internal/output"
	"github.com/spf13/cobra"
)

var (
	showCurrent bool
)

var showCmd = &cobra.Command{
	Use:   "show [profile]",
	Short: "Show details of a specific profile or current settings",
	Long:  "This command allows you to view details of a specific profile or current Claude settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesPath := cmd.Flag("profiles").Value.String()
		settingsPath := cmd.Flag("settings").Value.String()

		// Show current settings if requested or if no profile name provided
		if showCurrent || len(args) == 0 {
			settingsPath = cmdutil.ResolveSettingsPath(settingsPath, profilesPath)

			currentSettings, err := cmdutil.LoadSettings(settingsPath)
			if err != nil {
				return err
			}

			fmt.Println("Current Claude Settings:")
			fmt.Printf("  Settings file: %s\n", settingsPath)
			fmt.Println()

			if currentSettings.Model != "" {
				fmt.Printf("  Model: %s\n", currentSettings.Model)
			} else {
				fmt.Println("  Model: (default)")
			}

			output.PrintEnvVariables(currentSettings.Env)

			return nil
		}

		// Show specific profile
		profileName := args[0]
		profs, err := cmdutil.LoadProfiles(profilesPath)
		if err != nil {
			return err
		}

		if err := cmdutil.ValidateProfile(profs, profileName); err != nil {
			return err
		}

		profileData := profs.Get(profileName)
		descriptions := profs.Data.Descriptions

		fmt.Printf("Profile: %s\n", profileName)

		if profileName == profs.Default() {
			fmt.Println("  (default profile)")
		}

		if desc, ok := descriptions[profileName]; ok {
			fmt.Printf("  Description: %s\n", desc)
		}

		fmt.Println("\nConfiguration:")

		if len(profileData) > 0 {
			for key, value := range profileData {
				// Hide sensitive information
				if output.IsSensitiveKey(key) {
					value = output.MaskSensitiveValue(value)
				}
				fmt.Printf("  %s: %s\n", key, value)
			}
		} else {
			fmt.Println("  (no custom configuration)")
		}

		return nil
	},
}

func init() {
	showCmd.Flags().BoolVarP(&showCurrent, "current", "c", false, "Show current Claude settings instead of a profile")
}
