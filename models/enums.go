// Contains enums for niri models.
//
// These are not really structs, nor are they actions or events either,
// but can be used in models.
package models

// ModeToSet is the output mode to set.
type ModeToSet struct {
	// Automatic tells that niri will pick the mode automatically.
	Automatic string `json:"Automatic,omitempty"`
	// Specific tells that niri should pick a specific mode.
	Specific ConfiguredMode `json:"Specific,omitempty"`
}

// OutputAction is the output actions that niri can perform.
type OutputAction struct {
	// Off tells niri to turn off the output.
	Off string `json:"Off,omitempty"`
	// On tells niri to turn on the output.
	On string `json:"On,omitempty"`
	// Mode tells niri which output mode to set.
	Mode struct {
		// Mode is the mode to set, or "auto" for automatic selection.
		//
		// Run niri msg outputs to see the available modes.
		Mode ModeToSet `json:"mode"`
	} `json:"Mode"`
	// Scale tells niri which output scale to set.
	Scale struct {
		// Scale is the scale factor to set, or "auto" for automatic selection.
		Scale ScaleToSet `json:"scale"`
	} `json:"Scale"`
	// Transform tells niri which output transform to set.
	Transform struct {
		// Transform is the transform to set, counter-clockwise.
		Transform Transform `json:"transform"`
	} `json:"Transform"`
	// Position tells niri which output position to set.
	Position struct {
		// Position is the position to set, or "auto" for automatic selection.
		Position PositionToSet `json:"position"`
	} `json:"Position"`
	// Vrr tells niri which variable refresh rate to set.
	Vrr struct {
		// Vrr is the variable refresh rate mode to set.
		Vrr VrrToSet `json:"vrr"`
	} `json:"Vrr"`
}

// OutputConfigChanged is the output configuration change result.
type OutputConfigChanged struct {
	// Applied tells if the target output was connected and the change was applied.
	Applied string `json:"Applied"`
	// OutputWasMissing tells if the target output was not found, the change will be applied when it's connected.
	OutputWasMissing string `json:"OutputWasMissing"`
}

// PositionToSet is the output position to set.
type PositionToSet struct {
	// Automatic tells niri to position the output automatically.
	Automatic string `json:"Automatic,omitempty"`
	// Specific tells niri to use a specific position.
	Specific ConfiguredPosition `json:"Specific,omitempty"`
}

// ScaleToSet is the output scale to set.
type ScaleToSet struct {
	// Automatic tells niri to pick the scale automatically.
	Automatic string `json:"Automatic,omitempty"`
	// Specific tells niri to set a specific scale.
	Specific float64 `json:"Specific,omitempty"`
}

// Layer is the layer-shell layer.
type Layer struct {
	// Background is the background layer.
	Background string `json:"Background,omitempty"`
	// Bottom is the bottom layer.
	Bottom string `json:"Bottom,omitempty"`
	// Top is the top layer.
	Top string `json:"Top,omitempty"`
	// Overlay is the overlay layer.
	Overlay string `json:"Overlay,omitempty"`
}

// LayerSurfaceKeyboardInteractivity is the keyboard interactivity modes for a layer-shell surface.
type LayerSurfaceKeyboardInteractivity struct {
	// None tells that the surface cannot receive keyboard focus.
	None string `json:"None,omitempty"`
	// Exclusive tells that the surface receives keyboard focus whenever possible.
	Exclusive string `json:"Exclusive,omitempty"`
	// OnDemand tells that the surface receives keyboard focus on demand, e.g. when clicked.
	OnDemand string `json:"OnDemand,omitempty"`
}

// Transform is the output transformation, which goes counter-clockwise.
type Transform struct {
	// Normal is untransformed.
	Normal string `json:"Normal,omitempty"`
	// Rotate90 is rotated by 90 deg.
	Rotate90 string `json:"_90,omitempty"`
	// Rotate180 is rotated by 180 deg.
	Rotate180 string `json:"_180,omitempty"`
	// Rotate270 is rotated by 270 deg.
	Rotate270 string `json:"_270,omitempty"`
	// Flipped is flipped horizontally.
	Flipped string `json:"Flipped,omitempty"`
	// Flipped90 is rotated by 90 deg and flipped horizontally.
	Flipped90 string `json:"Flipped90,omitempty"`
	// Flipped180 is flipped vertically.
	Flipped180 string `json:"Flipped180,omitempty"`
	// Flipped270 is rotated by 270 deg and flipped horizontally.
	Flipped270 string `json:"Flipped270,omitempty"`
}
