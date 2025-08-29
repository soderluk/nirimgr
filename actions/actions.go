// Package actions contains all the actions Niri currently supports.
//
// A thing to note: we need to add the models.AName embedded struct to all the actions:
//
//	models.AName{Name: "ActionName"}
//	Example:
//	type Quit struct{
//		models.AName
//		SkipConfirmation bool `json:"skip_confirmation"`
//	}
//
// because we don't want to add the receiver functions for all 130+ actions.
// Supporting just one GetName() function for the interface, gives us some leeway when
// working with the actions.
//
// See: https://yalter.github.io/niri/niri_ipc/enum.Action.html# for more details.
package actions

import (
	"encoding/json"
	"log/slog"
	"reflect"

	"github.com/soderluk/nirimgr/internal/common"
	"github.com/soderluk/nirimgr/models"
)

// HandleDynamicIDs assigns the given possible keys to the action.
//
// These are IDs or references we need to dynamically assign to the action,
// if they exist. For the reference, the ID takes precedence, then Index, and at
// last the Name.
// This is used for the matching of workspaces and windows. Since we only get to know
// the necessary IDs during runtime, we want to be able to set them dynamically.
func HandleDynamicIDs(a Action, possibleKeys models.PossibleKeys) Action {
	value := reflect.ValueOf(a)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Set direct fields if present
	if possibleKeys.ID != 0 {
		common.SetUintField(value, "ID", possibleKeys.ID)
	}
	if possibleKeys.WindowID != 0 {
		common.SetUintField(value, "WindowID", possibleKeys.WindowID)
	}
	if possibleKeys.ActiveWindowID != 0 {
		common.SetUintField(value, "ActiveWindowID", possibleKeys.ActiveWindowID)
	}
	if possibleKeys.WorkspaceID != 0 {
		common.SetUintField(value, "WorkspaceID", possibleKeys.WorkspaceID)
	}
	if possibleKeys.Index != 0 {
		common.SetUintField(value, "Index", possibleKeys.Index)
	}

	// Handle Reference struct if present
	referenceField := value.FieldByName("Reference")
	if referenceField.IsValid() && referenceField.CanSet() && referenceField.Kind() == reflect.Struct {
		if possibleKeys.Reference.ID != 0 {
			common.SetUintField(referenceField, "ID", possibleKeys.Reference.ID)
		} else if possibleKeys.Reference.Index != 0 {
			common.SetUintField(referenceField, "Index", possibleKeys.Reference.Index)
		} else if possibleKeys.Reference.Name != "" {
			common.SetStringField(referenceField, "Name", possibleKeys.Reference.Name)
		}
	}
	return a
}

// FromRegistry returns the populated model from the ActionRegistry by given name.
func FromRegistry(name string, data []byte) Action {
	model, ok := ActionRegistry[name]
	if !ok {
		slog.Error("Could not get action model for action", "name", name)
		return nil
	}
	action := model()
	if err := json.Unmarshal(data, action); err != nil {
		slog.Error("Could not unmarshal action", "name", name, "error", err.Error())
		return nil
	}
	return action
}

// ParseRawActions parses the actions into their respective structs.
func ParseRawActions(rawActions map[string]json.RawMessage) []Action {
	var actionList []Action

	for name, raw := range rawActions {
		action := FromRegistry(name, raw)
		if action == nil {
			continue
		}
		actionList = append(actionList, action)
	}
	return actionList
}

