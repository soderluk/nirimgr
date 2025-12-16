package scratchpad

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/internal/common"
	"github.com/soderluk/nirimgr/internal/connection"
	"github.com/soderluk/nirimgr/models"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:          "show",
	Short:        "Show a window from the scratchpad",
	Long:         `Moves the last window from the scratchpad workspace to the currently active workspace. This requires niri to have a named workspace called "scratchpad" (configurable in the config.json). See README.md for more information.`,
	SilenceUsage: true, // If there's an error during running the command, don't show usage.
	RunE: func(cmd *cobra.Command, args []string) error {
		return showScratchpad()
	},
}

func init() {
	ScratchCmd.AddCommand(showCmd)
}

// showScratchpad moves the last window from the scratchpad workspace to the currently active workspace.
//
// This requires niri to have a named workspace called "scratchpad". See README.md for more information.
func showScratchpad() error {
	// Show scratchpad:
	// 1. get current workspace
	// 2. list all windows in scratchpad workspace
	// 3. take latest window and move it to the current workspace

	scratchpad, _ := getWorkspace(config.Config.ScratchpadWorkspace)
	focusedWorkspace, _ := getWorkspace("focused")

	windows, err := connection.ListWindows()
	if err != nil {
		slog.Error("Could not unmarshal windows", "error", err.Error())
		return errors.New("could not unmarshal windows")
	}
	// Filter the scratchpad windows.
	workspaceWindows := filterWindows(windows, func(w *models.Window) bool {
		return w.WorkspaceID == scratchpad.ID
	})

	var window *models.Window
	// Return the last window in the list, if we have any.
	if len(workspaceWindows) > 0 {
		// If we have more than one window in the scratchpad, open a launcher to select the window.
		if len(workspaceWindows) > 1 {
			// Build the list of windows to pass to the launcher.
			command := ""
			for idx, window := range workspaceWindows {
				command += fmt.Sprintf("%d - %s\n", idx, window.Title)
			}
			// Build the full command to run.
			fullCommand := fmt.Sprintf("echo \"%s\" | %s %s | awk '{print $1}'", command, config.Config.Launcher, config.Config.LauncherOptions)
			result, err := common.RunCommand(fullCommand)
			if err != nil {
				slog.Error("Error running command", slog.String("command", fullCommand), slog.Any("error", err.Error()))
				return errors.New("could not run command")
			}
			idx, err := strconv.ParseUint(strings.Replace(string(result), "\n", "", 1), 10, 64)
			if err != nil {
				slog.Error("Could not convert result to uint64", slog.String("result", string(result)), slog.Any("error", err.Error()))
				return errors.New("could not run command")
			}
			window = workspaceWindows[idx]
		} else {
			window = workspaceWindows[len(workspaceWindows)-1]
		}

		actionList := []actions.Action{
			actions.MoveWindowToWorkspace{
				AName:    actions.AName{Name: "MoveWindowToWorkspace"},
				WindowID: window.ID,
				Reference: actions.WorkspaceReferenceArg{
					ID: focusedWorkspace.ID,
				},
				Focus: true,
			},
			// Manually focus the window, since the `Focus: true` does nothing in the above action.
			actions.FocusWindow{
				AName: actions.AName{Name: "FocusWindow"},
				ID:    window.ID,
			},
		}

		// If we have actions, append them to the list.
		if len(config.Config.ShowScratchpadActions) > 0 {
			for _, action := range actions.ParseRawActions(config.Config.ShowScratchpadActions) {
				a := actions.HandleDynamicIDs(action, models.PossibleKeys{
					ID:       window.ID,
					WindowID: window.ID,
				})
				actionList = append(actionList, a)
			}
		}

		for _, action := range actionList {
			connection.PerformAction(action)
		}
	}
	return nil
}
