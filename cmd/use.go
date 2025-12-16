package cmd

import (
	"fmt"

	"github.com/huangdijia/ccswitch/internal/cmdutil"
	"github.com/huangdijia/ccswitch/internal/output"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [profile]",
	Short: "Switch the active Claude API profile",
	Long:  "This command allows you to set the active Claude API profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesPath := cmd.Flag("profiles").Value.String()
		settingsPath := cmd.Flag("settings").Value.String()

		profs, err := cmdutil.LoadProfiles(profilesPath)
		if err != nil {
			return err
		}

		var profileName string
		if len(args) == 0 {
			// Interactive selection when no profile is specified
			availableProfiles := profs.GetAll()
			if len(availableProfiles) == 0 {
				return fmt.Errorf("no profiles available")
			}

			fmt.Println("Available profiles:")
			for i, name := range availableProfiles {
				fmt.Printf("  %d. %s\n", i+1, name)
			}

			fmt.Print("\nSelect profile (enter number or name): ")
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return fmt.Errorf("invalid selection")
			}

			// Try to parse as number first
			var index int
			_, err = fmt.Sscanf(input, "%d", &index)
			if err == nil && index >= 1 && index <= len(availableProfiles) {
				profileName = availableProfiles[index-1]
			} else {
				// Treat as profile name
				profileName = input
			}
		} else {
			profileName = args[0]
		}

		if err := cmdutil.ValidateProfile(profs, profileName); err != nil {
			return err
		}

		settingsPath = cmdutil.ResolveSettingsPath(settingsPath, profilesPath)

		currentSettings, err := cmdutil.LoadSettings(settingsPath)
		if err != nil {
			return err
		}

		// Get the environment variables for the selected profile
		env := profs.Get(profileName)

		// Convert map[string]string to map[string]any
		envInterface := make(map[string]any)
		for k, v := range env {
			envInterface[k] = v
		}

		currentSettings.Env = envInterface

		// Handle model setting
		if model, ok := env["ANTHROPIC_MODEL"]; ok {
			currentSettings.Model = model
		} else {
			currentSettings.Model = ""
		}

		// Write settings
		if err := currentSettings.Write(); err != nil {
			return err
		}

		output.Success("Successfully switched to profile: %s", profileName)

		// Show profile details
		output.PrintProfileDetails(env)

		return nil
	},
}
