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
		profilesPath := cmd.Flag("profiles").Value.String()
		settingsPath := cmd.Flag("settings").Value.String()

		profs, err := profiles.New(profilesPath)
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
