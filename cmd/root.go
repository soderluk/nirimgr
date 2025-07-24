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
//
//	Usage: nirimgr scratch [move|show]
//
// # Version
//
// The version command prints out the version and build info of nirimgr.
//
//	Usage: nirimgr version
package cmd

import (
	"os"

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
