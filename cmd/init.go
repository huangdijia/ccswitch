package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	initForce bool
	initFull  bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ccswitch configuration",
	Long:  "This command initializes the ccswitch configuration by creating the necessary configuration files and directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesPath := cmd.Flag("profiles").Value.String()
		configDir := filepath.Dir(profilesPath)

		// Check if configuration already exists
		if _, err := os.Stat(profilesPath); err == nil && !initForce {
			return fmt.Errorf("configuration file already exists. Use --force to overwrite")
		}

		// Create config directory if it doesn't exist
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			fmt.Printf("Created directory: %s\n", configDir)
		}

		// Determine source config file
		sourceConfig := "config/ccs.json"
		if initFull {
			sourceConfig = "config/ccs-full.json"
		}

		// Try to find config in various locations
		possiblePaths := []string{
			sourceConfig,
			filepath.Join("/usr/local/share/ccswitch", filepath.Base(sourceConfig)),
			filepath.Join("/usr/share/ccswitch", filepath.Base(sourceConfig)),
		}

		// Try to find executable path and look for config relative to it
		if exePath, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exePath)
			possiblePaths = append([]string{
				filepath.Join(exeDir, sourceConfig),
				filepath.Join(exeDir, "..", sourceConfig),
			}, possiblePaths...)
		}

		var configContent []byte
		var foundPath string
		for _, path := range possiblePaths {
			if data, err := os.ReadFile(path); err == nil {
				configContent = data
				foundPath = path
				break
			}
		}

		if configContent == nil {
			// Try to download from GitHub
			repo := "huangdijia/ccswitch"
			branch := "main"
			configFile := "ccs.json"
			if initFull {
				configFile = "ccs-full.json"
			}
			githubURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/config/%s", repo, branch, configFile)

			fmt.Printf("Downloading configuration from GitHub...\n")
			resp, err := http.Get(githubURL)
			if err != nil {
				return fmt.Errorf("failed to download configuration from GitHub: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("GitHub API returned status: %s", resp.Status)
			}

			configContent, err = io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read downloaded configuration: %w", err)
			}
			foundPath = githubURL
		}

		// Write configuration file
		if err := os.WriteFile(profilesPath, configContent, 0644); err != nil {
			return fmt.Errorf("failed to write configuration file: %w", err)
		}

		configType := "default"
		if initFull {
			configType = "full"
		}
		fmt.Printf("✓ %s configuration file created successfully: %s\n", configType, profilesPath)
		if strings.HasPrefix(foundPath, "http") {
			fmt.Printf("  (downloaded from: %s)\n", foundPath)
		} else {
			fmt.Printf("  (copied from: %s)\n", foundPath)
		}

		return nil
	},
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Force overwrite existing configuration")
	initCmd.Flags().BoolVar(&initFull, "full", false, "Use full configuration with all available providers")
}
