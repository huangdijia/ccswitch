package cmd

import (
	"fmt"

	"github.com/huangdijia/ccswitch/internal/profiles"
	"github.com/huangdijia/ccswitch/internal/settings"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [profile]",
	Short: "Switch the active Claude API profile",
	Long:  "This command allows you to set the active Claude API profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please specify a profile name")
		}

		profileName := args[0]
		profilesPath := cmd.Flag("profiles").Value.String()
		settingsPath := cmd.Flag("settings").Value.String()

		profs, err := profiles.New(profilesPath)
		if err != nil {
			return err
		}

		if !profs.Has(profileName) {
			fmt.Printf("Error: Profile '%s' not found.\n", profileName)
			fmt.Println("Available profiles:")
			for _, name := range profs.GetAll() {
				fmt.Printf("  - %s\n", name)
			}
			return fmt.Errorf("profile not found")
		}

		if settingsPath == "" {
			settingsPath = profs.GetSettingsPath()
		}
		if settingsPath == "" {
			settingsPath = "~/.claude/settings.json"
		}

		currentSettings, err := settings.New(settingsPath)
		if err != nil {
			return err
		}

		// Get the environment variables for the selected profile
		env := profs.Get(profileName)

		// Convert map[string]string to map[string]interface{}
		envInterface := make(map[string]interface{})
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

		fmt.Printf("✓ Successfully switched to profile: %s\n", profileName)

		// Show profile details
		if len(env) > 0 {
			fmt.Println("\nProfile details:")
			if url, ok := env["ANTHROPIC_BASE_URL"]; ok {
				fmt.Printf("  URL: %s\n", url)
			}
			if model, ok := env["ANTHROPIC_MODEL"]; ok {
				fmt.Printf("  Model: %s\n", model)
			}
			if fastModel, ok := env["ANTHROPIC_SMALL_FAST_MODEL"]; ok {
				fmt.Printf("  Fast Model: %s\n", fastModel)
			}
		}

		return nil
	},
}
