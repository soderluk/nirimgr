// Package events handles events from Niri event-stream.
//
// Listens to the event stream on the NIRI_SOCKET and reacts to
// events in a specified manner.
package events

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/internal/common"
	"github.com/soderluk/nirimgr/internal/connection"
	"github.com/soderluk/nirimgr/models"
)

// Run starts listening on the event stream, and handle the events.
//
// Currently we only handle a few events, `WindowsChanged`, `WindowOpenedOrChanged` and `WindowClosed`.
// If you need to handle more events, add them in the switch statement.
// Initially the thought was to support the "Dynamic open-float script, for Bitwarden and other windows that set title/app-id late":
// https://github.com/YaLTeR/niri/discussions/1599
// But it doesn't stop us from handling other types of events as well.
func Run() {
	events, err := EventStream()
	if err != nil {
		slog.Error("Could not get events", "error", err.Error())
		panic(err)
	}
	existingWindows := make(map[uint64]*models.Window)
	existingWorkspaces := make(map[uint64]*models.Workspace)

	for event := range events {
		switch ev := event.(type) {
		case *WindowsChanged:
			slog.Debug("Handling event", "name", common.Repr(ev))
			for _, win := range ev.Windows {
				matchWindowAndPerformActions(win, existingWindows)
				existingWindows[win.ID] = win
			}
		case *WindowOpenedOrChanged:
			slog.Debug("Handling event", "name", common.Repr(ev))
			matchWindowAndPerformActions(ev.Window, existingWindows)
			existingWindows[ev.Window.ID] = ev.Window
		case *WindowClosed:
			slog.Debug("Handling event", "name", common.Repr(ev))
			delete(existingWindows, ev.ID)
		case *WorkspacesChanged:
			slog.Debug("Handling event", "name", common.Repr(ev))
			// Remove workspaces that are no longer present
			newWorkspaceIDs := make(map[uint64]struct{})
			for _, workspace := range ev.Workspaces {
				newWorkspaceIDs[workspace.ID] = struct{}{}
			}
			for id := range existingWorkspaces {
				if _, found := newWorkspaceIDs[id]; !found {
					slog.Debug("Removing workspace from existing workspaces", "id", id)
					delete(existingWorkspaces, id)
				}
			}
			for _, workspace := range ev.Workspaces {
				matchWorkspaceAndPerformActions(workspace, existingWorkspaces)
				existingWorkspaces[workspace.ID] = workspace
			}
		default:
		}
	}
}

// EventStream listens on the events in Niri event-stream.
//
// The function will use a goroutine to return the event models.
// Inspiration from: https://github.com/probeldev/niri-float-sticky
func EventStream() (<-chan any, error) {
	stream := make(chan any)
	socket := connection.Socket()

	go func() {
		defer connection.PutSocket(socket)
		defer socket.Close()
		defer close(stream)

		for line := range socket.Recv() {
			if len(line) < 2 {
				continue
			}

			var event map[string]json.RawMessage

			if err := json.Unmarshal(line, &event); err != nil {
				slog.Error("Error decoding JSON", "error", err.Error())
				continue
			}
			_, model, err := ParseEvent(event)
			if err != nil {
				slog.Error("Could not parse event!", "event", event, "error", err.Error())
			}
			stream <- model
		}
	}()

	if err := socket.Send(fmt.Sprintf("\"%s\"", models.EventStream)); err != nil {
		return nil, fmt.Errorf("error requesting event stream: %w", err)
	}

	return stream, nil
}

// ParseEvent parses the given event into it's struct.
//
// Returns the name, model, error. The name is the name of the event, model is the populated struct.
func ParseEvent(event map[string]json.RawMessage) (string, any, error) {
	for name, raw := range event {
		model := FromRegistry(name, raw)
		if model == nil {
			continue
		}
		slog.Debug("Parsed event into", "model", common.Repr(model))
		return name, model, nil
	}
	return "", nil, fmt.Errorf("no event found")
}

