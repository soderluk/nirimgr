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
	"CenterColumn":              func() Action { return &CenterColumn{AName: AName{Name: "CenterColumn"}} },
	"CenterVisibleColumns":      func() Action { return &CenterVisibleColumns{AName: AName{Name: "CenterVisibleColumns"}} },
	"CenterWindow":              func() Action { return &CenterWindow{AName: AName{Name: "CenterWindow"}} },
	"ClearDynamicCastTarget":    func() Action { return &ClearDynamicCastTarget{AName: AName{Name: "ClearDynamicCastTarget"}} },
	"CloseOverview":             func() Action { return &CloseOverview{AName: AName{Name: "CloseOverview"}} },
	"CloseWindow":               func() Action { return &CloseWindow{AName: AName{Name: "CloseWindow"}} },
	"ConsumeOrExpelWindowLeft":  func() Action { return &ConsumeOrExpelWindowLeft{AName: AName{Name: "ConsumeOrExpelWindowLeft"}} },
	"ConsumeOrExpelWindowRight": func() Action { return &ConsumeOrExpelWindowRight{AName: AName{Name: "ConsumeOrExpelWindowRight"}} },
	"ConsumeWindowIntoColumn":   func() Action { return &ConsumeWindowIntoColumn{AName: AName{Name: "ConsumeWindowIntoColumn"}} },
	"DebugToggleDamage":         func() Action { return &DebugToggleDamage{AName: AName{Name: "DebugToggleDamage"}} },
	"DebugToggleOpaqueRegions":  func() Action { return &DebugToggleOpaqueRegions{AName: AName{Name: "DebugToggleOpaqueRegions"}} },
	"DoScreenTransition":        func() Action { return &DoScreenTransition{AName: AName{Name: "DoScreenTransition"}} },
	"ExpandColumnToAvailableWidth": func() Action {
		return &ExpandColumnToAvailableWidth{AName: AName{Name: "ExpandColumnToAvailableWidth"}}
	},
	"ExpelWindowFromColumn":       func() Action { return &ExpelWindowFromColumn{AName: AName{Name: "ExpelWindowFromColumn"}} },
	"FocusColumn":                 func() Action { return &FocusColumn{AName: AName{Name: "FocusColumn"}} },
	"FocusColumnFirst":            func() Action { return &FocusColumnFirst{AName: AName{Name: "FocusColumnFirst"}} },
	"FocusColumnLast":             func() Action { return &FocusColumnLast{AName: AName{Name: "FocusColumnLast"}} },
	"FocusColumnLeft":             func() Action { return &FocusColumnLeft{AName: AName{Name: "FocusColumnLeft"}} },
	"FocusColumnLeftOrLast":       func() Action { return &FocusColumnLeftOrLast{AName: AName{Name: "FocusColumnLeftOrLast"}} },
	"FocusColumnOrMonitorLeft":    func() Action { return &FocusColumnOrMonitorLeft{AName: AName{Name: "FocusColumnOrMonitorLeft"}} },
	"FocusColumnOrMonitorRight":   func() Action { return &FocusColumnOrMonitorRight{AName: AName{Name: "FocusColumnOrMonitorRight"}} },
	"FocusColumnRight":            func() Action { return &FocusColumnRight{AName: AName{Name: "FocusColumnRight"}} },
	"FocusColumnRightOrFirst":     func() Action { return &FocusColumnRightOrFirst{AName: AName{Name: "FocusColumnRightOrFirst"}} },
	"FocusFloating":               func() Action { return &FocusFloating{AName: AName{Name: "FocusFloating"}} },
	"FocusMonitor":                func() Action { return &FocusMonitor{AName: AName{Name: "FocusMonitor"}} },
	"FocusMonitorDown":            func() Action { return &FocusMonitorDown{AName: AName{Name: "FocusMonitorDown"}} },
	"FocusMonitorLeft":            func() Action { return &FocusMonitorLeft{AName: AName{Name: "FocusMonitorLeft"}} },
	"FocusMonitorNext":            func() Action { return &FocusMonitorNext{AName: AName{Name: "FocusMonitorNext"}} },
	"FocusMonitorPrevious":        func() Action { return &FocusMonitorPrevious{AName: AName{Name: "FocusMonitorPrevious"}} },
	"FocusMonitorRight":           func() Action { return &FocusMonitorRight{AName: AName{Name: "FocusMonitorRight"}} },
	"FocusMonitorUp":              func() Action { return &FocusMonitorUp{AName: AName{Name: "FocusMonitorUp"}} },
	"FocusTiling":                 func() Action { return &FocusTiling{AName: AName{Name: "FocusTiling"}} },
	"FocusWindow":                 func() Action { return &FocusWindow{AName: AName{Name: "FocusWindow"}} },
	"FocusWindowBottom":           func() Action { return &FocusWindowBottom{AName: AName{Name: "FocusWindowBottom"}} },
	"FocusWindowDown":             func() Action { return &FocusWindowDown{AName: AName{Name: "FocusWindowDown"}} },
	"FocusWindowDownOrColumnLeft": func() Action { return &FocusWindowDownOrColumnLeft{AName: AName{Name: "FocusWindowDownOrColumnLeft"}} },
	"FocusWindowDownOrColumnRight": func() Action {
		return &FocusWindowDownOrColumnRight{AName: AName{Name: "FocusWindowDownOrColumnRight"}}
	},
	"FocusWindowDownOrTop":       func() Action { return &FocusWindowDownOrTop{AName: AName{Name: "FocusWindowDownOrTop"}} },
	"FocusWindowInColumn":        func() Action { return &FocusWindowInColumn{AName: AName{Name: "FocusWindowInColumn"}} },
	"FocusWindowOrMonitorDown":   func() Action { return &FocusWindowOrMonitorDown{AName: AName{Name: "FocusWindowOrMonitorDown"}} },
	"FocusWindowOrMonitorUp":     func() Action { return &FocusWindowOrMonitorUp{AName: AName{Name: "FocusWindowOrMonitorUp"}} },
	"FocusWindowOrWorkspaceDown": func() Action { return &FocusWindowOrWorkspaceDown{AName: AName{Name: "FocusWindowOrWorkspaceDown"}} },
	"FocusWindowOrWorkspaceUp":   func() Action { return &FocusWindowOrWorkspaceUp{AName: AName{Name: "FocusWindowOrWorkspaceUp"}} },
	"FocusWindowPrevious":        func() Action { return &FocusWindowPrevious{AName: AName{Name: "FocusWindowPrevious"}} },
	"FocusWindowTop":             func() Action { return &FocusWindowTop{AName: AName{Name: "FocusWindowTop"}} },
	"FocusWindowUp":              func() Action { return &FocusWindowUp{AName: AName{Name: "FocusWindowUp"}} },
	"FocusWindowUpOrBottom":      func() Action { return &FocusWindowUpOrBottom{AName: AName{Name: "FocusWindowUpOrBottom"}} },
	"FocusWindowUpOrColumnLeft":  func() Action { return &FocusWindowUpOrColumnLeft{AName: AName{Name: "FocusWindowUpOrColumnLeft"}} },
	"FocusWindowUpOrColumnRight": func() Action { return &FocusWindowUpOrColumnRight{AName: AName{Name: "FocusWindowUpOrColumnRight"}} },
	"FocusWorkspace":             func() Action { return &FocusWorkspace{AName: AName{Name: "FocusWorkspace"}} },
	"FocusWorkspaceDown":         func() Action { return &FocusWorkspaceDown{AName: AName{Name: "FocusWorkspaceDown"}} },
	"FocusWorkspacePrevious":     func() Action { return &FocusWorkspacePrevious{AName: AName{Name: "FocusWorkspacePrevious"}} },
	"FocusWorkspaceUp":           func() Action { return &FocusWorkspaceUp{AName: AName{Name: "FocusWorkspaceUp"}} },
	"FullscreenWindow":           func() Action { return &FullscreenWindow{AName: AName{Name: "FullscreenWindow"}} },
	"LoadConfigFile":             func() Action { return &LoadConfigFile{AName: AName{Name: "LoadConfigFile"}} },
	"MaximizeColumn":             func() Action { return &MaximizeColumn{AName: AName{Name: "MaximizeColumn"}} },
	"MaximizeWindowToEdges":      func() Action { return &MaximizeWindowToEdges{AName: AName{Name: "MaximizeWindowToEdges"}} },
	"MoveColumnLeft":             func() Action { return &MoveColumnLeft{AName: AName{Name: "MoveColumnLeft"}} },
	"MoveColumnLeftOrToMonitorLeft": func() Action {
		return &MoveColumnLeftOrToMonitorLeft{AName: AName{Name: "MoveColumnLeftOrToMonitorLeft"}}
	},
	"MoveColumnRight": func() Action { return &MoveColumnRight{AName: AName{Name: "MoveColumnRight"}} },
	"MoveColumnRightOrToMonitorRight": func() Action {
		return &MoveColumnRightOrToMonitorRight{AName: AName{Name: "MoveColumnRightOrToMonitorRight"}}
	},
	"MoveColumnToFirst":           func() Action { return &MoveColumnToFirst{AName: AName{Name: "MoveColumnToFirst"}} },
	"MoveColumnToIndex":           func() Action { return &MoveColumnToIndex{AName: AName{Name: "MoveColumnToIndex"}} },
	"MoveColumnToLast":            func() Action { return &MoveColumnToLast{AName: AName{Name: "MoveColumnToLast"}} },
	"MoveColumnToMonitor":         func() Action { return &MoveColumnToMonitor{AName: AName{Name: "MoveColumnToMonitor"}} },
	"MoveColumnToMonitorDown":     func() Action { return &MoveColumnToMonitorDown{AName: AName{Name: "MoveColumnToMonitorDown"}} },
	"MoveColumnToMonitorLeft":     func() Action { return &MoveColumnToMonitorLeft{AName: AName{Name: "MoveColumnToMonitorLeft"}} },
	"MoveColumnToMonitorNext":     func() Action { return &MoveColumnToMonitorNext{AName: AName{Name: "MoveColumnToMonitorNext"}} },
	"MoveColumnToMonitorPrevious": func() Action { return &MoveColumnToMonitorPrevious{AName: AName{Name: "MoveColumnToMonitorPrevious"}} },
	"MoveColumnToMonitorRight":    func() Action { return &MoveColumnToMonitorRight{AName: AName{Name: "MoveColumnToMonitorRight"}} },
	"MoveColumnToMonitorUp":       func() Action { return &MoveColumnToMonitorUp{AName: AName{Name: "MoveColumnToMonitorUp"}} },
	"MoveColumnToWorkspace":       func() Action { return &MoveColumnToWorkspace{AName: AName{Name: "MoveColumnToWorkspace"}} },
	"MoveColumnToWorkspaceDown":   func() Action { return &MoveColumnToWorkspaceDown{AName: AName{Name: "MoveColumnToWorkspaceDown"}} },
	"MoveColumnToWorkspaceUp":     func() Action { return &MoveColumnToWorkspaceUp{AName: AName{Name: "MoveColumnToWorkspaceUp"}} },
	"MoveFloatingWindow":          func() Action { return &MoveFloatingWindow{AName: AName{Name: "MoveFloatingWindow"}} },
	"MoveWindowDown":              func() Action { return &MoveWindowDown{AName: AName{Name: "MoveWindowDown"}} },
	"MoveWindowDownOrToWorkspaceDown": func() Action {
		return &MoveWindowDownOrToWorkspaceDown{AName: AName{Name: "MoveWindowDownOrToWorkspaceDown"}}
	},
	"MoveWindowToFloating":        func() Action { return &MoveWindowToFloating{AName: AName{Name: "MoveWindowToFloating"}} },
	"MoveWindowToMonitor":         func() Action { return &MoveWindowToMonitor{AName: AName{Name: "MoveWindowToMonitor"}} },
	"MoveWindowToMonitorDown":     func() Action { return &MoveWindowToMonitorDown{AName: AName{Name: "MoveWindowToMonitorDown"}} },
	"MoveWindowToMonitorLeft":     func() Action { return &MoveWindowToMonitorLeft{AName: AName{Name: "MoveWindowToMonitorLeft"}} },
	"MoveWindowToMonitorNext":     func() Action { return &MoveWindowToMonitorNext{AName: AName{Name: "MoveWindowToMonitorNext"}} },
	"MoveWindowToMonitorPrevious": func() Action { return &MoveWindowToMonitorPrevious{AName: AName{Name: "MoveWindowToMonitorPrevious"}} },
	"MoveWindowToMonitorRight":    func() Action { return &MoveWindowToMonitorRight{AName: AName{Name: "MoveWindowToMonitorRight"}} },
	"MoveWindowToMonitorUp":       func() Action { return &MoveWindowToMonitorUp{AName: AName{Name: "MoveWindowToMonitorUp"}} },
	"MoveWindowToTiling":          func() Action { return &MoveWindowToTiling{AName: AName{Name: "MoveWindowToTiling"}} },
	"MoveWindowToWorkspace":       func() Action { return &MoveWindowToWorkspace{AName: AName{Name: "MoveWindowToWorkspace"}} },
	"MoveWindowToWorkspaceDown":   func() Action { return &MoveWindowToWorkspaceDown{AName: AName{Name: "MoveWindowToWorkspaceDown"}} },
	"MoveWindowToWorkspaceUp":     func() Action { return &MoveWindowToWorkspaceUp{AName: AName{Name: "MoveWindowToWorkspaceUp"}} },
	"MoveWindowUp":                func() Action { return &MoveWindowUp{AName: AName{Name: "MoveWindowUp"}} },
	"MoveWindowUpOrToWorkspaceUp": func() Action { return &MoveWindowUpOrToWorkspaceUp{AName: AName{Name: "MoveWindowUpOrToWorkspaceUp"}} },
	"MoveWorkspaceDown":           func() Action { return &MoveWorkspaceDown{AName: AName{Name: "MoveWorkspaceDown"}} },
	"MoveWorkspaceToIndex":        func() Action { return &MoveWorkspaceToIndex{AName: AName{Name: "MoveWorkspaceToIndex"}} },
	"MoveWorkspaceToMonitor":      func() Action { return &MoveWorkspaceToMonitor{AName: AName{Name: "MoveWorkspaceToMonitor"}} },
	"MoveWorkspaceToMonitorDown":  func() Action { return &MoveWorkspaceToMonitorDown{AName: AName{Name: "MoveWorkspaceToMonitorDown"}} },
	"MoveWorkspaceToMonitorLeft":  func() Action { return &MoveWorkspaceToMonitorLeft{AName: AName{Name: "MoveWorkspaceToMonitorLeft"}} },
	"MoveWorkspaceToMonitorNext":  func() Action { return &MoveWorkspaceToMonitorNext{AName: AName{Name: "MoveWorkspaceToMonitorNext"}} },
	"MoveWorkspaceToMonitorPrevious": func() Action {
		return &MoveWorkspaceToMonitorPrevious{AName: AName{Name: "MoveWorkspaceToMonitorPrevious"}}
	},
	"MoveWorkspaceToMonitorRight": func() Action { return &MoveWorkspaceToMonitorRight{AName: AName{Name: "MoveWorkspaceToMonitorRight"}} },
	"MoveWorkspaceToMonitorUp":    func() Action { return &MoveWorkspaceToMonitorUp{AName: AName{Name: "MoveWorkspaceToMonitorUp"}} },
	"MoveWorkspaceUp":             func() Action { return &MoveWorkspaceUp{AName: AName{Name: "MoveWorkspaceUp"}} },
	"OpenOverview":                func() Action { return &OpenOverview{AName: AName{Name: "OpenOverview"}} },
	"PowerOffMonitors":            func() Action { return &PowerOffMonitors{AName: AName{Name: "PowerOffMonitors"}} },
	"PowerOnMonitors":             func() Action { return &PowerOnMonitors{AName: AName{Name: "PowerOnMonitors"}} },
	"Quit":                        func() Action { return &Quit{AName: AName{Name: "Quit"}} },
	"ResetWindowHeight":           func() Action { return &ResetWindowHeight{AName: AName{Name: "ResetWindowHeight"}} },
	"Screenshot":                  func() Action { return &Screenshot{AName: AName{Name: "Screenshot"}} },
	"ScreenshotScreen":            func() Action { return &ScreenshotScreen{AName: AName{Name: "ScreenshotScreen"}} },
	"ScreenshotWindow":            func() Action { return &ScreenshotWindow{AName: AName{Name: "ScreenshotWindow"}} },
	"SetColumnDisplay":            func() Action { return &SetColumnDisplay{AName: AName{Name: "SetColumnDisplay"}} },
	"SetColumnWidth":              func() Action { return &SetColumnWidth{AName: AName{Name: "SetColumnWidth"}} },
	"SetDynamicCastMonitor":       func() Action { return &SetDynamicCastMonitor{AName: AName{Name: "SetDynamicCastMonitor"}} },
	"SetDynamicCastWindow":        func() Action { return &SetDynamicCastWindow{AName: AName{Name: "SetDynamicCastWindow"}} },
	"SetWindowHeight":             func() Action { return &SetWindowHeight{AName: AName{Name: "SetWindowHeight"}} },
	"SetWindowUrgent":             func() Action { return &SetWindowUrgent{AName: AName{Name: "SetWindowUrgent"}} },
	"SetWindowWidth":              func() Action { return &SetWindowWidth{AName: AName{Name: "SetWindowWidth"}} },
	"SetWorkspaceName":            func() Action { return &SetWorkspaceName{AName: AName{Name: "SetWorkspaceName"}} },
	"ShowHotkeyOverlay":           func() Action { return &ShowHotkeyOverlay{AName: AName{Name: "ShowHotkeyOverlay"}} },
	"Spawn":                       func() Action { return &Spawn{AName: AName{Name: "Spawn"}} },
	"SpawnSh":                     func() Action { return &SpawnSh{AName: AName{Name: "SpawnSh"}} },
	"SwapWindowLeft":              func() Action { return &SwapWindowLeft{AName: AName{Name: "SwapWindowLeft"}} },
	"SwapWindowRight":             func() Action { return &SwapWindowRight{AName: AName{Name: "SwapWindowRight"}} },
	"SwitchFocusBetweenFloatingAndTiling": func() Action {
		return &SwitchFocusBetweenFloatingAndTiling{AName: AName{Name: "SwitchFocusBetweenFloatingAndTiling"}}
	},
	"SwitchLayout":                func() Action { return &SwitchLayout{AName: AName{Name: "SwitchLayout"}} },
	"SwitchPresetColumnWidth":     func() Action { return &SwitchPresetColumnWidth{AName: AName{Name: "SwitchPresetColumnWidth"}} },
	"SwitchPresetColumnWidthBack": func() Action { return &SwitchPresetColumnWidthBack{AName: AName{Name: "SwitchPresetColumnWidthBack"}} },
	"SwitchPresetWindowHeight":    func() Action { return &SwitchPresetWindowHeight{AName: AName{Name: "SwitchPresetWindowHeight"}} },
	"SwitchPresetWindowHeightBack": func() Action {
		return &SwitchPresetWindowHeightBack{AName: AName{Name: "SwitchPresetWindowHeightBack"}}
	},
	"SwitchPresetWindowWidth":     func() Action { return &SwitchPresetWindowWidth{AName: AName{Name: "SwitchPresetWindowWidth"}} },
	"SwitchPresetWindowWidthBack": func() Action { return &SwitchPresetWindowWidthBack{AName: AName{Name: "SwitchPresetWindowWidthBack"}} },
	"ToggleColumnTabbedDisplay": func() Action {
		return &ToggleColumnTabbedDisplay{AName: AName{Name: "ToggleColumnTabbedDisplay"}}
	},
	"ToggleDebugTint": func() Action { return &ToggleDebugTint{AName: AName{Name: "ToggleDebugTint"}} },
	"ToggleKeyboardShortcutsInhibit": func() Action {
		return &ToggleKeyboardShortcutsInhibit{AName: AName{Name: "ToggleKeyboardShortcutsInhibit"}}
	},
	"ToggleOverview":          func() Action { return &ToggleOverview{AName: AName{Name: "ToggleOverview"}} },
	"ToggleWindowFloating":    func() Action { return &ToggleWindowFloating{AName: AName{Name: "ToggleWindowFloating"}} },
	"ToggleWindowRuleOpacity": func() Action { return &ToggleWindowRuleOpacity{AName: AName{Name: "ToggleWindowRuleOpacity"}} },
	"ToggleWindowUrgent":      func() Action { return &ToggleWindowUrgent{AName: AName{Name: "ToggleWindowUrgent"}} },
	"ToggleWindowedFullscreen": func() Action {
		return &ToggleWindowedFullscreen{AName: AName{Name: "ToggleWindowedFullscreen"}}
	},
	"UnsetWindowUrgent":  func() Action { return &UnsetWindowUrgent{AName: AName{Name: "UnsetWindowUrgent"}} },
	"UnsetWorkspaceName": func() Action { return &UnsetWorkspaceName{AName: AName{Name: "UnsetWorkspaceName"}} },
}
