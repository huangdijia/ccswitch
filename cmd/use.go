package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/huangdijia/ccswitch/internal/cmdutil"
	"github.com/huangdijia/ccswitch/internal/output"
	"github.com/huangdijia/ccswitch/internal/termui"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [profile]",
	Short: "Switch the active Claude API profile",
	Long:  "This command allows you to set the active Claude API profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		profilesPath := cmd.Flag("profiles").Value.String()
		settingsPath := cmd.Flag("settings").Value.String()

		profs, err := cmdutil.LoadProfiles(profilesPath)
		if err != nil {
			return err
		}

		var profileName string
		if len(args) == 0 {
			// Interactive selection when no profile is specified.
			availableProfiles := profs.GetAll()
			if len(availableProfiles) == 0 {
				return fmt.Errorf("no profiles available")
			}

			sort.Strings(availableProfiles)

			// Default selection: default profile (if exists), else first.
			defaultIndex := 0
			defaultProfile := profs.Default()
			for i, name := range availableProfiles {
				if name == defaultProfile {
					defaultIndex = i
					break
				}
			}

			inFile, inOK := cmd.InOrStdin().(*os.File)
			outFile, outOK := cmd.OutOrStdout().(*os.File)
			if !inOK || !outOK {
				return fmt.Errorf("no profile specified (use 'ccswitch use <profile>')")
			}

			selected, err := termui.SelectString(termui.SelectConfig{
				In:           inFile,
				Out:          outFile,
				Prompt:       "Select profile:",
				Hint:         "↑/↓ to move, Enter to select, q to cancel",
				Items:        availableProfiles,
				DefaultIndex: defaultIndex,
			})
			if err != nil {
				if err == termui.ErrCanceled {
					return nil
				}
				if err == termui.ErrNotTerminal {
					return fmt.Errorf("no profile specified (use 'ccswitch use <profile>')")
				}
				return err
			}
			profileName = selected
		} else {
			profileName = args[0]
		}

		if err := cmdutil.ValidateProfile(profs, profileName); err != nil {
			return err
		}

		settingsPath = cmdutil.ResolveSettingsPath(settingsPath, profilesPath)

		currentSettings, err := cmdutil.LoadSettings(settingsPath)
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

		output.Success("Successfully switched to profile: %s", profileName)

		// Show profile details
		output.PrintProfileDetails(env)

		return nil
	},
}
