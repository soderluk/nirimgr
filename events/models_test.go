package events

import (
	"testing"

	"github.com/soderluk/nirimgr/models"
)

func TestEName_GetName(t *testing.T) {
	name := "TestEvent"
	e := EName{Name: name}
	if e.GetName() != name {
		t.Errorf("expected %s, got %s", name, e.GetName())
	}
}

func TestWorkspacesChanged(t *testing.T) {
	ws := &models.Workspace{ID: 1}
	e := WorkspacesChanged{
		EName:      EName{Name: "WorkspacesChanged"},
		Workspaces: []*models.Workspace{ws},
	}
	if len(e.Workspaces) != 1 || e.Workspaces[0].ID != 1 {
		t.Errorf("unexpected workspace data")
	}
}

func TestWorkspaceUrgencyChanged(t *testing.T) {
	e := WorkspaceUrgencyChanged{
		EName:  EName{Name: "WorkspaceUrgencyChanged"},
		ID:     42,
		Urgent: true,
	}
	if !e.Urgent || e.ID != 42 {
		t.Errorf("unexpected urgency or ID")
	}
}

func TestWorkspaceActivated(t *testing.T) {
	e := WorkspaceActivated{
		EName:   EName{Name: "WorkspaceActivated"},
		ID:      7,
		Focused: true,
	}
	if !e.Focused || e.ID != 7 {
		t.Errorf("unexpected focus or ID")
	}
}

func TestWorkspaceActiveWindowChanged(t *testing.T) {
	e := WorkspaceActiveWindowChanged{
		EName:          EName{Name: "WorkspaceActiveWindowChanged"},
		WorkspaceID:    3,
		ActiveWindowID: 99,
	}
	if e.WorkspaceID != 3 || e.ActiveWindowID != 99 {
		t.Errorf("unexpected workspace or window ID")
	}
}

func TestWindowsChanged(t *testing.T) {
	w := &models.Window{ID: 5}
	e := WindowsChanged{
		EName:   EName{Name: "WindowsChanged"},
		Windows: []*models.Window{w},
	}
	if len(e.Windows) != 1 || e.Windows[0].ID != 5 {
		t.Errorf("unexpected windows data")
	}
}

func TestWindowOpenedOrChanged(t *testing.T) {
	w := &models.Window{ID: 8}
	e := WindowOpenedOrChanged{
		EName:  EName{Name: "WindowOpenedOrChanged"},
		Window: w,
	}
	if e.Window == nil || e.Window.ID != 8 {
		t.Errorf("unexpected window data")
	}
}

func TestWindowClosed(t *testing.T) {
	e := WindowClosed{
		EName: EName{Name: "WindowClosed"},
		ID:    11,
	}
	if e.ID != 11 {
		t.Errorf("unexpected closed window ID")
	}
}

func TestWindowFocusChanged(t *testing.T) {
	e := WindowFocusChanged{
		EName: EName{Name: "WindowFocusChanged"},
		ID:    13,
	}
	if e.ID != 13 {
		t.Errorf("unexpected focused window ID")
	}
}

func TestWindowUrgencyChanged(t *testing.T) {
	e := WindowUrgencyChanged{
		EName:  EName{Name: "WindowUrgencyChanged"},
		ID:     15,
		Urgent: true,
	}
	if !e.Urgent || e.ID != 15 {
		t.Errorf("unexpected urgency or ID")
	}
}

func TestKeyboardLayoutsChanged(t *testing.T) {
	kl := models.KeyboardLayouts{Names: []string{"us", "se"}}
	e := KeyboardLayoutsChanged{
		EName:           EName{Name: "KeyboardLayoutsChanged"},
		KeyboardLayouts: kl,
	}
	if len(e.KeyboardLayouts.Names) != 2 {
		t.Errorf("unexpected keyboard layouts")
	}
}

func TestKeyboardLayoutSwitched(t *testing.T) {
	e := KeyboardLayoutSwitched{
		EName: EName{Name: "KeyboardLayoutSwitched"},
		Idx:   2,
	}
	if e.Idx != 2 {
		t.Errorf("unexpected layout index")
	}
}

func TestOverviewOpenedOrClosed(t *testing.T) {
	e := OverviewOpenedOrClosed{
		EName:  EName{Name: "OverviewOpenedOrClosed"},
		IsOpen: true,
	}
	if !e.IsOpen {
		t.Errorf("overview should be open")
	}
}
