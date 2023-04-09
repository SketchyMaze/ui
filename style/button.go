// Package style provides style definitions for UI components.
package style

import "git.kirsle.net/go/render"

// Default styles for widgets without a theme.
var (
	DefaultWindow = Window{
		ActiveTitleBackground:   render.Blue,
		ActiveTitleForeground:   render.White,
		InactiveTitleBackground: render.DarkGrey,
		InactiveTitleForeground: render.Grey,
		ActiveBackground:        render.Grey,
		InactiveBackground:      render.Grey,
	}

	DefaultLabel = Label{
		Background: render.Invisible,
		Foreground: render.Black,
	}

	DefaultButton = Button{
		Background:      render.RGBA(200, 200, 200, 255),
		Foreground:      render.Black,
		OutlineColor:    render.Black,
		OutlineSize:     1,
		HoverBackground: render.RGBA(200, 255, 255, 255),
		HoverForeground: render.Black,
		BorderStyle:     BorderRaised,
		BorderSize:      2,
	}

	DefaultListBox = ListBox{
		Background:         render.White,
		Foreground:         render.Black,
		HoverBackground:    render.Cyan,
		HoverForeground:    render.Orange,
		SelectedBackground: render.Blue,
		SelectedForeground: render.White,
		BorderStyle:        BorderSunken,
		// BorderColor: render.RGBA(200, 200, 200, 255),
		BorderSize: 2,
	}

	DefaultTooltip = Tooltip{
		Background: render.RGBA(0, 0, 0, 230),
		Foreground: render.White,
	}
)

// Window style configuration.
type Window struct {
	ActiveTitleBackground   render.Color
	ActiveTitleForeground   render.Color
	ActiveBackground        render.Color
	InactiveTitleBackground render.Color
	InactiveTitleForeground render.Color
	InactiveBackground      render.Color
}

// Label style configuration.
type Label struct {
	Background render.Color
	Foreground render.Color
}

// Button style configuration.
type Button struct {
	Background      render.Color
	Foreground      render.Color // Labels only
	OutlineColor    render.Color
	OutlineSize     int
	HoverBackground render.Color
	HoverForeground render.Color
	BorderStyle     BorderStyle
	BorderSize      int
}

// Tooltip style configuration.
type Tooltip struct {
	Background render.Color
	Foreground render.Color
}

// ListBox style configuration.
type ListBox struct {
	Background         render.Color
	Foreground         render.Color // Labels only
	SelectedBackground render.Color
	SelectedForeground render.Color
	HoverBackground    render.Color
	HoverForeground    render.Color
	BorderStyle        BorderStyle
	BorderSize         int
}
