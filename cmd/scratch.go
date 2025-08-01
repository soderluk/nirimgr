package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/internal/connection"
	"github.com/soderluk/nirimgr/models"

	"github.com/spf13/cobra"
)

// scratchCmd is the command for handling the scratchpad.
//
// Depending on the arguments, we either show the scratchpad window,
// or move the currently focused window to scratchpad.
var scratchCmd = &cobra.Command{
	Use:   "scratch",
	Short: "Simple support for a scratchpad in Niri",
	Long: `An i3wm inspired simple scratchpad functionality for Niri.
		nirimgr scratch move - moves the currently focused window to the scratchpad workspace.
		nirimgr scratch show - moves the last window in the scratchpad workspace to the currently focused workspace.
		nirimgr scratch spawn-or-focus app-id - Spawns the specified app or focuses it if it's already running. Requires configuration for the commands and app IDs.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]
		switch command {
		case "move":
			moveToScratchpad()
		case "show":
			showScratchpad()
		case "spawn-or-focus":
			spawnOrFocus(args[1])
		default:
			slog.Error("Unknown command", "cmd", command)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(scratchCmd)
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

// showScratchpad moves the last window from the scratchpad workspace to the currently active workspace.
//
// This requires niri to have a named workspace called "scratchpad". See README.md for more information.
func showScratchpad() {
	// Show scratchpad:
	// 1. get current workspace
	// 2. list all windows in scratchpad workspace
	// 3. take latest window and move it to the current workspace

	scratchpad, _ := getWorkspace(config.Config.ScratchpadWorkspace)
	focusedWorkspace, _ := getWorkspace("focused")

	windows, err := connection.ListWindows()
	if err != nil {
		slog.Error("Could not unmarshal windows", "error", err.Error())
		os.Exit(1)
	}
	// Filter the scratchpad windows.
	workspaceWindows := filterWindows(windows, func(w *models.Window) bool {
		return w.WorkspaceID == scratchpad.ID
	})

	// Return the last window in the list, if we have any.
	if len(workspaceWindows) > 0 {
		window := workspaceWindows[len(workspaceWindows)-1]
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
}

// spawnOrFocus will spawn a specified window or focus it if it's already open.
//
// This is heavily inspired by the discussion over here: https://github.com/YaLTeR/niri/discussions/329#discussioncomment-13378697
// I adapted the functionality to be supported in nirimgr. The commands and app id's are configurable in the config.json.
func spawnOrFocus(arg string) {
	windows, err := connection.ListWindows()
	if err != nil {
		slog.Error("Could not list windows", "error", err.Error())
		os.Exit(1)
	}

	var matchedWindow *models.Window
	command, err := config.Config.SpawnOrFocus.Command(arg)
	if err != nil {
		slog.Error("Could not get command", "error", err.Error())
		os.Exit(1)
	}
	for _, window := range windows {
		for _, rule := range config.Config.SpawnOrFocus.Rules {
			if rule.WindowMatches(*window) {
				slog.Debug("Trying window", "appId", window.AppID)
				// If the app id doesn't contain the given argument, skip this window.
				if !strings.Contains(window.AppID, arg) {
					continue
				}
				slog.Debug("Window matches rule", "window", window.AppID, "rule", rule)
				matchedWindow = window
				break
			}
		}
	}
	if matchedWindow != nil {
		slog.Debug("matched window", "window", matchedWindow.Title)
		if matchedWindow.IsFocused {
			connection.PerformAction(actions.FocusWindowPrevious{AName: actions.AName{Name: "FocusWindowPrevious"}})
		} else {
			connection.PerformAction(actions.FocusWindow{AName: actions.AName{Name: "FocusWindow"}, ID: matchedWindow.ID})
		}
	} else {
		slog.Debug("Didn't match any window, spawning command", "cmd", command)
		connection.PerformAction(actions.Spawn{AName: actions.AName{Name: "Spawn"}, Command: command})
	}
}

// getWorkspace returns a workspace.
//
// Given the wtype, we return either the named workspace, focused or active workspace.
func getWorkspace(wtype string) (*models.Workspace, error) {
	workspaces, err := connection.ListWorkspaces()
	if err != nil {
		return nil, err
	}

	// If the wtype is not given, default to "scratchpad"
	if wtype == "" {
		wtype = "scratchpad"
	}
	// If the scratchpad workspace is not configured, default to "scratchpad"
	scratchpadWorkspace := config.Config.ScratchpadWorkspace
	if scratchpadWorkspace == "" {
		scratchpadWorkspace = "scratchpad"
	}
	for _, workspace := range workspaces {
		switch wtype {
		case scratchpadWorkspace:
			if workspace.Name == scratchpadWorkspace {
				return workspace, nil
			}
		case "focused":
			if workspace.IsFocused {
				return workspace, nil
			}
		case "active":
			if workspace.IsActive {
				return workspace, nil
			}
		default:
			return nil, fmt.Errorf("invalid type '%v'", wtype)
		}
	}
	return nil, fmt.Errorf("no workspace found matching '%v'", wtype)
}

// getFocusedWindow returns the currently focused window.
func getFocusedWindow() (*models.Window, error) {
	response, err := connection.PerformRequest(models.FocusedWindow)
	if err != nil {
		return nil, err
	}

	resp := <-response

	var window *models.Window
	if err := json.Unmarshal(resp.Ok["FocusedWindow"], &window); err != nil {
		return nil, err
	}
	return window, nil
}

// filterWindows returns a slice of window models depending on the filtering function.
func filterWindows(data []*models.Window, f func(*models.Window) bool) []*models.Window {
	w := make([]*models.Window, 0)
	for _, e := range data {
		if f(e) {
			w = append(w, e)
		}
	}

	return w
}
