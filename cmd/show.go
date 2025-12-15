package cmd

import (
	"fmt"
	"strings"

	"github.com/huangdijia/ccswitch/internal/profiles"
	"github.com/huangdijia/ccswitch/internal/settings"
	"github.com/spf13/cobra"
)

// Token prefixes that should not be masked when empty or just the prefix
var emptyTokenPrefixes = []string{"sk-", "ms-", "sk-kimi-"}

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
			profs, err := profiles.New(profilesPath)
			if err != nil {
				return err
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

			fmt.Println("Current Claude Settings:")
			fmt.Printf("  Settings file: %s\n", settingsPath)
			fmt.Println()

			if currentSettings.Model != "" {
				fmt.Printf("  Model: %s\n", currentSettings.Model)
			} else {
				fmt.Println("  Model: (default)")
			}

			if len(currentSettings.Env) > 0 {
				fmt.Println("\nEnvironment Variables:")
				for key, value := range currentSettings.Env {
					valStr := fmt.Sprintf("%v", value)
					// Hide sensitive information
					if strings.Contains(strings.ToLower(key), "token") || strings.Contains(strings.ToLower(key), "key") {
						valStr = maskSensitiveValue(valStr)
					}
					fmt.Printf("  %s: %s\n", key, valStr)
				}
			}

			return nil
		}

		// Show specific profile
		profileName := args[0]
		profs, err := profiles.New(profilesPath)
		if err != nil {
			return err
		}

		if !profs.Has(profileName) {
			fmt.Printf("Error: Profile '%s' not found.\n", profileName)
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
				if strings.Contains(strings.ToLower(key), "token") || strings.Contains(strings.ToLower(key), "key") {
					value = maskSensitiveValue(value)
				}
				fmt.Printf("  %s: %s\n", key, value)
			}
		} else {
			fmt.Println("  (no custom configuration)")
		}

		return nil
	},
}

func maskSensitiveValue(value string) string {
	// Check if value is empty or just a prefix
	if value == "" {
		return "(not set)"
	}
	for _, prefix := range emptyTokenPrefixes {
		if value == prefix {
			return "(not set)"
		}
	}

	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}

	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}

func init() {
	showCmd.Flags().BoolVarP(&showCurrent, "current", "c", false, "Show current Claude settings instead of a profile")
}
