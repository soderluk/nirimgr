package actions

// ColumnDisplay sets the column display to either Normal (tiled) or Tabbed (tabs).
type ColumnDisplay struct {
	Normal int `json:"Normal,omitempty"`
	Tabbed int `json:"Tabbed,omitempty"`
}

// LayoutSwitchTarget defines the layout to switch to.
type LayoutSwitchTarget struct {
	// Next the next configured layout.
	Next int `json:"Next,omitempty"`
	// Prev the previous configured layout.
	Prev int `json:"Prev,omitempty"`
	// Index the specific layout by index.
	Index uint8 `json:"Index,omitempty"`
}

// PositionChange defines how we want to position a window.
type PositionChange struct {
	// SetFixed sets the position in logical pixels.
	SetFixed float64 `json:"SetFixed,omitempty"`
	// AdjustFixed adds or subtracts the current position in logical pixels.
	AdjustFixed float64 `json:"AdjustFixed,omitempty"`
}

// SizeChange defines how we want to change the size of a window.
type SizeChange struct {
	// SetFixed sets the size in logical pixels.
	SetFixed int32 `json:"SetFixed,omitempty"`
	// SetProportion sets the size as a proportion of the working area.
	SetProportion float64 `json:"SetProportion,omitempty"`
	// AdjustFixed adds or subtracts the current size in logical pixels.
	AdjustFixed int32 `json:"AdjustFixed,omitempty"`
	// AdjustProportion adds or subtracts the current size as a proportion of the working area.
	AdjustProportion float64 `json:"AdjustProportion,omitempty"`
}

// WorkspaceReferenceArg takes either the ID, Index or Name of the workspace.
type WorkspaceReferenceArg struct {
	// ID the ID of the workspace.
	ID uint64 `json:"Id,omitempty"`
	// Index the index of the workspace.
	Index uint8 `json:"Index,omitempty"`
	// Name the name of the workspace.
	Name string `json:"Name,omitempty"`
}
