package theme

import (
	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/style"
)

// Color schemes.
var (
	ButtonBackgroundColor = render.RGBA(200, 200, 200, 255)
	ButtonHoverColor      = render.RGBA(200, 255, 255, 255)
	ButtonOutlineColor    = render.Black

	BorderColorOffset = 40
)

// Theme is a collection of styles for various built-in widgets.
type Theme struct {
	Name    string
	Window  *style.Window
	Label   *style.Label
	Button  *style.Button
	Tooltip *style.Tooltip
}

// Default theme.
var Default = Theme{
	Name:    "Default",
	Label:   &style.DefaultLabel,
	Button:  &style.DefaultButton,
	Tooltip: &style.DefaultTooltip,
}

// DefaultFlat is a flat version of the default theme.
var DefaultFlat = Theme{
	Name: "DefaultFlat",
	Button: &style.Button{
		Background:      style.DefaultButton.Background,
		Foreground:      style.DefaultButton.Foreground,
		OutlineColor:    style.DefaultButton.OutlineColor,
		OutlineSize:     1,
		HoverBackground: style.DefaultButton.HoverBackground,
		HoverForeground: style.DefaultButton.HoverForeground,
		BorderStyle:     style.BorderSolid,
		BorderSize:      2,
	},
}

// DefaultDark is a dark version of the default theme.
var DefaultDark = Theme{
	Name: "DefaultDark",
	Label: &style.Label{
		Foreground: render.Grey,
	},
	Window: &style.Window{
		ActiveTitleBackground:   render.Red,
		ActiveTitleForeground:   render.White,
		InactiveTitleBackground: render.DarkGrey,
		InactiveTitleForeground: render.Grey,
		ActiveBackground:        render.Black,
		InactiveBackground:      render.Black,
	},
	Button: &style.Button{
		Background:      render.Black,
		Foreground:      render.Grey,
		OutlineColor:    render.DarkGrey,
		OutlineSize:     1,
		HoverBackground: render.Grey,
		BorderStyle:     style.BorderRaised,
		BorderSize:      2,
	},
	Tooltip: &style.Tooltip{
		Background: render.RGBA(60, 60, 60, 230),
		Foreground: render.Cyan,
	},
}
