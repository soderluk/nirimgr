package events

import (
	"github.com/soderluk/nirimgr/models"
)

// If more events are added in Niri, we must define them here, and add them to the EventRegistry.

// Event defines the "base" interface for all the events.
//
// NOTE: We have to use GetName, since the field is called Name.
type Event interface {
	GetName() string
	GetPossibleKeys() models.PossibleKeys
}

// EName defines the name of the event.
type EName struct {
	Name string
}

// GetName returns the event name.
func (e EName) GetName() string {
	return e.Name
}

// GetPossibleKeys extracts relevant IDs and fields from any event.
// This is a default implementation that should work for most events.
func (e EName) GetPossibleKeys() models.PossibleKeys {
	// This default implementation returns empty keys since EName itself has no useful fields
	// Individual event types can override this method if they have specific fields to extract
	return models.PossibleKeys{}
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

// GetPossibleKeys extracts the workspace ID from this event.
func (w WorkspaceUrgencyChanged) GetPossibleKeys() models.PossibleKeys {
	return models.PossibleKeys{
		ID:          w.ID,
		WorkspaceID: w.ID,
	}
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

// GetPossibleKeys extracts the workspace ID from this event.
func (w WorkspaceActivated) GetPossibleKeys() models.PossibleKeys {
	return models.PossibleKeys{
		ID:          w.ID,
		WorkspaceID: w.ID,
	}
}

// WorkspaceActiveWindowChanged when an active window changed on a workspace.
type WorkspaceActiveWindowChanged struct {
	EName
	// WorkspaceID the ID of the workspace on which the active window changed.
	WorkspaceID uint64 `json:"workspace_id"`
	// ActiveWindowID the ID of the new active window, if any.
	ActiveWindowID uint64 `json:"active_window_id"`
}

// GetPossibleKeys extracts the workspace and window IDs from this event.
func (w WorkspaceActiveWindowChanged) GetPossibleKeys() models.PossibleKeys {
	return models.PossibleKeys{
		WorkspaceID:    w.WorkspaceID,
		ActiveWindowID: w.ActiveWindowID,
	}
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

// GetPossibleKeys extracts the window ID from this event.
func (w WindowClosed) GetPossibleKeys() models.PossibleKeys {
	return models.PossibleKeys{
		ID:       w.ID,
		WindowID: w.ID,
	}
}

// WindowFocusChanged when a window focus changed.
//
// All other windows are no longer focused.
type WindowFocusChanged struct {
	EName
	// ID the ID of the newly focused window, or omitted if no window is now focused.
	ID uint64 `json:"id"`
}

// GetPossibleKeys extracts the window ID from this event.
func (w WindowFocusChanged) GetPossibleKeys() models.PossibleKeys {
	return models.PossibleKeys{
		ID:       w.ID,
		WindowID: w.ID,
	}
}

// WindowUrgencyChanged when a window urgency changed.
type WindowUrgencyChanged struct {
	EName
	// ID the ID of the window.
	ID uint64 `json:"id"`
	// Urgent the new urgency state of the window.
	Urgent bool `json:"urgent"`
}

// GetPossibleKeys extracts the window ID from this event.
func (w WindowUrgencyChanged) GetPossibleKeys() models.PossibleKeys {
	return models.PossibleKeys{
		ID:       w.ID,
		WindowID: w.ID,
	}
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

// GetPossibleKeys extracts the layout index from this event.
func (k KeyboardLayoutSwitched) GetPossibleKeys() models.PossibleKeys {
	return models.PossibleKeys{
		Index: k.Idx,
	}
}

// OverviewOpenedOrClosed when the overview was opened or closed.
type OverviewOpenedOrClosed struct {
	EName
	// IsOpen contains the new state of the overview.
	IsOpen bool `json:"is_open"`
}

// ConfigLoaded when the configuration was reloaded
//
// This will always be received when connecting to the event stream,
// indicating the last config load attempt
type ConfigLoaded struct {
	// Failed indicates that the configuration couldn't be reloaded.
	//
	// This can happen e.g. when the config validation
	// fails.
	Failed bool `json:"failed"`
}
