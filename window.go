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

	// Title bar colors. Sensible defaults are chosen in NewWindow but you
	// may customize after the fact.
	ActiveTitleBackground   render.Color
	ActiveTitleForeground   render.Color
	InactiveTitleBackground render.Color
	InactiveTitleForeground render.Color

	// Private widgets.
	body       *Frame
	titleBar   *Frame
	titleLabel *Label
	content    *Frame

	// Window manager controls.
	dragging      bool
	startDragAt   render.Point // cursor position when drag began
	dragOrigPoint render.Point // original position of window at drag start
	focused       bool
	managed       bool // window is managed by Supervisor
}

// NewWindow creates a new window.
func NewWindow(title string) *Window {
	w := &Window{
		Title: title,
		body:  NewFrame("body:" + title),

		// Default title bar colors.
		ActiveTitleBackground:   render.Blue,
		ActiveTitleForeground:   render.White,
		InactiveTitleBackground: render.Grey,
		InactiveTitleForeground: render.Black,
	}
	w.IDFunc(func() string {
		return fmt.Sprintf("Window<%s %+v>",
			w.Title, w.focused,
		)
	})

	w.body.Configure(Config{
		Background:  render.Grey,
		BorderSize:  2,
		BorderStyle: BorderRaised,
	})

	// Title bar widget.
	titleBar, titleLabel := w.setupTitleBar()
	w.body.Pack(titleBar, Pack{
		Side: N,
		Fill: true,
	})
	w.titleBar = titleBar
	w.titleLabel = titleLabel

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

// setupTitlebar creates the title bar frame of the window.
func (w *Window) setupTitleBar() (*Frame, *Label) {
	frame := NewFrame("Titlebar for Windows: " + w.Title)
	frame.Configure(Config{
		Background: w.ActiveTitleBackground,
	})

	label := NewLabel(Label{
		TextVariable: &w.Title,
		Font: render.Text{
			Color:   w.ActiveTitleForeground,
			Size:    10,
			Stroke:  w.ActiveTitleBackground.Darken(40),
			Padding: 2,
		},
	})
	frame.Pack(label, Pack{
		Side: W,
	})

	return frame, label
}

// Supervise enables the window to be dragged around by its title bar by
// adding its relevant event hooks to your Supervisor.
func (w *Window) Supervise(s *Supervisor) {
	// Add a click handler to the title bar to enable dragging.
	w.titleBar.Handle(MouseDown, func(ed EventData) error {
		w.startDragAt = ed.Point
		w.dragOrigPoint = w.Point()
		fmt.Printf("Clicked at %s   window at %s!\n", ed.Point, w.dragOrigPoint)

		s.DragStartWidget(w)
		return nil
	})

	// Clicking anywhere in the window focuses the window.
	w.Handle(MouseDown, func(ed EventData) error {
		s.FocusWindow(w)
		fmt.Printf("%s handles click event\n", w)
		return nil
	})
	w.Handle(Click, func(ed EventData) error {
		return nil
	})

	// Window as a whole receives DragMove events while being dragged.
	w.Handle(DragMove, func(ed EventData) error {
		// Get the delta of movement from where we began.
		delta := w.startDragAt.Compare(ed.Point)
		if delta != render.Origin {
			fmt.Printf("    Dragged to: %s   Delta: %s\n", ed.Point, delta)
			moveTo := w.dragOrigPoint
			moveTo.Add(delta)
			w.MoveTo(moveTo)
		}
		return nil
	})

	// Add the title bar to the supervisor.
	s.Add(w.titleBar)
	s.Add(w)

	// Add the window to the focus list of the supervisor.
	s.addWindow(w)
}

// Focused returns whether the window is focused.
func (w *Window) Focused() bool {
	return w.focused
}

// SetFocus sets the window's focus value. Note: if you're using the Supervisor
// to manage the windows, do NOT call this method -- window focus is managed
// by the Supervisor.
func (w *Window) SetFocus(v bool) {
	w.focused = v

	// Update the title bar colors.
	var (
		bg = w.ActiveTitleBackground
		fg = w.ActiveTitleForeground
	)
	if !w.focused {
		bg = w.InactiveTitleBackground
		fg = w.InactiveTitleForeground
	}
	w.titleBar.SetBackground(bg)
	w.titleLabel.Font.Color = fg
	w.titleLabel.Font.Stroke = bg.Darken(40)
}

// Children returns the window's child widgets.
func (w *Window) Children() []Widget {
	return []Widget{
		w.body,
	}
}

// Pack a child widget into the window's main frame.
func (w *Window) Pack(child Widget, config ...Pack) {
	w.content.Pack(child, config...)
}

// Place a child widget into the window's main frame.
func (w *Window) Place(child Widget, config Place) {
	w.content.Place(child, config)
}

// TitleBar returns the title bar widget.
func (w *Window) TitleBar() *Frame {
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

	// Call the BaseWidget Compute in case we have subscribers.
	w.BaseWidget.Compute(e)
}

// Present the window.
func (w *Window) Present(e render.Engine, P render.Point) {
	w.body.Present(e, P)

	// Call the BaseWidget Present in case we have subscribers.
	w.BaseWidget.Present(e, P)
}
