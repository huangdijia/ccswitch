package cmd

import (
	"fmt"
	"strings"

	"github.com/huangdijia/ccswitch/internal/profiles"
	"github.com/spf13/cobra"
)

var profilesCmd = &cobra.Command{
	Use:     "profiles",
	Aliases: []string{"ls"},
	Short:   "List all available profiles",
	Long:    "This command allows you to list all configured Claude API profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesPath := cmd.Flag("profiles").Value.String()

		profs, err := profiles.New(profilesPath)
		if err != nil {
			return err
		}

		defaultProfile := profs.Default()
		profileData := profs.Data.Profiles
		descriptions := profs.Data.Descriptions

		fmt.Println("Available Claude API Profiles:")
		fmt.Println()

		// Print header
		fmt.Printf("┌%-20s┬%-30s┬%-40s┬%-20s┬%-10s┐\n",
			repeat("─", 20), repeat("─", 30), repeat("─", 40), repeat("─", 20), repeat("─", 10))
		fmt.Printf("│ %-18s │ %-28s │ %-38s │ %-18s │ %-8s │\n",
			"Profile", "Description", "URL", "Model", "Status")
		fmt.Printf("├%-20s┼%-30s┼%-40s┼%-20s┼%-10s┤\n",
			repeat("─", 20), repeat("─", 30), repeat("─", 40), repeat("─", 20), repeat("─", 10))

		// Print rows
		for name, config := range profileData {
			status := ""
			if name == defaultProfile {
				status = "Default"
			}
			url := config["ANTHROPIC_BASE_URL"]
			model := config["ANTHROPIC_MODEL"]
			description := descriptions[name]

			// Truncate long values
			if len(description) > 28 {
				description = description[:25] + "..."
			}
			if len(url) > 38 {
				url = url[:35] + "..."
			}
			if len(model) > 18 {
				model = model[:15] + "..."
			}

			fmt.Printf("│ %-18s │ %-28s │ %-38s │ %-18s │ %-8s │\n",
				name, description, url, model, status)
		}

		// Print footer
		fmt.Printf("└%-20s┴%-30s┴%-40s┴%-20s┴%-10s┘\n",
			repeat("─", 20), repeat("─", 30), repeat("─", 40), repeat("─", 20), repeat("─", 10))

		fmt.Println()
		fmt.Printf("Total profiles: %d\n", len(profileData))

		return nil
	},
}

func repeat(s string, count int) string {
	return strings.Repeat(s, count)
}