// matchWindowAndPerformActions updates the window struct if it matches the rule as configured in the config file.
//
// If the matching window has any defined actions in the config, run them sequentially on the matched window.
// The functionality is taken from the "Dynamic open-float script, for Bitwarden and other windows that set title/app-id late":
// https://github.com/YaLTeR/niri/discussions/1599
func matchWindowAndPerformActions(window *models.Window, existingWindows map[uint64]*models.Window) {
	window.Matched = false
	if existing, ok := existingWindows[window.ID]; ok {
		window.Matched = existing.Matched
	}

	matchedBefore := window.Matched
	window.Matched = false
	var rawActions map[string]json.RawMessage
	for _, r := range config.Config.GetRules() {
		if r.Type != "window" && r.Type != "" {
			continue
		}
		if r.WindowMatches(*window) {
			window.Matched = true
			if len(r.Actions) > 0 {
				rawActions = r.Actions
			}
			break
		}
	}
	actionList := actions.ParseRawActions(rawActions)
	if window.Matched && !matchedBefore {
		for _, a := range actionList {
			// Set the action Id dynamically here.
			a = actions.SetActionID(a, window.ID)
			connection.PerformAction(a)
		}
	}
}

// matchWorkspaceAndPerformActions updates the workspace struct if it matches the rule as configured in the config file.
//
// If the matching workspace has any defined actions in the config, run them sequentially on the matched workspace.
func matchWorkspaceAndPerformActions(workspace *models.Workspace, existingWorkspaces map[uint64]*models.Workspace) {
	workspace.Matched = false
	if existing, ok := existingWorkspaces[workspace.ID]; ok {
		workspace.Matched = existing.Matched
	}
	matchedBefore := workspace.Matched

	workspace.Matched = false
	var rawActions map[string]json.RawMessage
	for _, r := range config.Config.GetRules() {
		if r.Type != "workspace" {
			continue
		}
		if r.WorkspaceMatches(*workspace) {
			workspace.Matched = true
			if len(r.Actions) > 0 {
				rawActions = r.Actions
			}
			break
		}
	}
	actionList := actions.ParseRawActions(rawActions)
	if workspace.Matched && !matchedBefore {
		for _, a := range actionList {
			a = actions.SetActionID(a, workspace.ID)
			connection.PerformAction(a)
		}
	}
}

// FromRegistry returns the populated model from the EventRegistry by given name.
func FromRegistry(name string, data []byte) Event {
	model, ok := EventRegistry[name]
	if !ok {
		slog.Error("Could not get event model for event", "name", name)
		return nil
	}
	event := model()
	if err := json.Unmarshal(data, event); err != nil {
		slog.Error("Could not unmarshal event", "name", name, "error", err.Error())
		return nil
	}
	return event
}

// EventRegistry contains all the events Niri currently sends.
//
// The key needs to be the event name, and it should return the correct event model, and set
// its EName embedded struct. If you know of a better way to handle this, please let me know.
var EventRegistry = map[string]func() Event{
	"WorkspacesChanged":            func() Event { return &WorkspacesChanged{EName: EName{Name: "WorkspacesChanged"}} },
	"WorkspaceUrgencyChanged":      func() Event { return &WorkspaceUrgencyChanged{EName: EName{Name: "WorkspaceUrgencyChanged"}} },
	"WorkspaceActivated":           func() Event { return &WorkspaceActivated{EName: EName{Name: "WorkspaceActivated"}} },
	"WorkspaceActiveWindowChanged": func() Event { return &WorkspaceActiveWindowChanged{EName: EName{Name: "WorkspaceActiveWindowChanged"}} },
	"WindowsChanged":               func() Event { return &WindowsChanged{EName: EName{Name: "WindowsChanged"}} },
	"WindowOpenedOrChanged":        func() Event { return &WindowOpenedOrChanged{EName: EName{Name: "WindowOpenedOrChanged"}} },
	"WindowClosed":                 func() Event { return &WindowClosed{EName: EName{Name: "WindowClosed"}} },
	"WindowFocusChanged":           func() Event { return &WindowFocusChanged{EName: EName{Name: "WindowFocusChanged"}} },
	"WindowUrgencyChanged":         func() Event { return &WindowUrgencyChanged{EName: EName{Name: "WindowUrgencyChanged"}} },
	"KeyboardLayoutsChanged":       func() Event { return &KeyboardLayoutsChanged{EName: EName{Name: "KeyboardLayoutsChanged"}} },
	"KeyboardLayoutSwitched":       func() Event { return &KeyboardLayoutSwitched{EName: EName{Name: "KeyboardLayoutSwitched"}} },
	"OverviewOpenedOrClosed":       func() Event { return &OverviewOpenedOrClosed{EName: EName{Name: "OverviewOpenedOrClosed"}} },
}
