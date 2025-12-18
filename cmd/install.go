package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	installForce bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a profile from the full configuration",
	Long: `This command allows you to interactively select and install a profile from
the full configuration (ccs-full.json). It will download the configuration,
let you choose a profile, enter the authentication token, and save it to
your local configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesPath := cmd.Flag("profiles").Value.String()

		// Load existing profiles
		profs, err := cmdutil.LoadProfiles(profilesPath)
		if err != nil {
			return err
		}

		// Download ccs-full.json to temporary directory
		tmpDir, err := os.MkdirTemp("", "ccswitch-install-*")
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %w", err)
		}
		defer os.RemoveAll(tmpDir) // Clean up temporary files

		tmpConfigPath := filepath.Join(tmpDir, "ccs-full.json")
		
		repo := "huangdijia/ccswitch"
		branch := "main"
		githubURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/config/ccs-full.json", repo, branch)

		fmt.Println("Downloading full configuration from GitHub...")
		configContent, err := httputil.FetchBytes(githubURL)
		if err != nil {
			return fmt.Errorf("failed to download configuration from GitHub: %w", err)
		}

		// Write to temporary file
		if err := os.WriteFile(tmpConfigPath, configContent, 0644); err != nil {
			return fmt.Errorf("failed to write temporary configuration file: %w", err)
		}

		// Parse the full configuration
		var fullConfig profiles.Config
		if err := json.Unmarshal(configContent, &fullConfig); err != nil {
			return fmt.Errorf("failed to parse configuration: %w", err)
		}

		if len(fullConfig.Profiles) == 0 {
			return fmt.Errorf("no profiles found in the full configuration")
		}

		// Get list of profile names
		profileNames := make([]string, 0, len(fullConfig.Profiles))
		for name := range fullConfig.Profiles {
			profileNames = append(profileNames, name)
		}
		sort.Strings(profileNames)

		// Interactive profile selection
		inFile, inOK := cmd.InOrStdin().(*os.File)
		outFile, outOK := cmd.OutOrStdout().(*os.File)
		if !inOK || !outOK {
			return fmt.Errorf("interactive mode not available (use non-interactive flags instead)")
		}

		selected, err := termui.SelectString(termui.SelectConfig{
			In:           inFile,
			Out:          outFile,
			Prompt:       "Select profile to install:",
			Hint:         "↑/↓ to move, Enter to select, q to cancel",
			Items:        profileNames,
			DefaultIndex: 0,
		})
		if err != nil {
			if err == termui.ErrCanceled {
				fmt.Println("Installation canceled.")
				return nil
			}
			return err
		}

		profileName := selected

		// Check if profile already exists
		if profs.Has(profileName) && !installForce {
			return fmt.Errorf("profile '%s' already exists. Use --force to overwrite", profileName)
		}

		// Get the profile from full config
		selectedProfile := fullConfig.Profiles[profileName]
		description := fullConfig.Descriptions[profileName]

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

		// Remove old profile if it exists and force flag is set
		if profs.Has(profileName) && installForce {
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

		output.Success("Profile '%s' installed successfully!", profileName)

		// Show profile details
		output.PrintProfileDetails(env)

		return nil
	},
}

func init() {
	installCmd.Flags().BoolVarP(&installForce, "force", "f", false, "Force overwrite existing profile")
}
