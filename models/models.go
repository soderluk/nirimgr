// Package models contains all the necessary models for different niri objects.
package models

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
)

// NiriRequest is the representation of a simple niri request.
//
// The request is sent to the socket as a string, e.g. "Windows" returns all the current windows.
type NiriRequest string

const (
	// Outputs lists connected outputs
	Outputs NiriRequest = "Outputs"
	// Workspaces lists workspaces
	Workspaces NiriRequest = "Workspaces"
	// Windows lists open windows
	Windows NiriRequest = "Windows"
	// Layers lists open layer-shell surfaces
	Layers NiriRequest = "Layers"
	// ListKeyboardLayouts lists the configured keyboard layouts
	ListKeyboardLayouts NiriRequest = "KeyboardLayouts"
	// FocusedOutput prints information about the focused output
	FocusedOutput NiriRequest = "FocusedOutput"
	// FocusedWindow prints information about the focused window
	FocusedWindow NiriRequest = "FocusedWindow"
	// PickWindow to pick a window with the mouse and print information about it. Not applicable to nirimgr.
	PickWindow NiriRequest = "PickWindow"
	// PickColor to pick a color from the screen with the mouse. Not applicable to nirimgr.
	PickColor NiriRequest = "PickColor"
	// RunAction performs an action
	RunAction NiriRequest = "Action"
	// ChangeOutput changes output configuration temporarily
	ChangeOutput NiriRequest = "Output"
	// EventStream starts continuously receiving events from the compositor
	EventStream NiriRequest = "EventStream"
	// Version prints the version of the running niri instance
	Version NiriRequest = "Version"
	// RequestError requests an error from the running niri instance
	RequestError NiriRequest = "RequestError"
	// OverviewState prints the overview state
	OverviewState NiriRequest = "OverviewState"
)

// Config contains the configuration for nirimgr.
type Config struct {
	// LogLevel is the log level to use. One of "DEBUG", "INFO", "WARN", "ERROR" should be used. Defaults to "INFO".
	LogLevel string `json:"logLevel"`
	// Rules contains the rules to match windows, and the actions to perform on them.
	Rules []Rule `json:"rules,omitempty"`
	// ScratchpadWorkspace is the name of the scratchpad workspace. Defaults to "scratchpad".
	//
	// NOTE: The named workspace must be defined in niri config.
	ScratchpadWorkspace string `json:"scratchpadWorkspace,omitempty"`
	// SpawnOrFocus defines the configuration for the spawn-or-focus command.
	SpawnOrFocus SpawnOrFocus `json:"spawnOrFocus,omitempty"`
	// ShowScratchpadActions lists actions that should be performed on the shown scratchpad window.
	//
	// The `scratch show` command will always run MoveWindowToWorkspace and FocusWindow, but in addition can perform the following actions,
	// e.g. if you want to center the window or resize it or something.
	ShowScratchpadActions map[string]json.RawMessage `json:"showScratchpadActions,omitempty"`
}

// GetRules returns the configured rules.
//
// NOTE: We cannot use the name Rules() because we already define the Rules in the struct.
func (c *Config) GetRules() []Rule {
	var rules []Rule
	rules = append(rules, c.Rules...)
	return rules
}

// SpawnOrFocus defines the rules and commands to run for the spawn-or-focus command.
type SpawnOrFocus struct {
	Rules []Rule `json:"rules,omitempty"`
	// Command is the command to spawn for the spawnOrFocus command.
	Commands map[string][]string `json:"commands,omitempty"`
}

// Command returns the specified command to run for the given key.
func (s *SpawnOrFocus) Command(key string) ([]string, error) {
	command, ok := s.Commands[key]
	if !ok {
		return nil, fmt.Errorf("could not read command for %s", key)
	}
	return command, nil
}

// Match is used to match a window.
type Match struct {
	// Title matches the title of the window. Used only for rules with type "window".
	Title string `json:"title,omitempty"`
	// AppID matches the app-id of the window. Used only for rules with type "window".
	AppID string `json:"appId,omitempty"`
	// Name is used for rule types "workspace" to match on the workspace name.
	Name string `json:"name,omitempty"`
	// Output is used for rule types "workspace" to match on the workspace output name.
	Output string `json:"output,omitempty"`
}

// WindowMatches checks if the window matches the specified rule match.
func (m Match) WindowMatches(window Window) bool {
	if m.Title == "" && m.AppID == "" {
		slog.Debug("Title and AppID empty for window", "window", window.ID)
		return false
	}
	matched := true

	if m.Title != "" {
		titleMatch, err := regexp.MatchString(m.Title, window.Title)
		if err != nil {
			slog.Error("Could not match title", "error", err.Error())
			return false
		}
		matched = matched && titleMatch
	}
	if m.AppID != "" {
		appMatch, err := regexp.MatchString(m.AppID, window.AppID)
		if err != nil {
			slog.Error("Could not match AppID", "error", err.Error())
			return false
		}
		matched = matched && appMatch
	}

	return matched
}

