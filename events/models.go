package events

import "github.com/soderluk/nirimgr/models"

// If more events are added in Niri, we must define them here, and add them to the EventRegistry.

// Event defines the "base" interface for all the events.
//
// NOTE: We have to use GetName, since the field is called Name.
type Event interface {
	GetName() string
}

// EName defines the name of the event.
type EName struct {
	Name string
}

// GetName returns the event name.
func (e EName) GetName() string {
	return e.Name
}

// WorkspacesChanged when the workspace configuration has changed.
type WorkspacesChanged struct {
	EName
	// Workspaces contains the new workspace configuration.
	//
	// This configuration completely replaces the previous configuration. If any workspaces
	// are missing from here, then they were deleted.
	Workspaces []*models.Workspace `json:"workspaces"`
}

// WorkspaceUrgencyChanged when the workspace urgency changed.
type WorkspaceUrgencyChanged struct {
	EName
	// ID the ID of the workspace.
	ID uint64 `json:"id"`
	// Urgent tells if this workspace has an urgent window.
	Urgent bool `json:"urgent"`
}

// WorkspaceActivated when a workspace was activated on an output.
type WorkspaceActivated struct {
	EName
	// ID the ID of the newly active workspace.
	ID uint64 `json:"id"`
	// Focused tells if this workspace also became focused.
	//
	// If true, this is now the single focused workspace. All other workspaces are no longer
	// focused, but they may remain active on their respective outputs.
	Focused bool `json:"focused"`
}

// WorkspaceActiveWindowChanged when an active window changed on a workspace.
type WorkspaceActiveWindowChanged struct {
	EName
	// WorkspaceID the ID of the workspace on which the active window changed.
	WorkspaceID uint64 `json:"workspace_id"`
	// ActiveWindowID the ID of the new active window, if any.
	ActiveWindowID uint64 `json:"active_window_id"`
}

// WindowsChanged when the window configuration has changed.
type WindowsChanged struct {
	EName
	// Windows contains the new window configuration.
	//
	// This configuration completely replaces the previous configuration. If any windows
	// are missing from here, then they were closed.
	Windows []*models.Window `json:"windows"`
}

// WindowOpenedOrChanged when a new toplevel window was opened, or an existing toplevel window changed.
type WindowOpenedOrChanged struct {
	EName
	// Window contains the new or updated window.
	//
	// If the window is focused, all other windows are no longer focused.
	Window *models.Window `json:"window"`
}

// WindowClosed when a toplevel window was closed.
type WindowClosed struct {
	EName
	// ID the ID of the removed window.
	ID uint64 `json:"id"`
}

// WindowFocusChanged when a window focus changed.
//
// All other windows are no longer focused.
type WindowFocusChanged struct {
	EName
	// ID the ID of the newly focused window, or omitted if no window is now focused.
	ID uint64 `json:"id"`
}

// WindowUrgencyChanged when a window urgency changed.
type WindowUrgencyChanged struct {
	EName
	// ID the ID of the window.
	ID uint64 `json:"id"`
	// Urgent the new urgency state of the window.
	Urgent bool `json:"urgent"`
}

// KeyboardLayoutsChanged when the configured keyboard layouts have changed.
type KeyboardLayoutsChanged struct {
	EName
	// KeyboardLayouts contains the new keyboard layout configuration.
	KeyboardLayouts models.KeyboardLayouts `json:"keyboard_layouts"`
}

// KeyboardLayoutSwitched when the keyboard layout switched.
type KeyboardLayoutSwitched struct {
	EName
	// Idx contains the index of the newly active layout.
	Idx uint8 `json:"idx"`
}

// OverviewOpenedOrClosed when the overview was opened or closed.
type OverviewOpenedOrClosed struct {
	EName
	// IsOpen contains the new state of the overview.
	IsOpen bool `json:"is_open"`
}
