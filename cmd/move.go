package cmd

import (
	"errors"
	"strconv"

	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/internal/common"
	"github.com/soderluk/nirimgr/internal/connection"
	"github.com/soderluk/nirimgr/models"
	"github.com/spf13/cobra"
)

// moveCmd moves a floating window to the left/right/top/bottom edges.
//
// This is from https://github.com/YaLTeR/niri/discussions/1656#discussioncomment-14268880
var moveCmd = &cobra.Command{
	Use:          "move",
	Short:        "Moves a floating window to the left/right/top/bottom edges of the screen.",
	SilenceUsage: true, // If there's an error during running the command, don't show usage.
	Args:         cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		direction := args[0]
		// Note: The border must be at minimum 1, because if it's 0, the "newX" will be empty.
		border := float64(1)
		if len(args) > 1 {
			border, _ = strconv.ParseFloat(args[1], 64)
		}

		windows, err := connection.ListWindows()
		if err != nil {
			return errors.New("could not get windows")
		}
		window, err := common.FilterWindowsChain(windows, func(w *models.Window) bool {
			return w.IsFocused && w.IsFloating
		}).First()
		if err != nil {
			return errors.New("no active floating window")
		}

		workspaces, err := connection.ListWorkspaces()
		if err != nil {
			return errors.New("could not get workspaces")
		}

		workspace, err := common.FilterWorkspacesChain(workspaces, func(w *models.Workspace) bool {
			return w.IsFocused
		}).First()
		if err != nil {
			return errors.New("could not get focused workspace")
		}
		outputName := workspace.Output
		outputs, err := connection.ListOutputs()
		if err != nil {
			return errors.New("could not get outputs")
		}
		output, err := common.FilterOutputsChain(outputs, func(o *models.Output) bool {
			return o.Name == outputName
		}).First()
		if err != nil {
			return errors.New("could not get output")
		}

		width := float64(window.Layout.WindowSize[0])
		height := float64(window.Layout.WindowSize[1])
		x := float64(window.Layout.TilePosInWorkspaceView[0])
		y := float64(window.Layout.TilePosInWorkspaceView[1])

		screenWidth := float64(output.Logical.Width)
		screenHeight := float64(output.Logical.Height)

		var newX, newY float64
		switch direction {
		case "left":
			newX = border
			newY = y - 34
		case "right":
			newX = (screenWidth - width - border)
			newY = y - 34
		case "up":
			newX = x
			newY = border
		case "down":
			newX = x
			newY = (screenHeight - height - border) - 34
		default:
			return errors.New("invalid direction provided")
		}

		data := map[string]any{
			"ID": window.ID,
			"x": map[string]float64{
				"SetFixed": newX,
			},
			"y": map[string]float64{
				"SetFixed": newY,
			},
		}
		action, _ := actions.FromName("MoveFloatingWindow", data)
		connection.PerformAction(action)
		return nil
	},
}

func init() {
	floatingCmd.AddCommand(moveCmd)
}