// WorkspaceMatches checks if the workspace matches the specified rule match.
func (m Match) WorkspaceMatches(workspace Workspace) bool {
	if m.Name == "" && m.Output == "" {
		slog.Debug("Name and Output empty for workspace", "workspace", workspace.ID)
		return false
	}
	matched := true

	if m.Name != "" {
		titleMatch, err := regexp.MatchString(m.Name, workspace.Name)
		if err != nil {
			slog.Error("Could not match Name", "error", err.Error())
			return false
		}
		matched = matched && titleMatch
	}
	if m.Output != "" {
		appMatch, err := regexp.MatchString(m.Output, workspace.Output)
		if err != nil {
			slog.Error("Could not match Output", "error", err.Error())
			return false
		}
		matched = matched && appMatch
	}

	return matched
}

// Rule contains the matches, excludes and actions for a window.
type Rule struct {
	// Type is the type of object we want to match, e.g. window or workspace. Defaults to window.
	Type string `json:"type"`
	// Match list of matches to target a window.
	Match []Match `json:"match,omitempty"`
	// Exclude list of matches to target a window, to be excluded from the match.
	Exclude []Match `json:"exclude,omitempty"`
	// Actions defines the action to do on the matching window.
	//
	// This is a json.RawMessage on purpose, since we need to
	// dynamically create the action struct.
	Actions map[string]json.RawMessage `json:"actions,omitempty"`
}

