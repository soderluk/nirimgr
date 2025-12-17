package scratchpad

import (
	"github.com/soderluk/nirimgr/cmd"
	"github.com/spf13/cobra"
)

// ScratchCmd is the parent command for scratchpad operations.
var ScratchCmd = &cobra.Command{
	Use:   "scratch",
	Short: "Simple support for a scratchpad in Niri",
	Long: `An i3wm inspired simple scratchpad functionality for Niri.
		nirimgr scratch move - moves the currently focused window to the scratchpad workspace.
		nirimgr scratch show - moves the last window in the scratchpad workspace to the currently focused workspace.
		nirimgr scratch spawn-or-focus app-id - Spawns the specified app or focuses it if it's already running. Requires configuration for the commands and app IDs.`,
}

func init() {
	cmd.RootCmd.AddCommand(ScratchCmd)
}
