// Package style provides style definitions for UI components.
package style

// BorderStyle options for widget.SetBorderStyle()
type BorderStyle string

// Styles for a widget border.
const (
	BorderNone   BorderStyle = ""
	BorderSolid  BorderStyle = "solid"
	BorderRaised             = "raised"
	BorderSunken             = "sunken"
)