// WindowMatches checks if the window matches the given rule.
func (r Rule) WindowMatches(window Window) bool {
	if r.Type != "window" && r.Type != "" {
		return false
	}

	if len(r.Match) > 0 {
		matched := false
		for _, m := range r.Match {
			if m.WindowMatches(window) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	for _, m := range r.Exclude {
		if m.WindowMatches(window) {
			return false
		}
	}
	return true
}

// WorkspaceMatches checks if the workspace matches the given rule.
func (r *Rule) WorkspaceMatches(workspace Workspace) bool {
	if r.Type != "workspace" {
		return false
	}

	if len(r.Match) > 0 {
		matched := false
		for _, m := range r.Match {
			if m.WorkspaceMatches(workspace) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	for _, m := range r.Exclude {
		if m.WorkspaceMatches(workspace) {
			return false
		}
	}
	return true
}

// Response contains the response from the Niri Socket.
type Response struct {
	Ok map[string]json.RawMessage `json:"Ok"`
}

// ConfiguredMode is the output mode as set in the config file.
type ConfiguredMode struct {
	// Width is the width in physical pixels.
	Width uint16 `json:"width"`
	// Height is the height in physical pixels.
	Height uint16 `json:"height"`
	// Refresh is the refresh rate.
	Refresh float64 `json:"refresh,omitempty"`
}

// ConfiguredPosition is the output position as set in the config file.
type ConfiguredPosition struct {
	// X is the logical x position.
	X int32 `json:"x"`
	// Y is the logical y position.
	Y int32 `json:"y"`
}

// KeyboardLayouts is the configured keyboard layouts.
type KeyboardLayouts struct {
	// Names is the XKB names of the configured layouts.
	Names []string `json:"names"`
	// CurrentIdx is the index of the currently active layout in Names.
	CurrentIdx uint8 `json:"current_idx"`
}

// LayerSurface is the layer-shell surface.
type LayerSurface struct {
	// Namespace is the namespace provided by the layer-shell client.
	Namespace string `json:"namespace"`
	// Output is the name of the output the surface is on.
	Output string `json:"output"`
	// Layer is the layer that the surface is on.
	Layer Layer `json:"layer"`
	// KeyboardInteractivity is the surface's keyboard interactivity mode.
	KeyboardInteractivity LayerSurfaceKeyboardInteractivity `json:"keyboard_interactivity"`
}

// LogicalOutput is the logical output in the compositor's coordinate space.
type LogicalOutput struct {
	// X is the logical x position.
	X int `json:"x"`
	// Y is the logical y position.
	Y int `json:"y"`
	// Width is the width in logical pixels.
	Width int `json:"width"`
	// Height is the height in logical pixels.
	Height int `json:"height"`
	// Scale is the scale factor.
	Scale float64 `json:"scale"`
	// Transform sets the transformation of the output.
	Transform Transform `json:"transform"`
}

// Mode is the output mode.
type Mode struct {
	// Width is the width in physical pixels.
	Width int `json:"width"`
	// Height is the height in physical pixels.
	Height int `json:"height"`
	// RefreshRate is the refresh rate in millihertz.
	RefreshRate int `json:"refresh_rate"`
	// IsPreferred tells whether this mode is preferred by the monitor.
	IsPreferred bool `json:"is_preferred"`
}

// Output is the connected output.
type Output struct {
	// Name is the name of the output.
	Name string `json:"name"`
	// Make is the textual description of the manufacturer.
	Make string `json:"make"`
	// Model is the textual description of the model.
	Model string `json:"model"`
	// Serial is the serial of the output, if known.
	Serial string `json:"serial"`
	// PhysicalSize is the physical width and height of the output in mm, if known.
	PhysicalSize []int `json:"physical_size"`
	// Modes is the available modes for the output.
	Modes []Mode `json:"modes"`
	// CurrentMode is the current mode. None if the output is disabled.
	CurrentMode int `json:"current_mode"`
	// VrrSupported tells whether the output supports variable refresh rate.
	VrrSupported bool `json:"vrr_supported"`
	// VrrEnabled tells whether the variable refresh rate is enabled on the output.
	VrrEnabled bool `json:"vrr_enabled"`
	// Logical is the logical output information. None if the output is not mapped to any logical output (e.g. if it's disabled).
	Logical LogicalOutput `json:"logical"`
}

// Overview is the overview information.
type Overview struct {
	// IsOpen tells whether the overview is currently open or not.
	IsOpen bool `json:"is_open"`
}

// PickedColor is the color picked from the screen.
type PickedColor struct {
	// RGB is the color values as red, green, blue, each ranging from 0.0 to 1.0.
	RGB float64 `json:"rgb"`
}

// VrrToSet is the output variable refresh rate to set.
type VrrToSet struct {
	// Vrr tells whether to enable variable refresh rate or not.
	Vrr bool `json:"vrr"`
	// OnDemand tells to only enable when the output shows a window matching the variable-refresh-rate window rule.
	OnDemand bool `json:"on_demand"`
}

// Window contains the details of a window.
type Window struct {
	// ID is the unique ID of this window.
	//
	// This ID remains constant while this window is open.
	//
	// Do not assume that window IDs will always increase without wrapping, or start at 1.
	// That is an implementation detail subject to change. For example, IDs may change to be
	// randomly generated for each new window.
	ID uint64 `json:"id"`
	// Title is the window title, if set.
	Title string `json:"title"`
	// AppID is the application ID, if set.
	AppID string `json:"app_id"`
	// Pid is the process ID that created the Wayland connection for this window, if known.
	//
	// Currently, windows created by xdg-desktop-portal-gnome will have a None PID, but this
	// may change in the future.
	Pid int `json:"pid"`
	// WorkspaceID is the ID of the workspace this window is on, if any.
	WorkspaceID uint64 `json:"workspace_id"`
	// IsFocused tell whether this window is currently focused.
	//
	// There can either be one focused window, or zero (e.g. when a layer-shell surface has focus).
	IsFocused bool `json:"is_focused"`
	// IsFloating tells whether this window is currently floating.
	//
	// If the window isn't floating, then it's in the tiling layout.
	IsFloating bool `json:"is_floating"`
	// IsUrgent tells whether this window requests your attention.
	IsUrgent bool `json:"is_urgent"`
	// Matched tells if the window matches a rule defined by nirimgr rules.
	//
	// This is not a part of the Niri Window model.
	Matched bool
}

// Workspace is the workspace.
type Workspace struct {
	// ID is the unique ID of this workspace.
	//
	// This id remains constant regardless of the workspace moving around and across monitors.
	// Do not assume that workspace IDs will always increase without wrapping, or start at 1.
	// That is an implementation detail subject to change.
	// For example, IDs may change to be randomly generated for each new workspace.
	ID uint64 `json:"id"`
	// Idx is the index of the workspace on this monitor.
	//
	// This is the same index you can use for requests like niri msg action focus-workspace.
	// This index will change as you move and re-order workspace. It is merely the workspace's
	// current position on its monitor. Workspaces on different monitors can have the same index.
	// If you need a unique workspace id that doesn’t change, see Id.
	Idx uint8 `json:"idx"`
	// Name is the optional name of the workspace.
	Name string `json:"name"`
	// Output is the name of the output that the workspace is on.
	//
	// Can be None if no outputs are currently connected.
	Output string `json:"output"`
	// IsUrgent tells whether the workspace currently has an urgent window in its output.
	IsUrgent bool `json:"is_urgent"`
	// IsActive tells whether the workspace is currently active on its output.
	//
	// Every output has one active workspace, the one that is currently visible on that output.
	IsActive bool `json:"is_active"`
	// IsFocused tells whether the workspace is currently focused.
	//
	// There's only one focused workspace across all outputs.
	IsFocused bool `json:"is_focused"`
	// ActiveWindowID is the ID of the active window on this workspace, if any.
	ActiveWindowID uint64 `json:"active_window_id"`
	// Matched tells if the workspace matches a rule defined by nirimgr rules.
	//
	// This is not a part of the Niri Workspace model.
	Matched bool
}

// ReferenceKeys contains the possible keys a WorkspaceReferenceArg can have.
//
// This is used when setting the reference dynamically on matching workspaces.
type ReferenceKeys struct {
	ID    uint64
	Index uint8
	Name  string
}

// PossibleKeys contains the possible keys an action could have.
//
// This is used when setting the action IDs dynamically during matching of windows.
type PossibleKeys struct {
	ID             uint64
	WindowID       uint64
	ActiveWindowID uint64
	WorkspaceID    uint64
	Index          uint8
	Reference      ReferenceKeys
}
