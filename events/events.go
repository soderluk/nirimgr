// Package events handles events from Niri event-stream.
//
// Listens to the event stream on the NIRI_SOCKET and reacts to
// events in a specified manner.
package events

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/expr-lang/expr"
	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/internal/common"
	"github.com/soderluk/nirimgr/internal/connection"
	"github.com/soderluk/nirimgr/models"
)

// Run starts listening on the event stream, and handle the events.
//
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

	// Any events we want to specifically listen to and perform actions on the event Window/Workspace/whatever.
	listenToEvents := config.Config.Events

	for event := range events {
		// These events are specific for the matching logic of nirimgr.
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
			// Any events we're not specifically listening to, let's check if there are any configured events.
			if ev != nil {
				// Handle the event if it exists in the map
				if actionConfigs, exists := listenToEvents[ev.GetName()]; exists {
					for actionName, actionConfig := range actionConfigs {
						rawAction := map[string]json.RawMessage{
							actionName: actionConfig.Params,
						}
						// Perform each defined action on the event.
						for _, a := range ActionsFromRaw(rawAction) {
							evaluationResult, err := EvaluateCondition(actionConfig.When, ev)
							if err != nil {
								slog.Error("Error in EvaluateCondition", slog.Any("error", err))
							}
							if evaluationResult {
								possibleKeys := ev.GetPossibleKeys()
								a = actions.HandleDynamicIDs(a, possibleKeys)
								connection.PerformAction(a)
							} else {
								slog.Debug(
									"Not performing action",
									slog.String("name", actionName),
									slog.Bool("EvaluateCondition", evaluationResult),
								)
							}
						}
					}
				}
			}
		}
	}
}

