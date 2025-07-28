package actions

import (
	"testing"
)

func TestPowerOnMonitors(t *testing.T) {
	a := PowerOnMonitors{AName{"PowerOnMonitors"}}
	if a.GetName() != "PowerOnMonitors" {
		t.Errorf("expected name PowerOnMonitors, got %s", a.GetName())
	}
}

func TestSpawn(t *testing.T) {
	cmd := []string{"ls", "-l"}
	a := Spawn{AName{"Spawn"}, cmd}
	if a.GetName() != "Spawn" {
		t.Errorf("expected name Spawn, got %s", a.GetName())
	}
	if len(a.Command) != 2 || a.Command[0] != "ls" || a.Command[1] != "-l" {
		t.Errorf("unexpected command: %v", a.Command)
	}
}

func TestDoScreenTransition(t *testing.T) {
	a := DoScreenTransition{AName{"DoScreenTransition"}, 100}
	if a.GetName() != "DoScreenTransition" {
		t.Errorf("expected name DoScreenTransition, got %s", a.GetName())
	}
	if a.DelayMs != 100 {
		t.Errorf("expected DelayMs 100, got %d", a.DelayMs)
	}
}

func TestScreenshot(t *testing.T) {
	a := Screenshot{AName{"Screenshot"}, true}
	if a.GetName() != "Screenshot" {
		t.Errorf("expected name Screenshot, got %s", a.GetName())
	}
	if !a.ShowPointer {
		t.Errorf("expected ShowPointer true")
	}
}

func TestScreenshotScreen(t *testing.T) {
	a := ScreenshotScreen{AName{"ScreenshotScreen"}, true, false}
	if a.GetName() != "ScreenshotScreen" {
		t.Errorf("expected name ScreenshotScreen, got %s", a.GetName())
	}
	if !a.WriteToDisk {
		t.Errorf("expected WriteToDisk true")
	}
	if a.ShowPointer {
		t.Errorf("expected ShowPointer false")
	}
}

func TestScreenshotWindow(t *testing.T) {
	a := ScreenshotWindow{AName{"ScreenshotWindow"}, 42, true}
	if a.GetName() != "ScreenshotWindow" {
		t.Errorf("expected name ScreenshotWindow, got %s", a.GetName())
	}
	if a.ID != 42 {
		t.Errorf("expected ID 42, got %d", a.ID)
	}
	if !a.WriteToDisk {
		t.Errorf("expected WriteToDisk true")
	}
}

func TestToggleKeyboardShortcutsInhibit(t *testing.T) {
	a := ToggleKeyboardShortcutsInhibit{AName{"ToggleKeyboardShortcutsInhibit"}}
	if a.GetName() != "ToggleKeyboardShortcutsInhibit" {
		t.Errorf("expected name ToggleKeyboardShortcutsInhibit, got %s", a.GetName())
	}
}

func TestCloseWindow(t *testing.T) {
	a := CloseWindow{AName{"CloseWindow"}, 99}
	if a.GetName() != "CloseWindow" {
		t.Errorf("expected name CloseWindow, got %s", a.GetName())
	}
	if a.ID != 99 {
		t.Errorf("expected ID 99, got %d", a.ID)
	}
}

func TestFullscreenWindow(t *testing.T) {
	a := FullscreenWindow{AName{"FullscreenWindow"}, 77}
	if a.GetName() != "FullscreenWindow" {
		t.Errorf("expected name FullscreenWindow, got %s", a.GetName())
	}
	if a.ID != 77 {
		t.Errorf("expected ID 77, got %d", a.ID)
	}
}

func TestToggleWindowedFullscreen(t *testing.T) {
	a := ToggleWindowedFullscreen{AName{"ToggleWindowedFullscreen"}, 88}
	if a.GetName() != "ToggleWindowedFullscreen" {
		t.Errorf("expected name ToggleWindowedFullscreen, got %s", a.GetName())
	}
	if a.ID != 88 {
		t.Errorf("expected ID 88, got %d", a.ID)
	}
}

func TestFocusWindow(t *testing.T) {
	a := FocusWindow{AName{"FocusWindow"}, 123}
	if a.GetName() != "FocusWindow" {
		t.Errorf("expected name FocusWindow, got %s", a.GetName())
	}
	if a.ID != 123 {
		t.Errorf("expected ID 123, got %d", a.ID)
	}
}

func TestFocusWindowInColumn(t *testing.T) {
	a := FocusWindowInColumn{AName{"FocusWindowInColumn"}, 2}
	if a.GetName() != "FocusWindowInColumn" {
		t.Errorf("expected name FocusWindowInColumn, got %s", a.GetName())
	}
	if a.Index != 2 {
		t.Errorf("expected Index 2, got %d", a.Index)
	}
}
