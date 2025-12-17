package scratchpad

import (
	"encoding/json"
	"fmt"

	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/internal/connection"
	"github.com/soderluk/nirimgr/models"
)

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
