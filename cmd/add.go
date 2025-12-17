package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/huangdijia/ccswitch/internal/cmdutil"
	"github.com/huangdijia/ccswitch/internal/output"
	"github.com/spf13/cobra"
)

var (
	addAPIKey     string
	addBaseURL    string
	addModel      string
	addDescription string
	addForce      bool
)

var addCmd = &cobra.Command{
	Use:   "add [profile-name]",
	Short: "Add a new Claude API profile",
	Long:  "This command allows you to add a new Claude API profile with custom configuration",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesPath := cmd.Flag("profiles").Value.String()

		profs, err := cmdutil.LoadProfiles(profilesPath)
		if err != nil {
			return err
		}

		// Get profile name
		var profileName string
		if len(args) > 0 {
			profileName = args[0]
		} else {
			// Prompt for profile name
			fmt.Print("Enter profile name: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read profile name: %w", err)
			}
			profileName = strings.TrimSpace(input)
			if profileName == "" {
				return fmt.Errorf("profile name cannot be empty")
			}
		}

		// Check if profile already exists
		if profs.Has(profileName) && !addForce {
			return fmt.Errorf("profile '%s' already exists. Use --force to overwrite", profileName)
		}

		// Determine if we're in interactive mode (no flags provided)
		reader := bufio.NewReader(os.Stdin)
		
		// Get API key
		apiKey := addAPIKey
		if apiKey == "" && cmd.Flags().Changed("api-key") == false {
			// Only prompt if stdin is not a pipe
			if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
				fmt.Print("Enter ANTHROPIC_API_KEY (or press Enter to skip): ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read API key: %w", err)
				}
				apiKey = strings.TrimSpace(input)
			}
		}

		// Get base URL
		baseURL := addBaseURL
		if baseURL == "" {
			if cmd.Flags().Changed("base-url") == false {
				// Only prompt if stdin is not a pipe
				if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
					fmt.Print("Enter ANTHROPIC_BASE_URL [https://api.anthropic.com]: ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read base URL: %w", err)
					}
					baseURL = strings.TrimSpace(input)
				}
			}
			if baseURL == "" {
				baseURL = "https://api.anthropic.com"
			}
		}

		// Get model
		model := addModel
		if model == "" {
			if cmd.Flags().Changed("model") == false {
				// Only prompt if stdin is not a pipe
				if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
					fmt.Print("Enter ANTHROPIC_MODEL [opus]: ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read model: %w", err)
					}
					model = strings.TrimSpace(input)
				}
			}
			if model == "" {
				model = "opus"
			}
		}

		// Get description
		description := addDescription
		if description == "" && cmd.Flags().Changed("description") == false {
			// Only prompt if stdin is not a pipe
			if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
				fmt.Print("Enter description (optional): ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read description: %w", err)
				}
				description = strings.TrimSpace(input)
			}
		}

		// Create environment map
		env := make(map[string]string)
		if apiKey != "" {
			env["ANTHROPIC_API_KEY"] = apiKey
		}
		env["ANTHROPIC_BASE_URL"] = baseURL
		env["ANTHROPIC_MODEL"] = model

		// Add default model keys
		env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = "haiku"
		env["ANTHROPIC_DEFAULT_OPUS_MODEL"] = "opus"
		env["ANTHROPIC_DEFAULT_SONNET_MODEL"] = "sonnet"
		env["ANTHROPIC_SMALL_FAST_MODEL"] = "haiku"

		// If overwriting, remove the old profile first
		if profs.Has(profileName) && addForce {
			delete(profs.Data.Profiles, profileName)
			delete(profs.Data.Descriptions, profileName)
		}

		// Add the profile
		if err := profs.Add(profileName, env, description); err != nil {
			return err
		}

		// Save the profiles
		if err := profs.Save(); err != nil {
			return err
		}

		output.Success("Profile '%s' added successfully!", profileName)

		// Show profile details
		fmt.Println("\nProfile details:")
		output.PrintProfileDetails(env)

		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addAPIKey, "api-key", "k", "", "Anthropic API key")
	addCmd.Flags().StringVarP(&addBaseURL, "base-url", "u", "", "Anthropic base URL")
	addCmd.Flags().StringVarP(&addModel, "model", "m", "", "Anthropic model")
	addCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Profile description")
	addCmd.Flags().BoolVarP(&addForce, "force", "f", false, "Force overwrite existing profile")
}
