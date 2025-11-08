package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/internal/common"
	"github.com/soderluk/nirimgr/internal/connection"
	"github.com/spf13/cobra"
)

// switchWindowCmd is used to open fuzzel with a list of windows, and you can choose which window to focus on.
var switchWindowCmd = &cobra.Command{
	Use:          "switch-window",
	Short:        "Switch window focus with fuzzel.",
	Long:         `Opens fuzzel with a list of currently open windows. Choose one and you focus that window.`,
	SilenceUsage: true, // If there's an error during running the command, don't show usage.
	RunE: func(cmd *cobra.Command, args []string) error {
		windows, err := connection.ListWindows()
		if err != nil {
			return errors.New("could not get windows")
		}
		var command string
		for _, window := range windows {
			command += fmt.Sprintf("[WS: %d] [%s] (%s) -> id: %d\n", window.WorkspaceID, window.AppID, window.Title, window.ID)
		}
		// TODO: Support any dmenu launcher here. Make the launcher configurable.
		fullCommand := fmt.Sprintf("echo \"%s\" | sort | fuzzel -d -w 75 | awk '{print $NF}'", command)
		result, err := common.RunCommand(fullCommand)
		if err != nil {
			slog.Error("Error running command", slog.String("command", fullCommand), slog.Any("error", err.Error()))
			return errors.New("could not run command")
		}
		windowID, err := strconv.ParseUint(strings.Replace(string(result), "\n", "", 1), 10, 64)
		if err != nil {
			slog.Error("Could not convert result to uint64", slog.String("result", string(result)), slog.Any("error", err.Error()))
			return errors.New("could not run command")
		}
		action, _ := actions.FromName("FocusWindow", map[string]any{"ID": windowID})
		connection.PerformAction(action)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchWindowCmd)
}
