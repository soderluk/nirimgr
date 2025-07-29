// Package cmd contains all the commands nirimgr supports.
//
// The root command just specifies nirimgr cli-name. Use the sub-commands to use nirimgr.
//
// # Events
//
// The events command starts listening on the Niri event stream.
//
//	Usage: nirimgr events
//
// # List
//
// The list command lists all defined events and actions.
//
//	Usage: nirimgr list [actions|events]
//
// # Scratch
//
// The scratch is the command to move a window to the scratchpad workspace,
// or show a window from the scratchpad workspace.
// Added in v0.3.0: spawn-or-focus [app-id]
// Using the spawn-or-focus [app-id] will either spawn a specific app, or focus it if it's already open.
//
//	Usage: nirimgr scratch [move|show|spawn-or-focus [app-id]]
package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/soderluk/nirimgr/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nirimgr",
	Short: "Commands for managing Niri",
	Long: `The nirimgr command can be used to listen to the niri event stream, and
		do something when an event comes in. E.g. when a window title hasn't yet been set
		when the window opens the first time, niri's own window-rule might not pick up on it.
		This command can handle the actions to be done for such cases. E.g. set a window to
		floating, when the app id and title of the window matches a rule.
		There is also a "scratchpad" command that can be run on a key-bind.`,
	Version: getVersionInfo(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
//
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// getVersionInfo returns the version of the executable.
func getVersionInfo() string {
	if config.Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					config.CommitSHA = setting.Value[:8]
				}
				if setting.Key == "vcs.time" {
					config.BuildDate = setting.Value
				}
			}
			config.Version = info.Main.Version
			if config.Version == "(devel)" {
				return config.Version
			}
		}
	}

	if config.CommitSHA == "unknown" || config.BuildDate == "unknown" {
		return config.Version
	}

	return fmt.Sprintf("%s (commit: %s, built at: %s)", config.Version, config.CommitSHA, config.BuildDate)
}
