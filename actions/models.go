package actions

// If more actions are added in Niri, we must define them here, and add them to the ActionRegistry.

// Action is the "base" interface for all the actions.
//
// NOTE: We have to use GetName, since the field is called Name.
type Action interface {
	GetName() string
}

// AName defines the name of the action.
type AName struct {
	Name string
}

// GetName returns the action name.
func (a AName) GetName() string {
	return a.Name
}

// Quit exits niri.
type Quit struct {
	AName
	// SkipConfirmation skips the "Press Enter to confirm" prompt.
	SkipConfirmation bool `json:"skip_confirmation"`
}

// PowerOffMonitors powers off all monitors via DPMS.
type PowerOffMonitors struct {
	AName
}

// PowerOnMonitors powers on all monitors via DPMS.
type PowerOnMonitors struct {
	AName
}

// Spawn spawns a command.
type Spawn struct {
	AName
	// Command the command to spawn.
	Command []string `json:"command"`
}

// SpawnSh spawns a command through the shell.
type SpawnSh struct {
	AName
	// Command the command to run
	Command string `json:"command"`
}

// DoScreenTransition does a screen transition.
type DoScreenTransition struct {
	AName
	// DelayMs the delay in ms for the screen to freeze before starting the transition.
	DelayMs uint16 `json:"delay_ms"`
}

// Screenshot opens the screenshot UI.
type Screenshot struct {
	AName
	// ShowPointer whether to show the pointer by default in the screenshot UI.
	ShowPointer bool `json:"show_pointer"`
}

// ScreenshotScreen screenshots the focused screen.
type ScreenshotScreen struct {
	AName
	// WriteToDisk writes the screenshot to disk in addition to putting it in the clipboard.
	WriteToDisk bool `json:"write_to_disk"`
	// ShowPointer whether to include the mouse pointer in the screenshot or not.
	ShowPointer bool `json:"show_pointer"`
}

// ScreenshotWindow screenshots a window.
type ScreenshotWindow struct {
	AName
	// ID the ID of the window to screenshot.
	ID uint64 `json:"id"`
	// WriteToDisk writes the screenshot to disk in addition to putting it in the clipboard.
	WriteToDisk bool `json:"write_to_disk"`
}

// ToggleKeyboardShortcutsInhibit enables or disables the keyboard shortcuts inhibitor (if any) for the focused surface.
type ToggleKeyboardShortcutsInhibit struct {
	AName
}

