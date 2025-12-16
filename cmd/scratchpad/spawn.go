package scratchpad

import (
	"log/slog"
	"os"
	"strings"

	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/internal/connection"
	"github.com/soderluk/nirimgr/models"
	"github.com/spf13/cobra"
)

var spawnCmd = &cobra.Command{
	Use:   "spawn-or-focus [app-id]",
	Short: "Spawn an app or focus it if already running",
	Long:  `Spawns the specified app or focuses it if it's already running. Requires configuration for the commands and app IDs.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		spawnOrFocus(args[0])
	},
}

func init() {
	ScratchCmd.AddCommand(spawnCmd)
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