// EventStream listens on the events in Niri event-stream.
//
// The function will use a goroutine to return the event models.
// Inspiration from: https://github.com/probeldev/niri-float-sticky
func EventStream() (<-chan Event, error) {
	stream := make(chan Event)
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
func ParseEvent(event map[string]json.RawMessage) (string, Event, error) {
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
	var actionConfigs map[string]models.ActionConfig
	for _, r := range config.Config.GetRules() {
		if r.Type != "window" && r.Type != "" {
			continue
		}
		if r.WindowMatches(*window) {
			window.Matched = true
			if len(r.Actions) > 0 {
				actionConfigs = r.Actions
			}
			break
		}
	}
	if window.Matched && !matchedBefore {
		for actionName, actionConfig := range actionConfigs {
			rawAction := map[string]json.RawMessage{
				actionName: actionConfig.Params,
			}
			for _, a := range ActionsFromRaw(rawAction) {
				// If we have a condition defined, evaluate it, and perform the action if it evaluates to true.
				evaluationResult, err := EvaluateCondition(actionConfig.When, window)
				if err != nil {
					slog.Error("Error in EvaluateCondition", slog.Any("error", err))
				}
				if evaluationResult {
					a = actions.HandleDynamicIDs(a, models.PossibleKeys{
						ID:       window.ID,
						WindowID: window.ID,
					})
					connection.PerformAction(a)
				} else {
					slog.Debug("Not doing action", slog.String("name", actionName), slog.Bool("EvaluateCondition", evaluationResult))
				}
			}
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
	var actionConfigs map[string]models.ActionConfig
	for _, r := range config.Config.GetRules() {
		if r.Type != "workspace" {
			continue
		}
		if r.WorkspaceMatches(*workspace) {
			workspace.Matched = true
			if len(r.Actions) > 0 {
				actionConfigs = r.Actions
			}
			break
		}
	}
	if workspace.Matched && !matchedBefore {
		for actionName, actionConfig := range actionConfigs {
			rawAction := map[string]json.RawMessage{
				actionName: actionConfig.Params,
			}
			for _, a := range ActionsFromRaw(rawAction) {
				// If we have a condition defined, evaluate it, and perform the action if it evaluates to true.
				evaluationResult, err := EvaluateCondition(actionConfig.When, workspace)
				if err != nil {
					slog.Error("Error in EvaluateCondition", slog.Any("error", err))
				}
				if evaluationResult {
					a = actions.HandleDynamicIDs(a, models.PossibleKeys{
						ID:             workspace.ID,
						ActiveWindowID: workspace.ActiveWindowID,
						Reference: models.ReferenceKeys{
							ID:    workspace.ID,
							Index: workspace.Idx,
							Name:  workspace.Name,
						},
					})
					connection.PerformAction(a)
				} else {
					slog.Debug(
						"Not performing action",
						slog.String("name", actionName),
						slog.Bool("EvaluateCondition", evaluationResult),
					)
				}
			}
		}
	}
}

// ActionsFromRaw converts the raw actions from the config into a list of Action structs.
func ActionsFromRaw(rawActions map[string]json.RawMessage) []actions.Action {
	return actions.ParseRawActions(rawActions)
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

// EvaluateCondition evaluates the given condition on the given model.
//
// If the model is a "WindowUrgencyChanged" event, we know that it has a field called Urgent, so
// the condition could be "model.Urgent == true" to run an action only when the event urgency is set.
// Note: The model can be an event, action, window, workspace or any other model.
func EvaluateCondition(condition string, model any) (bool, error) {
	slog.Debug("EvaluateCondition", slog.String("condition", condition))
	// We always evaluate empty conditions to true.
	if condition == "" {
		return true, nil
	}

	env := map[string]any{
		"model": model,
	}
	slog.Debug("EvaluateCondition", slog.Any("model", model))
	program, err := expr.Compile(condition, expr.Env(env))
	if err != nil {
		return false, fmt.Errorf("invalid condition '%s': %w", condition, err)
	}
	result, err := expr.Run(program, env)
	if err != nil {
		return false, fmt.Errorf("error evaluating condition '%s': %w", condition, err)
	}

	boolResult, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("condition '%s' didn't evaluate to a boolean", condition)
	}

	return boolResult, nil
}

// EventRegistry contains all the events Niri currently sends.
//
// The key needs to be the event name, and it should return the correct event model, and set
// its EName embedded struct. If you know of a better way to handle this, please let me know.
var EventRegistry = map[string]func() Event{
	"ConfigLoaded":                 func() Event { return &ConfigLoaded{EName: EName{Name: "ConfigLoaded"}} },
	"KeyboardLayoutSwitched":       func() Event { return &KeyboardLayoutSwitched{EName: EName{Name: "KeyboardLayoutSwitched"}} },
	"KeyboardLayoutsChanged":       func() Event { return &KeyboardLayoutsChanged{EName: EName{Name: "KeyboardLayoutsChanged"}} },
	"OverviewOpenedOrClosed":       func() Event { return &OverviewOpenedOrClosed{EName: EName{Name: "OverviewOpenedOrClosed"}} },
	"WindowClosed":                 func() Event { return &WindowClosed{EName: EName{Name: "WindowClosed"}} },
	"WindowFocusChanged":           func() Event { return &WindowFocusChanged{EName: EName{Name: "WindowFocusChanged"}} },
	"WindowLayoutsChanged":         func() Event { return &WindowLayoutsChanged{EName: EName{Name: "WindowLayoutsChanged"}} },
	"WindowOpenedOrChanged":        func() Event { return &WindowOpenedOrChanged{EName: EName{Name: "WindowOpenedOrChanged"}} },
	"WindowUrgencyChanged":         func() Event { return &WindowUrgencyChanged{EName: EName{Name: "WindowUrgencyChanged"}} },
	"WindowsChanged":               func() Event { return &WindowsChanged{EName: EName{Name: "WindowsChanged"}} },
	"WorkspaceActivated":           func() Event { return &WorkspaceActivated{EName: EName{Name: "WorkspaceActivated"}} },
	"WorkspaceActiveWindowChanged": func() Event { return &WorkspaceActiveWindowChanged{EName: EName{Name: "WorkspaceActiveWindowChanged"}} },
	"WorkspaceUrgencyChanged":      func() Event { return &WorkspaceUrgencyChanged{EName: EName{Name: "WorkspaceUrgencyChanged"}} },
	"WorkspacesChanged":            func() Event { return &WorkspacesChanged{EName: EName{Name: "WorkspacesChanged"}} },
}
