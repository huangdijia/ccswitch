package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	profilesPath string
	settingsPath string
	// appVersion holds the version of the application
	appVersion string
)

var rootCmd = &cobra.Command{
	Use:   "ccswitch",
	Short: "A command-line tool for managing and switching between different Claude Code API profiles",
	Long: `CCSwitch allows you to easily manage multiple Claude Code API configurations (profiles)
and switch between them. This is useful when you need to use different API endpoints,
models, or authentication tokens for different projects or environments.`,
	Version: appVersion,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}

	defaultProfilesPath := homeDir + "/.ccswitch/ccs.json"

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&profilesPath, "profiles", "p", defaultProfilesPath, "Path to the profiles configuration file")
	rootCmd.PersistentFlags().StringVarP(&settingsPath, "settings", "s", "", "Path to the Claude settings file")

	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(profilesCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(updateCmd)
}

// SetVersion sets the application version
func SetVersion(v string) {
	appVersion = v
	rootCmd.Version = appVersion
}