// ActionRegistry contains all the actions Niri currently sends.
//
// The key needs to be the action name, and it should return the correct action model, and set
// its AName embedded struct. If you know of a better way to handle this, please let me know.
var ActionRegistry = map[string]func() Action{
	"Quit":               func() Action { return &Quit{AName: AName{Name: "Quit"}} },
	"PowerOffMonitors":   func() Action { return &PowerOffMonitors{AName: AName{Name: "PowerOffMonitors"}} },
	"PowerOnMonitors":    func() Action { return &PowerOnMonitors{AName: AName{Name: "PowerOnMonitors"}} },
	"Spawn":              func() Action { return &Spawn{AName: AName{Name: "Spawn"}} },
	"DoScreenTransition": func() Action { return &DoScreenTransition{AName: AName{Name: "DoScreenTransition"}} },
	"Screenshot":         func() Action { return &Screenshot{AName: AName{Name: "Screenshot"}} },
	"ScreenshotScreen":   func() Action { return &ScreenshotScreen{AName: AName{Name: "ScreenshotScreen"}} },
	"ScreenshotWindow":   func() Action { return &ScreenshotWindow{AName: AName{Name: "ScreenshotWindow"}} },
	"ToggleKeyboardShortcutsInhibit": func() Action {
		return &ToggleKeyboardShortcutsInhibit{AName: AName{Name: "ToggleKeyboardShortcutsInhibit"}}
	},
	"CloseWindow":                 func() Action { return &CloseWindow{AName: AName{Name: "CloseWindow"}} },
	"FullscreenWindow":            func() Action { return &FullscreenWindow{AName: AName{Name: "FullscreenWindow"}} },
	"ToggleWindowedFullscreen":    func() Action { return &ToggleWindowedFullscreen{AName: AName{Name: "ToggleWindowedFullscreen"}} },
	"FocusWindow":                 func() Action { return &FocusWindow{AName: AName{Name: "FocusWindow"}} },
	"FocusWindowInColumn":         func() Action { return &FocusWindowInColumn{AName: AName{Name: "FocusWindowInColumn"}} },
	"FocusWindowPrevious":         func() Action { return &FocusWindowPrevious{AName: AName{Name: "FocusWindowPrevious"}} },
	"FocusColumnLeft":             func() Action { return &FocusColumnLeft{AName: AName{Name: "FocusColumnLeft"}} },
	"FocusColumnRight":            func() Action { return &FocusColumnRight{AName: AName{Name: "FocusColumnRight"}} },
	"FocusColumnFirst":            func() Action { return &FocusColumnFirst{AName: AName{Name: "FocusColumnFirst"}} },
	"FocusColumnLast":             func() Action { return &FocusColumnLast{AName: AName{Name: "FocusColumnLast"}} },
	"FocusColumnRightOrFirst":     func() Action { return &FocusColumnRightOrFirst{AName: AName{Name: "FocusColumnRightOrFirst"}} },
	"FocusColumnLeftOrLast":       func() Action { return &FocusColumnLeftOrLast{AName: AName{Name: "FocusColumnLeftOrLast"}} },
	"FocusColumn":                 func() Action { return &FocusColumn{AName: AName{Name: "FocusColumn"}} },
	"FocusWindowOrMonitorUp":      func() Action { return &FocusWindowOrMonitorUp{AName: AName{Name: "FocusWindowOrMonitorUp"}} },
	"FocusWindowOrMonitorDown":    func() Action { return &FocusWindowOrMonitorDown{AName: AName{Name: "FocusWindowOrMonitorDown"}} },
	"FocusColumnOrMonitorLeft":    func() Action { return &FocusColumnOrMonitorLeft{AName: AName{Name: "FocusColumnOrMonitorLeft"}} },
	"FocusColumnOrMonitorRight":   func() Action { return &FocusColumnOrMonitorRight{AName: AName{Name: "FocusColumnOrMonitorRight"}} },
	"FocusWindowDown":             func() Action { return &FocusWindowDown{AName: AName{Name: "FocusWindowDown"}} },
	"FocusWindowUp":               func() Action { return &FocusWindowUp{AName: AName{Name: "FocusWindowUp"}} },
	"FocusWindowDownOrColumnLeft": func() Action { return &FocusWindowDownOrColumnLeft{AName: AName{Name: "FocusWindowDownOrColumnLeft"}} },
	"FocusWindowDownOrColumnRight": func() Action {
		return &FocusWindowDownOrColumnRight{AName: AName{Name: "FocusWindowDownOrColumnRight"}}
	},
	"FocusWindowUpOrColumnLeft":  func() Action { return &FocusWindowUpOrColumnLeft{AName: AName{Name: "FocusWindowUpOrColumnLeft"}} },
	"FocusWindowUpOrColumnRight": func() Action { return &FocusWindowUpOrColumnRight{AName: AName{Name: "FocusWindowUpOrColumnRight"}} },
	"FocusWindowOrWorkspaceDown": func() Action { return &FocusWindowOrWorkspaceDown{AName: AName{Name: "FocusWindowOrWorkspaceDown"}} },
	"FocusWindowOrWorkspaceUp":   func() Action { return &FocusWindowOrWorkspaceUp{AName: AName{Name: "FocusWindowOrWorkspaceUp"}} },
	"FocusWindowTop":             func() Action { return &FocusWindowTop{AName: AName{Name: "FocusWindowTop"}} },
	"FocusWindowBottom":          func() Action { return &FocusWindowBottom{AName: AName{Name: "FocusWindowBottom"}} },
	"FocusWindowDownOrTop":       func() Action { return &FocusWindowDownOrTop{AName: AName{Name: "FocusWindowDownOrTop"}} },
	"FocusWindowUpOrBottom":      func() Action { return &FocusWindowUpOrBottom{AName: AName{Name: "FocusWindowUpOrBottom"}} },
	"MoveColumnLeft":             func() Action { return &MoveColumnLeft{AName: AName{Name: "MoveColumnLeft"}} },
	"MoveColumnRight":            func() Action { return &MoveColumnRight{AName: AName{Name: "MoveColumnRight"}} },
	"MoveColumnToFirst":          func() Action { return &MoveColumnToFirst{AName: AName{Name: "MoveColumnToFirst"}} },
	"MoveColumnToLast":           func() Action { return &MoveColumnToLast{AName: AName{Name: "MoveColumnToLast"}} },
	"MoveColumnLeftOrToMonitorLeft": func() Action {
		return &MoveColumnLeftOrToMonitorLeft{AName: AName{Name: "MoveColumnLeftOrToMonitorLeft"}}
	},
	"MoveColumnRightOrToMonitorRight": func() Action {
		return &MoveColumnRightOrToMonitorRight{AName: AName{Name: "MoveColumnRightOrToMonitorRight"}}
	},
	"MoveColumnToIndex": func() Action { return &MoveColumnToIndex{AName: AName{Name: "MoveColumnToIndex"}} },
	"MoveWindowDown":    func() Action { return &MoveWindowDown{AName: AName{Name: "MoveWindowDown"}} },
	"MoveWindowUp":      func() Action { return &MoveWindowUp{AName: AName{Name: "MoveWindowUp"}} },
	"MoveWindowDownOrToWorkspaceDown": func() Action {
		return &MoveWindowDownOrToWorkspaceDown{AName: AName{Name: "MoveWindowDownOrToWorkspaceDown"}}
	},
	"MoveWindowUpOrToWorkspaceUp": func() Action { return &MoveWindowUpOrToWorkspaceUp{AName: AName{Name: "MoveWindowUpOrToWorkspaceUp"}} },
	"ConsumeOrExpelWindowLeft":    func() Action { return &ConsumeOrExpelWindowLeft{AName: AName{Name: "ConsumeOrExpelWindowLeft"}} },
	"ConsumeOrExpelWindowRight":   func() Action { return &ConsumeOrExpelWindowRight{AName: AName{Name: "ConsumeOrExpelWindowRight"}} },
	"ConsumeWindowIntoColumn":     func() Action { return &ConsumeWindowIntoColumn{AName: AName{Name: "ConsumeWindowIntoColumn"}} },
	"ExpelWindowFromColumn":       func() Action { return &ExpelWindowFromColumn{AName: AName{Name: "ExpelWindowFromColumn"}} },
	"SwapWindowRight":             func() Action { return &SwapWindowRight{AName: AName{Name: "SwapWindowRight"}} },
	"SwapWindowLeft":              func() Action { return &SwapWindowLeft{AName: AName{Name: "SwapWindowLeft"}} },
	"ToggleColumnTabbedDisplay":   func() Action { return &ToggleColumnTabbedDisplay{AName: AName{Name: "ToggleColumnTabbedDisplay"}} },
	"SetColumnDisplay":            func() Action { return &SetColumnDisplay{AName: AName{Name: "SetColumnDisplay"}} },
	"CenterColumn":                func() Action { return &CenterColumn{AName: AName{Name: "CenterColumn"}} },
	"CenterWindow":                func() Action { return &CenterWindow{AName: AName{Name: "CenterWindow"}} },
	"CenterVisibleColumns":        func() Action { return &CenterVisibleColumns{AName: AName{Name: "CenterVisibleColumns"}} },
	"FocusWorkspaceDown":          func() Action { return &FocusWorkspaceDown{AName: AName{Name: "FocusWorkspaceDown"}} },
	"FocusWorkspaceUp":            func() Action { return &FocusWorkspaceUp{AName: AName{Name: "FocusWorkspaceUp"}} },
	"FocusWorkspace":              func() Action { return &FocusWorkspace{AName: AName{Name: "FocusWorkspace"}} },
	"FocusWorkspacePrevious":      func() Action { return &FocusWorkspacePrevious{AName: AName{Name: "FocusWorkspacePrevious"}} },
	"MoveWindowToWorkspaceDown":   func() Action { return &MoveWindowToWorkspaceDown{AName: AName{Name: "MoveWindowToWorkspaceDown"}} },
	"MoveWindowToWorkspaceUp":     func() Action { return &MoveWindowToWorkspaceUp{AName: AName{Name: "MoveWindowToWorkspaceUp"}} },
	"MoveWindowToWorkspace":       func() Action { return &MoveWindowToWorkspace{AName: AName{Name: "MoveWindowToWorkspace"}} },
	"MoveColumnToWorkspaceDown":   func() Action { return &MoveColumnToWorkspaceDown{AName: AName{Name: "MoveColumnToWorkspaceDown"}} },
	"MoveColumnToWorkspaceUp":     func() Action { return &MoveColumnToWorkspaceUp{AName: AName{Name: "MoveColumnToWorkspaceUp"}} },
	"MoveColumnToWorkspace":       func() Action { return &MoveColumnToWorkspace{AName: AName{Name: "MoveColumnToWorkspace"}} },
	"MoveWorkspaceDown":           func() Action { return &MoveWorkspaceDown{AName: AName{Name: "MoveWorkspaceDown"}} },
	"MoveWorkspaceUp":             func() Action { return &MoveWorkspaceUp{AName: AName{Name: "MoveWorkspaceUp"}} },
	"MoveWorkspaceToIndex":        func() Action { return &MoveWorkspaceToIndex{AName: AName{Name: "MoveWorkspaceToIndex"}} },
	"SetWorkspaceName":            func() Action { return &SetWorkspaceName{AName: AName{Name: "SetWorkspaceName"}} },
	"UnsetWorkspaceName":          func() Action { return &UnsetWorkspaceName{AName: AName{Name: "UnsetWorkspaceName"}} },
	"FocusMonitorLeft":            func() Action { return &FocusMonitorLeft{AName: AName{Name: "FocusMonitorLeft"}} },
	"FocusMonitorRight":           func() Action { return &FocusMonitorRight{AName: AName{Name: "FocusMonitorRight"}} },
	"FocusMonitorDown":            func() Action { return &FocusMonitorDown{AName: AName{Name: "FocusMonitorDown"}} },
	"FocusMonitorUp":              func() Action { return &FocusMonitorUp{AName: AName{Name: "FocusMonitorUp"}} },
	"FocusMonitorPrevious":        func() Action { return &FocusMonitorPrevious{AName: AName{Name: "FocusMonitorPrevious"}} },
	"FocusMonitorNext":            func() Action { return &FocusMonitorNext{AName: AName{Name: "FocusMonitorNext"}} },
	"FocusMonitor":                func() Action { return &FocusMonitor{AName: AName{Name: "FocusMonitor"}} },
	"MoveWindowToMonitorLeft":     func() Action { return &MoveWindowToMonitorLeft{AName: AName{Name: "MoveWindowToMonitorLeft"}} },
	"MoveWindowToMonitorRight":    func() Action { return &MoveWindowToMonitorRight{AName: AName{Name: "MoveWindowToMonitorRight"}} },
	"MoveWindowToMonitorDown":     func() Action { return &MoveWindowToMonitorDown{AName: AName{Name: "MoveWindowToMonitorDown"}} },
	"MoveWindowToMonitorUp":       func() Action { return &MoveWindowToMonitorUp{AName: AName{Name: "MoveWindowToMonitorUp"}} },
	"MoveWindowToMonitorPrevious": func() Action { return &MoveWindowToMonitorPrevious{AName: AName{Name: "MoveWindowToMonitorPrevious"}} },
	"MoveWindowToMonitorNext":     func() Action { return &MoveWindowToMonitorNext{AName: AName{Name: "MoveWindowToMonitorNext"}} },
	"MoveWindowToMonitor":         func() Action { return &MoveWindowToMonitor{AName: AName{Name: "MoveWindowToMonitor"}} },
	"MoveColumnToMonitorLeft":     func() Action { return &MoveColumnToMonitorLeft{AName: AName{Name: "MoveColumnToMonitorLeft"}} },
	"MoveColumnToMonitorRight":    func() Action { return &MoveColumnToMonitorRight{AName: AName{Name: "MoveColumnToMonitorRight"}} },
	"MoveColumnToMonitorDown":     func() Action { return &MoveColumnToMonitorDown{AName: AName{Name: "MoveColumnToMonitorDown"}} },
	"MoveColumnToMonitorUp":       func() Action { return &MoveColumnToMonitorUp{AName: AName{Name: "MoveColumnToMonitorUp"}} },
	"MoveColumnToMonitorPrevious": func() Action { return &MoveColumnToMonitorPrevious{AName: AName{Name: "MoveColumnToMonitorPrevious"}} },
	"MoveColumnToMonitorNext":     func() Action { return &MoveColumnToMonitorNext{AName: AName{Name: "MoveColumnToMonitorNext"}} },
	"MoveColumnToMonitor":         func() Action { return &MoveColumnToMonitor{AName: AName{Name: "MoveColumnToMonitor"}} },
	"SetWindowWidth":              func() Action { return &SetWindowWidth{AName: AName{Name: "SetWindowWidth"}} },
	"SetWindowHeight":             func() Action { return &SetWindowHeight{AName: AName{Name: "SetWindowHeight"}} },
	"ResetWindowHeight":           func() Action { return &ResetWindowHeight{AName: AName{Name: "ResetWindowHeight"}} },
	"SwitchPresetColumnWidth":     func() Action { return &SwitchPresetColumnWidth{AName: AName{Name: "SwitchPresetColumnWidth"}} },
	"SwitchPresetWindowWidth":     func() Action { return &SwitchPresetWindowWidth{AName: AName{Name: "SwitchPresetWindowWidth"}} },
	"SwitchPresetWindowHeight":    func() Action { return &SwitchPresetWindowHeight{AName: AName{Name: "SwitchPresetWindowHeight"}} },
	"MaximizeColumn":              func() Action { return &MaximizeColumn{AName: AName{Name: "MaximizeColumn"}} },
	"SetColumnWidth":              func() Action { return &SetColumnWidth{AName: AName{Name: "SetColumnWidth"}} },
	"ExpandColumnToAvailableWidth": func() Action {
		return &ExpandColumnToAvailableWidth{AName: AName{Name: "ExpandColumnToAvailableWidth"}}
	},
	"SwitchLayout":                func() Action { return &SwitchLayout{AName: AName{Name: "SwitchLayout"}} },
	"ShowHotkeyOverlay":           func() Action { return &ShowHotkeyOverlay{AName: AName{Name: "ShowHotkeyOverlay"}} },
	"MoveWorkspaceToMonitorLeft":  func() Action { return &MoveWorkspaceToMonitorLeft{AName: AName{Name: "MoveWorkspaceToMonitorLeft"}} },
	"MoveWorkspaceToMonitorRight": func() Action { return &MoveWorkspaceToMonitorRight{AName: AName{Name: "MoveWorkspaceToMonitorRight"}} },
	"MoveWorkspaceToMonitorDown":  func() Action { return &MoveWorkspaceToMonitorDown{AName: AName{Name: "MoveWorkspaceToMonitorDown"}} },
	"MoveWorkspaceToMonitorUp":    func() Action { return &MoveWorkspaceToMonitorUp{AName: AName{Name: "MoveWorkspaceToMonitorUp"}} },
	"MoveWorkspaceToMonitorPrevious": func() Action {
		return &MoveWorkspaceToMonitorPrevious{AName: AName{Name: "MoveWorkspaceToMonitorPrevious"}}
	},
	"MoveWorkspaceToMonitorNext": func() Action { return &MoveWorkspaceToMonitorNext{AName: AName{Name: "MoveWorkspaceToMonitorNext"}} },
	"MoveWorkspaceToMonitor":     func() Action { return &MoveWorkspaceToMonitor{AName: AName{Name: "MoveWorkspaceToMonitor"}} },
	"ToggleDebugTint":            func() Action { return &ToggleDebugTint{AName: AName{Name: "ToggleDebugTint"}} },
	"DebugToggleOpaqueRegions":   func() Action { return &DebugToggleOpaqueRegions{AName: AName{Name: "DebugToggleOpaqueRegions"}} },
	"DebugToggleDamage":          func() Action { return &DebugToggleDamage{AName: AName{Name: "DebugToggleDamage"}} },
	"ToggleWindowFloating":       func() Action { return &ToggleWindowFloating{AName: AName{Name: "ToggleWindowFloating"}} },
	"MoveWindowToFloating":       func() Action { return &MoveWindowToFloating{AName: AName{Name: "MoveWindowToFloating"}} },
	"MoveWindowToTiling":         func() Action { return &MoveWindowToTiling{AName: AName{Name: "MoveWindowToTiling"}} },
	"FocusFloating":              func() Action { return &FocusFloating{AName: AName{Name: "FocusFloating"}} },
	"FocusTiling":                func() Action { return &FocusTiling{AName: AName{Name: "FocusTiling"}} },
	"SwitchFocusBetweenFloatingAndTiling": func() Action {
		return &SwitchFocusBetweenFloatingAndTiling{AName: AName{Name: "SwitchFocusBetweenFloatingAndTiling"}}
	},
	"MoveFloatingWindow":      func() Action { return &MoveFloatingWindow{AName: AName{Name: "MoveFloatingWindow"}} },
	"ToggleWindowRuleOpacity": func() Action { return &ToggleWindowRuleOpacity{AName: AName{Name: "ToggleWindowRuleOpacity"}} },
	"SetDynamicCastWindow":    func() Action { return &SetDynamicCastWindow{AName: AName{Name: "SetDynamicCastWindow"}} },
	"SetDynamicCastMonitor":   func() Action { return &SetDynamicCastMonitor{AName: AName{Name: "SetDynamicCastMonitor"}} },
	"ClearDynamicCastTarget":  func() Action { return &ClearDynamicCastTarget{AName: AName{Name: "ClearDynamicCastTarget"}} },
	"ToggleOverview":          func() Action { return &ToggleOverview{AName: AName{Name: "ToggleOverview"}} },
	"OpenOverview":            func() Action { return &OpenOverview{AName: AName{Name: "OpenOverview"}} },
	"CloseOverview":           func() Action { return &CloseOverview{AName: AName{Name: "CloseOverview"}} },
	"ToggleWindowUrgent":      func() Action { return &ToggleWindowUrgent{AName: AName{Name: "ToggleWindowUrgent"}} },
	"SetWindowUrgent":         func() Action { return &SetWindowUrgent{AName: AName{Name: "SetWindowUrgent"}} },
	"UnsetWindowUrgent":       func() Action { return &UnsetWindowUrgent{AName: AName{Name: "UnsetWindowUrgent"}} },
}
