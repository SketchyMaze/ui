package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
)

// Window is a frame with a title bar.
type Window struct {
	BaseWidget
	Title  string
	Active bool

	// Private widgets.
	body     *Frame
	titleBar *Label
	content  *Frame
}

// NewWindow creates a new window.
func NewWindow(title string) *Window {
	w := &Window{
		Title: title,
		body:  NewFrame("body:" + title),
	}
	w.IDFunc(func() string {
		return fmt.Sprintf("Window<%s>",
			w.Title,
		)
	})

	w.body.Configure(Config{
		Background:  render.Grey,
		BorderSize:  2,
		BorderStyle: BorderRaised,
	})

	// Title bar widget.
	titleBar := NewLabel(Label{
		TextVariable: &w.Title,
		Font: render.Text{
			Color:   render.White,
			Size:    10,
			Stroke:  render.DarkBlue,
			Padding: 2,
		},
	})
	titleBar.Configure(Config{
		Background: render.Blue,
	})
	w.body.Pack(titleBar, Pack{
		Side: N,
		Fill: true,
	})
	w.titleBar = titleBar

	// Window content frame.
	content := NewFrame("content:" + title)
	content.Configure(Config{
		Background: render.Grey,
	})
	w.body.Pack(content, Pack{
		Side: N,
		Fill: true,
	})
	w.content = content

	// Set up parent/child relationships
	w.body.SetParent(w)

	return w
}

// Children returns the window's child widgets.
func (w *Window) Children() []Widget {
	return []Widget{
		w.body,
	}
}

// TitleBar returns the title bar widget.
func (w *Window) TitleBar() *Label {
	return w.titleBar
}

// Configure the widget. Color and style changes are passed down to the inner
// content frame of the window.
func (w *Window) Configure(C Config) {
	w.BaseWidget.Configure(C)
	w.body.Configure(C)

	// Don't pass dimensions down any further than the body.
	C.Width = 0
	C.Height = 0
	w.content.Configure(C)
}

// ConfigureTitle configures the title bar widget.
func (w *Window) ConfigureTitle(C Config) {
	w.titleBar.Configure(C)
}

// Compute the window.
func (w *Window) Compute(e render.Engine) {
	w.body.Compute(e)
}

// Present the window.
func (w *Window) Present(e render.Engine, P render.Point) {
	w.body.Present(e, P)
}

// Pack a widget into the window's frame.
func (w *Window) Pack(child Widget, config ...Pack) {
	w.content.Pack(child, config...)
}
