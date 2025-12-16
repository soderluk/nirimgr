package scratchpad

import (
	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/internal/connection"
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move the currently focused window to the scratchpad",
	Long:  `Moves the currently focused window to the scratchpad workspace. This requires niri to have a named workspace called "scratchpad". See README.md for more information.`,
	Run: func(cmd *cobra.Command, args []string) {
		moveToScratchpad()
	},
}

func init() {
	ScratchCmd.AddCommand(moveCmd)
}

// moveToScratchpad moves the currently focused window to the scratchpad workspace.
//
// This requires niri to have a named workspace called "scratchpad". See README.md for more information.
func moveToScratchpad() {
	// Move to scratchpad:
	// 1. get scratch workspace
	// 2. get focused window
	// 3. move window to floating
	// 4. move floating window to scratchpad workspace, focus=false
	scratchpad, _ := getWorkspace(config.Config.ScratchpadWorkspace)
	focusedWindow, _ := getFocusedWindow()

	actionList := []actions.Action{
		actions.MoveWindowToWorkspace{
			AName:    actions.AName{Name: "MoveWindowToWorkspace"},
			WindowID: focusedWindow.ID,
			Reference: actions.WorkspaceReferenceArg{
				ID: scratchpad.ID,
			},
			Focus: false,
		},
		actions.MoveWindowToFloating{
			AName: actions.AName{Name: "MoveWindowToFloating"},
			ID:    focusedWindow.ID,
		},
	}

	for _, action := range actionList {
		connection.PerformAction(action)
	}
}