// CloseWindow closes a window.
type CloseWindow struct {
	AName
	// ID the ID of the window to close. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// FullscreenWindow toggles fullscreen on a window.
type FullscreenWindow struct {
	AName
	// ID the ID of the window to toggle. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// LoadConfigFile reloads the config file.
//
// Can be useful for scripts changing the config file, to avoid
// waiting the small duration for niri's config file watcher to
// notice the changes.
type LoadConfigFile struct {
	AName
}

// ToggleWindowedFullscreen toggles windowed (fake) fullscreen on a window.
type ToggleWindowedFullscreen struct {
	AName
	// ID the ID of the window to toggle. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// FocusWindow focuses a window by ID.
type FocusWindow struct {
	AName
	// ID the window ID to focus.
	ID uint64 `json:"id"`
}

// FocusWindowInColumn focuses a window in the focused column by index.
type FocusWindowInColumn struct {
	AName
	// Index the index of the window in the column. The index starts from 1 for the topmost window.
	Index uint8 `json:"index"`
}

// FocusWindowPrevious focuses the previously focused window.
type FocusWindowPrevious struct {
	AName
}

// FocusColumnLeft focuses the column to the left.
type FocusColumnLeft struct {
	AName
}

// FocusColumnRight focuses the column to the right.
type FocusColumnRight struct {
	AName
}

// FocusColumnFirst focuses the first column.
type FocusColumnFirst struct {
	AName
}

// FocusColumnLast focuses the last column.
type FocusColumnLast struct {
	AName
}

// FocusColumnRightOrFirst focuses the next column on the right, looping if at end.
type FocusColumnRightOrFirst struct {
	AName
}

// FocusColumnLeftOrLast focuses the next column on the left, looping if at start.
type FocusColumnLeftOrLast struct {
	AName
}

// FocusColumn focuses a column by index.
type FocusColumn struct {
	AName
	// Index the index of the column to focus. The index starts from 1 for the first column.
	Index uint `json:"index"`
}

// FocusWindowOrMonitorUp focuses the window or the monitor above.
type FocusWindowOrMonitorUp struct {
	AName
}

// FocusWindowOrMonitorDown focuses the window or the monitor below.
type FocusWindowOrMonitorDown struct {
	AName
}

// FocusColumnOrMonitorLeft focuses the column or monitor to the left.
type FocusColumnOrMonitorLeft struct {
	AName
}

// FocusColumnOrMonitorRight focuses the column or monitor to the right.
type FocusColumnOrMonitorRight struct {
	AName
}

// FocusWindowDown focuses the window below.
type FocusWindowDown struct {
	AName
}

// FocusWindowUp focuses the window above.
type FocusWindowUp struct {
	AName
}

// FocusWindowDownOrColumnLeft focuses the window below or the column to the left.
type FocusWindowDownOrColumnLeft struct {
	AName
}

// FocusWindowDownOrColumnRight focuses the window above or the column to the right.
type FocusWindowDownOrColumnRight struct {
	AName
}

// FocusWindowUpOrColumnLeft focuses the window above or the column to the left.
type FocusWindowUpOrColumnLeft struct {
	AName
}

// FocusWindowUpOrColumnRight focuses the window above or the column to the right.
type FocusWindowUpOrColumnRight struct {
	AName
}

// FocusWindowOrWorkspaceDown focuses the window or the workspace below.
type FocusWindowOrWorkspaceDown struct {
	AName
}

// FocusWindowOrWorkspaceUp focuses the window or the workspace above.
type FocusWindowOrWorkspaceUp struct {
	AName
}

// FocusWindowTop focuses the topmost window.
type FocusWindowTop struct {
	AName
}

// FocusWindowBottom focuses the bottommost window.
type FocusWindowBottom struct {
	AName
}

// FocusWindowDownOrTop focuses the window below or the topmost window.
type FocusWindowDownOrTop struct {
	AName
}

// FocusWindowUpOrBottom focuses the window above or the bottommost window.
type FocusWindowUpOrBottom struct {
	AName
}

// MoveColumnLeft moves the focused column to the left.
type MoveColumnLeft struct {
	AName
}

// MoveColumnRight moves the focused column to the right.
type MoveColumnRight struct {
	AName
}

// MoveColumnToFirst moves the focused column to the start of the workspace.
type MoveColumnToFirst struct {
	AName
}

// MoveColumnToLast moves the focused column to the end of the workspace.
type MoveColumnToLast struct {
	AName
}

// MoveColumnLeftOrToMonitorLeft moves the focused column to the left, or to the monitor to the left.
type MoveColumnLeftOrToMonitorLeft struct {
	AName
}

// MoveColumnRightOrToMonitorRight moves the focused column to the right, or to the monitor to the right.
type MoveColumnRightOrToMonitorRight struct {
	AName
}

// MoveColumnToIndex moves the focused column to a specific index on its workspace.
type MoveColumnToIndex struct {
	AName
	// Index is the new index for the column. The index starts from 1 for the first column.
	Index uint `json:"index"`
}

// MoveWindowDown moves the focused window down in a column.
type MoveWindowDown struct {
	AName
}

// MoveWindowUp moves the focused window up in a column.
type MoveWindowUp struct {
	AName
}

// MoveWindowDownOrToWorkspaceDown moves the focused window down in a column or the workspace below.
type MoveWindowDownOrToWorkspaceDown struct {
	AName
}

// MoveWindowUpOrToWorkspaceUp moves the focused window up in a column or to the workspace above.
type MoveWindowUpOrToWorkspaceUp struct {
	AName
}

// ConsumeOrExpelWindowLeft consumes or expels a window left.
type ConsumeOrExpelWindowLeft struct {
	AName
	// ID the ID of the window to consume or expel. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// ConsumeOrExpelWindowRight consumes or expels a window right.
type ConsumeOrExpelWindowRight struct {
	AName
	// ID the ID of the window to consume or expel. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// ConsumeWindowIntoColumn consumes the window to the right into the focused column.
type ConsumeWindowIntoColumn struct {
	AName
}

// ExpelWindowFromColumn expels the focused window from the column.
type ExpelWindowFromColumn struct {
	AName
}

// SwapWindowRight swaps the focused window with the one to the right.
type SwapWindowRight struct {
	AName
}

// SwapWindowLeft swaps the focused window with the one to the left.
type SwapWindowLeft struct {
	AName
}

// ToggleColumnTabbedDisplay toggles the focused column between normal and tabbed display.
type ToggleColumnTabbedDisplay struct {
	AName
}

// SetColumnDisplay sets the display mode of the focused column.
type SetColumnDisplay struct {
	AName
	// Display display mode to set.
	Display ColumnDisplay `json:"display"`
}

// CenterColumn centers the focused column on the screen.
type CenterColumn struct {
	AName
}

// CenterWindow centers a window on the screen.
type CenterWindow struct {
	AName
	// ID the ID of the window to center. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// CenterVisibleColumns centers all fully visible columns on the screen.
type CenterVisibleColumns struct {
	AName
}

// FocusWorkspaceDown focuses the workspace below.
type FocusWorkspaceDown struct {
	AName
}

// FocusWorkspaceUp focuses the workspace above.
type FocusWorkspaceUp struct {
	AName
}

// FocusWorkspace focuses a workspace by reference (id, index or name).
type FocusWorkspace struct {
	AName
	// Reference the reference (id, index or name) of the workspace to focus.
	Reference WorkspaceReferenceArg `json:"reference"`
}

// FocusWorkspacePrevious focuses the previous workspace.
type FocusWorkspacePrevious struct {
	AName
}

// MoveWindowToWorkspaceDown moves the focused window to the workspace below.
type MoveWindowToWorkspaceDown struct {
	AName
}

// MoveWindowToWorkspaceUp moves the focused window to the workspace above.
type MoveWindowToWorkspaceUp struct {
	AName
}

// MoveWindowToWorkspace moves a window to a workspace.
type MoveWindowToWorkspace struct {
	AName
	// WindowID the ID of the window to move. If omitted, uses the focused window.
	WindowID uint64 `json:"window_id"`
	// Reference the reference (id, index or name) of the workspace to move the window to.
	Reference WorkspaceReferenceArg `json:"reference"`
	// Focus follows the moved window.
	//
	// If true (default) and the window to move is focused, the focus will follow the window to the new workspace.
	// If false, the focus will remain on the original workspace.
	Focus bool `json:"focus"`
}

// MoveColumnToWorkspaceDown moves the focused column to the workspace below.
type MoveColumnToWorkspaceDown struct {
	AName
	// Focus follows the target workspace.
	//
	// If true (default), the focus will follow the column to the new workspace.
	// If false, the focus will remain on the original workspace.
	Focus bool `json:"focus"`
}

// MoveColumnToWorkspaceUp moves the focused column to the workspace above.
type MoveColumnToWorkspaceUp struct {
	AName
	// Focus follows the target workspace.
	//
	// If true (default), the focus will follow the column to the new workspace.
	// If false, the focus will remain on the original workspace.
	Focus bool `json:"focus"`
}

// MoveColumnToWorkspace moves the focused column to a workspace by reference (id, index or name).
type MoveColumnToWorkspace struct {
	AName
	// Reference the reference (id, index or name) of the workspace to move the column to.
	Reference WorkspaceReferenceArg `json:"reference"`
	// Focus follows the target workspace.
	//
	// If true (default), the focus will follow the column to the new workspace.
	// If false, the focus will remain on the original workspace.
	Focus bool `json:"focus"`
}

// MoveWorkspaceDown moves the focused workspace below.
type MoveWorkspaceDown struct {
	AName
}

// MoveWorkspaceUp moves the focused workspace above.
type MoveWorkspaceUp struct {
	AName
}

// MoveWorkspaceToIndex moves the focused workspace to a specific index on its monitor.
type MoveWorkspaceToIndex struct {
	AName
	// Index the new index for the workspace.
	Index uint `json:"index"`
	// Reference the reference (id, index or name) of the workspace to move. If omitted, uses the focused workspace.
	Reference WorkspaceReferenceArg `json:"reference"`
}

// SetWorkspaceName sets the name of a workspace.
type SetWorkspaceName struct {
	AName
	// Name the new name of the workspace.
	Name string `json:"name"`
	// Workspace the reference (id, index or name) of the workspace to name. If omitted, uses the focused workspace.
	Workspace WorkspaceReferenceArg `json:"workspace"`
}

// UnsetWorkspaceName unsets the name of a workspace.
type UnsetWorkspaceName struct {
	AName
	// Reference the reference (id, index or name) of the workspace to unname. If omitted, uses the focused workspace.
	Reference WorkspaceReferenceArg `json:"reference"`
}

// FocusMonitorLeft focuses the monitor to the left.
type FocusMonitorLeft struct {
	AName
}

// FocusMonitorRight focuses the monitor to the right.
type FocusMonitorRight struct {
	AName
}

// FocusMonitorDown focuses the monitor below.
type FocusMonitorDown struct {
	AName
}

// FocusMonitorUp focuses the monitor above.
type FocusMonitorUp struct {
	AName
}

// FocusMonitorPrevious focuses the previous monitor.
type FocusMonitorPrevious struct {
	AName
}

// FocusMonitorNext focuses the next monitor.
type FocusMonitorNext struct {
	AName
}

// FocusMonitor focuses a monitor by name.
type FocusMonitor struct {
	AName
	// Output the name of the output to focus.
	Output string `json:"output"`
}

// MoveWindowToMonitorLeft moves the focused window to the monitor to the left.
type MoveWindowToMonitorLeft struct {
	AName
}

// MoveWindowToMonitorRight moves the focused window to the monitor to the right.
type MoveWindowToMonitorRight struct {
	AName
}

// MoveWindowToMonitorDown moves the focused window to the monitor below.
type MoveWindowToMonitorDown struct {
	AName
}

// MoveWindowToMonitorUp moves the focused window to the monitor above.
type MoveWindowToMonitorUp struct {
	AName
}

// MoveWindowToMonitorPrevious moves the focused window to the previous monitor.
type MoveWindowToMonitorPrevious struct {
	AName
}

// MoveWindowToMonitorNext moves the focused window to the next monitor.
type MoveWindowToMonitorNext struct {
	AName
}

// MoveWindowToMonitor moves a window to a specific monitor.
type MoveWindowToMonitor struct {
	AName
	// ID the ID of the window to move. If omitted, uses the focused window.
	ID uint64 `json:"id"`
	// Output the target output name.
	Output string `json:"output"`
}

// MoveColumnToMonitorLeft moves the focused column to the monitor to the left.
type MoveColumnToMonitorLeft struct {
	AName
}

// MoveColumnToMonitorRight moves the focused column to the monitor to the right.
type MoveColumnToMonitorRight struct {
	AName
}

// MoveColumnToMonitorDown moves the focused column to the monitor below.
type MoveColumnToMonitorDown struct {
	AName
}

// MoveColumnToMonitorUp moves the focused column to the monitor above.
type MoveColumnToMonitorUp struct {
	AName
}

// MoveColumnToMonitorPrevious moves the focused column to the previous monitor.
type MoveColumnToMonitorPrevious struct {
	AName
}

// MoveColumnToMonitorNext moves the focused column to the next monitor.
type MoveColumnToMonitorNext struct {
	AName
}

// MoveColumnToMonitor moves the focused column to a specific monitor.
type MoveColumnToMonitor struct {
	AName
	// Output the target output name.
	Output string `json:"output"`
}

// SetWindowWidth changes the width of a window.
type SetWindowWidth struct {
	AName
	// ID the ID of the window to change the width for. If omitted, uses the focused window.
	ID uint64 `json:"id"`
	// Change tells how to change the width.
	Change SizeChange `json:"change"`
}

// SetWindowHeight changes the height of a window.
type SetWindowHeight struct {
	AName
	// ID the ID of the window to change the height for. If omitted, uses the focused window.
	ID uint64 `json:"id"`
	// Change tells how to change the height.
	Change SizeChange `json:"change"`
}

// ResetWindowHeight resets the height of a window back to automatic.
type ResetWindowHeight struct {
	AName
	// ID the ID of the window to reset the height for. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// SwitchPresetColumnWidth switches between preset column widths.
type SwitchPresetColumnWidth struct {
	AName
}

// SwitchPresetColumnWidthBack switches between preset column widths backwards.
type SwitchPresetColumnWidthBack struct {
	AName
}

// SwitchPresetWindowWidth switches between preset window widths.
type SwitchPresetWindowWidth struct {
	AName
	// ID the ID of the window to switch the width for. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// SwitchPresetWindowWidthBack switches between preset window widths backwards.
type SwitchPresetWindowWidthBack struct {
	AName
	// ID the ID of the window whose width to switch. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// SwitchPresetWindowHeight switches between preset window heights.
type SwitchPresetWindowHeight struct {
	AName
	// ID the ID of the window to switch the height for. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// SwitchPresetWindowHeightBack switches between preset window heights backwards.
type SwitchPresetWindowHeightBack struct {
	AName
	// ID the ID of the window whose height to switch. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// MaximizeColumn toggles the maximized state of the focused column.
type MaximizeColumn struct {
	AName
}

// MaximizeWindowToEdges toggles the maximized-to-edges state of the focused window.
type MaximizeWindowToEdges struct {
	AName

	// ID of the window to maximize.
	ID uint64 `json:"id,omitempty"`
}

// SetColumnWidth changes the width of the focused column.
type SetColumnWidth struct {
	AName
	// Change tells how to change the width.
	Change SizeChange `json:"change"`
}

// ExpandColumnToAvailableWidth expands the focused column to space not taken up by other fully visible columns.
type ExpandColumnToAvailableWidth struct {
	AName
}

// SwitchLayout switches between keyboard layouts.
type SwitchLayout struct {
	AName
	// Layout the layout to switch to.
	Layout LayoutSwitchTarget `json:"layout"`
}

// ShowHotkeyOverlay shows the hotkey overlay.
type ShowHotkeyOverlay struct {
	AName
}

// MoveWorkspaceToMonitorLeft moves the focused workspace to the monitor to the left.
type MoveWorkspaceToMonitorLeft struct {
	AName
}

// MoveWorkspaceToMonitorRight moves the focused workspace to the monitor to the right.
type MoveWorkspaceToMonitorRight struct {
	AName
}

// MoveWorkspaceToMonitorDown moves the focused workspace to the monitor below.
type MoveWorkspaceToMonitorDown struct {
	AName
}

// MoveWorkspaceToMonitorUp moves the focused workspace to the monitor above.
type MoveWorkspaceToMonitorUp struct {
	AName
}

// MoveWorkspaceToMonitorPrevious moves the focused workspace to the previous monitor.
type MoveWorkspaceToMonitorPrevious struct {
	AName
}

// MoveWorkspaceToMonitorNext moves the focused workspace to the next monitor.
type MoveWorkspaceToMonitorNext struct {
	AName
}

// MoveWorkspaceToMonitor moves a workspace to a specific monitor.
type MoveWorkspaceToMonitor struct {
	AName
	// Output the target output name.
	Output string `json:"output"`
	// Reference the reference (id, index or name) of the workspace to move. If omitted, uses the focused workspace.
	Reference WorkspaceReferenceArg `json:"reference"`
}

// ToggleDebugTint toggles a debug tint on windows.
type ToggleDebugTint struct {
	AName
}

// DebugToggleOpaqueRegions toggles visualization of render element opaque regions.
type DebugToggleOpaqueRegions struct {
	AName
}

// DebugToggleDamage toggles visualization of output damage.
type DebugToggleDamage struct {
	AName
}

// ToggleWindowFloating toggles the focused window between floating and tiling layout.
type ToggleWindowFloating struct {
	AName
	// ID the ID of the window to toggle. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// MoveWindowToFloating moves a window to the floating layout.
type MoveWindowToFloating struct {
	AName
	// ID the ID of the window to toggle. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// MoveWindowToTiling moves a window to the tiling layout.
type MoveWindowToTiling struct {
	AName
	// ID the ID of the window to toggle. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// FocusFloating switches focus to the floating layout.
type FocusFloating struct {
	AName
}

// FocusTiling switches focus to the tiling layout.
type FocusTiling struct {
	AName
}

// SwitchFocusBetweenFloatingAndTiling toggles the focus between floating and tiling layout.
type SwitchFocusBetweenFloatingAndTiling struct {
	AName
}

// MoveFloatingWindow moves a floating window or screen.
type MoveFloatingWindow struct {
	AName
	// ID the ID of the window to move. If omitted, uses the focused window.
	ID uint64 `json:"id"`
	// X tells how to change the x position.
	X PositionChange `json:"x"`
	// Y tells how to change the y position.
	Y PositionChange `json:"y"`
}

// ToggleWindowRuleOpacity toggles the opacity of a window.
type ToggleWindowRuleOpacity struct {
	AName
	// ID the ID of the window to toggle. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// SetDynamicCastWindow sets the dynamic cast target to a window.
type SetDynamicCastWindow struct {
	AName
	// ID the ID of the window to target. If omitted, uses the focused window.
	ID uint64 `json:"id"`
}

// SetDynamicCastMonitor sets the dynamic cast target to a monitor.
type SetDynamicCastMonitor struct {
	AName
	// Output the name of the output to target. If omitted, uses the focused output.
	Output string `json:"output"`
}

// ClearDynamicCastTarget clears the dynamic cast target, making it show nothing.
type ClearDynamicCastTarget struct {
	AName
}

// ToggleOverview toggles the overview.
type ToggleOverview struct {
	AName
}

// OpenOverview open the overview.
type OpenOverview struct {
	AName
}

// CloseOverview closes the overview.
type CloseOverview struct {
	AName
}

// ToggleWindowUrgent toggles the urgent status of a window.
type ToggleWindowUrgent struct {
	AName
	// ID the ID of the window to toggle.
	ID uint64 `json:"id"`
}

// SetWindowUrgent sets the urgent status of a window.
type SetWindowUrgent struct {
	AName
	// ID the ID of the window to set urgent.
	ID uint64 `json:"id"`
}

// UnsetWindowUrgent unsets the urgent status of a window.
type UnsetWindowUrgent struct {
	AName
	// ID the ID of the window to unset urgent.
	ID uint64 `json:"id"`
}
