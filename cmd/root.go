package cmd

import (
	"fmt"
	"os"

	"github.com/huangdijia/ccswitch/internal/pathutil"
	"github.com/spf13/cobra"
)

var (
	profilesPath string
	settingsPath string
	// appVersion holds the version of the application
	appVersion string
	// appCommit holds the git commit hash
	appCommit string
	// appDate holds the build date
	appDate string
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
	defaultProfilesPath := pathutil.DefaultProfilesPath()

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&profilesPath, "profiles", "p", defaultProfilesPath, "Path to the profiles configuration file")
	rootCmd.PersistentFlags().StringVarP(&settingsPath, "settings", "s", "", "Path to the Claude settings file")

	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(updateCmd)
}

// SetVersion sets the application version, commit and build date
func SetVersion(v, c, d string) {
	appVersion = v
	appCommit = c
	appDate = d
	// Include commit and date in version string if available
	versionStr := appVersion
	if appCommit != "none" && appCommit != "" {
		versionStr += " (commit: " + appCommit[:8] + ")"
	}
	if appDate != "unknown" && appDate != "" {
		versionStr += ", built: " + appDate
	}
	rootCmd.Version = versionStr
}
