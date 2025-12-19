package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/huangdijia/ccswitch/internal/cmdutil"
	"github.com/huangdijia/ccswitch/internal/httputil"
	"github.com/huangdijia/ccswitch/internal/output"
	"github.com/huangdijia/ccswitch/internal/profiles"
	"github.com/huangdijia/ccswitch/internal/termui"
	"github.com/spf13/cobra"
)

var (
	addAPIKey      string
	addBaseURL     string
	addModel       string
	addDescription string
	addForce       bool
	addOnline      bool
)

var addCmd = &cobra.Command{
	Use:     "add [profile-name]",
	Aliases: []string{"install"},
	Short:   "Add a new Claude API profile",
	Long: `Add a new Claude API profile with custom configuration or install from preset profiles.

Use --online flag to select from preset profiles available on GitHub.
Without --online flag, you can create a custom profile by providing your own configuration.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesPath := cmd.Flag("profiles").Value.String()

		profs, err := cmdutil.LoadProfiles(profilesPath)
		if err != nil {
			return err
		}

		// If --online flag is set, use online profile installation
		if addOnline {
			return installOnlineProfile(cmd, profs)
		}

		// Custom profile creation (existing functionality)
		return addCustomProfile(cmd, args, profs)
	},
}

// installOnlineProfile handles installation of profiles from preset.json
func installOnlineProfile(cmd *cobra.Command, profs *profiles.Profiles) error {
	// Download preset.json to temporary directory
	tmpDir, err := os.MkdirTemp("", "ccswitch-add-online-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up temporary files

	repo := "huangdijia/ccswitch"
	branch := "main"
	githubURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/config/preset.json", repo, branch)

	fmt.Println("Downloading preset configuration from GitHub...")
	configContent, err := httputil.FetchBytes(githubURL)
	if err != nil {
		return fmt.Errorf("failed to download configuration from GitHub: %w", err)
	}

	// Parse the preset configuration
	var presetConfig profiles.Config
	if err := json.Unmarshal(configContent, &presetConfig); err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	if len(presetConfig.Profiles) == 0 {
		return fmt.Errorf("no profiles found in the preset configuration")
	}

	// Get list of profile names
	profileNames := make([]string, 0, len(presetConfig.Profiles))
	for name := range presetConfig.Profiles {
		profileNames = append(profileNames, name)
	}
	sort.Strings(profileNames)

	// Interactive profile selection
	inFile, inOK := cmd.InOrStdin().(*os.File)
	outFile, outOK := cmd.OutOrStdout().(*os.File)
	if !inOK || !outOK {
		return fmt.Errorf("interactive mode requires an interactive terminal on stdin and stdout")
	}

	selected, err := termui.SelectString(termui.SelectConfig{
		In:           inFile,
		Out:          outFile,
		Prompt:       "Select profile to add:",
		Hint:         "↑/↓ to move, Enter to select, q to cancel",
		Items:        profileNames,
		DefaultIndex: 0,
	})
	if err != nil {
		if err == termui.ErrCanceled {
			fmt.Println("Operation canceled.")
			return nil
		}
		return err
	}

	profileName := selected

	// Check if profile already exists
	if profs.Has(profileName) && !addForce {
		return fmt.Errorf("profile '%s' already exists. Use --force to overwrite", profileName)
	}

	// Get the profile from preset config
	selectedProfile := presetConfig.Profiles[profileName]
	description := presetConfig.Descriptions[profileName]

	// Make a copy of the profile
	env := make(map[string]string)
	for k, v := range selectedProfile {
		env[k] = v
	}

	// Prompt for authentication token
	fmt.Printf("\nEnter authentication token for profile '%s'", profileName)
	if authKey, ok := env["ANTHROPIC_AUTH_TOKEN"]; ok && authKey != "" {
		fmt.Printf(" [current: %s]", output.MaskSensitiveValue(authKey))
	}
	fmt.Print(": ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read authentication token: %w", err)
	}
	authToken := strings.TrimSpace(input)

	// Update the auth token if provided
	if authToken != "" {
		env["ANTHROPIC_AUTH_TOKEN"] = authToken
	}

	// If overwriting, remove the old profile first
	if addForce && profs.Has(profileName) {
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
	output.PrintProfileDetails(env)

	return nil
}

// addCustomProfile handles creation of custom profiles
func addCustomProfile(cmd *cobra.Command, args []string, profs *profiles.Profiles) error {
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

	// Check if profile already exists and force flag is not set
	if profs.Has(profileName) && !addForce {
		return fmt.Errorf("profile '%s' already exists. Use --force to overwrite", profileName)
	}

	// If overwriting, remove the old profile first
	if profs.Has(profileName) && addForce {
		delete(profs.Data.Profiles, profileName)
		delete(profs.Data.Descriptions, profileName)
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

	// Add the profile (now guaranteed not to exist due to earlier checks)
	if err := profs.Add(profileName, env, description); err != nil {
		return err
	}

	// Save the profiles
	if err := profs.Save(); err != nil {
		return err
	}

	output.Success("Profile '%s' added successfully!", profileName)

	// Show profile details
	output.PrintProfileDetails(env)

	return nil
}

func init() {
	addCmd.Flags().BoolVarP(&addOnline, "online", "o", false, "Install a profile from online preset configuration")
	addCmd.Flags().StringVarP(&addAPIKey, "api-key", "k", "", "Anthropic API key (for custom profiles)")
	addCmd.Flags().StringVarP(&addBaseURL, "base-url", "u", "", "Anthropic base URL (for custom profiles)")
	addCmd.Flags().StringVarP(&addModel, "model", "m", "", "Anthropic model (for custom profiles)")
	addCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Profile description (for custom profiles)")
	addCmd.Flags().BoolVarP(&addForce, "force", "f", false, "Force overwrite existing profile")
}
