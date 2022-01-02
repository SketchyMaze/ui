package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/style"
)

// Window is a frame with a title bar.
type Window struct {
	BaseWidget
	Title string

	// Title bar colors. Sensible defaults are chosen in NewWindow but you
	// may customize after the fact.
	ActiveTitleBackground   render.Color
	ActiveTitleForeground   render.Color
	InactiveTitleBackground render.Color
	InactiveTitleForeground render.Color

	// Private widgets.
	style        *style.Window
	body         *Frame
	titleBar     *Frame
	titleLabel   *Label
	titleButtons []*Button
	content      *Frame

	// Configured title bar buttons.
	buttonsEnabled int

	// Window manager controls.
	dragging      bool
	startDragAt   render.Point // cursor position when drag began
	dragOrigPoint render.Point // original position of window at drag start
	focused       bool
	managed       bool          // window is managed by Supervisor
	maximized     bool          // toggled by MaximizeButton
	origPoint     render.Point  // placement before a maximize
	origSize      render.Rect   // size before a maximize
	engine        render.Engine // hang onto the render engine, for Maximize support.
}

// NewWindow creates a new window.
func NewWindow(title string) *Window {
	w := &Window{
		Title: title,
		body:  NewFrame("body:" + title),

		// Default title bar colors.
		ActiveTitleBackground:   render.Blue,
		ActiveTitleForeground:   render.White,
		InactiveTitleBackground: render.DarkGrey,
		InactiveTitleForeground: render.Grey,
	}
	w.IDFunc(func() string {
		return fmt.Sprintf("Window<%s %+v>",
			w.Title, w.focused,
		)
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

	w.SetStyle(Theme.Window)

	return w
}

// SetStyle sets the window style.
func (w *Window) SetStyle(v *style.Window) {
	if v == nil {
		v = &style.DefaultWindow
	}

	w.style = v
	w.body.Configure(Config{
		Background:  w.style.ActiveBackground,
		BorderSize:  2,
		BorderStyle: BorderRaised,
	})
	if w.focused {
		w.titleBar.SetBackground(w.style.ActiveTitleBackground)
		w.titleLabel.Font.Color = w.style.ActiveTitleForeground
	} else {
		w.titleBar.SetBackground(w.style.InactiveTitleBackground)
		w.titleLabel.Font.Color = w.style.InactiveTitleForeground
	}
}

// setupTitlebar creates the title bar frame of the window.
func (w *Window) setupTitleBar() (*Frame, *Label) {
	frame := NewFrame("Titlebar for Window: " + w.Title)
	frame.Configure(Config{
		Background: w.ActiveTitleBackground,
	})

	// Title label.
	label := NewLabel(Label{
		TextVariable: &w.Title,
		Font: render.Text{
			Color:   w.ActiveTitleForeground,
			Size:    11,
			Stroke:  w.ActiveTitleBackground.Darken(40),
			Padding: 2,
		},
	})
	frame.Pack(label, Pack{
		Side: W,
	})

	// Window buttons.
	var buttons = []struct {
		If    bool
		Label string
		Event Event
	}{
		{
			Label: "×",
			Event: CloseWindow,
		},
		{
			Label: "+",
			Event: MaximizeWindow,
		},
		{
			Label: "_",
			Event: MinimizeWindow,
		},
	}
	w.titleButtons = make([]*Button, len(buttons))
	for i, cfg := range buttons {
		cfg := cfg
		btn := NewButton(
			fmt.Sprintf("Title Button %d for Window: %s", i, w.Title),
			NewLabel(Label{
				Text: cfg.Label,
				Font: render.Text{
					// Color: w.ActiveTitleForeground,
					Size:    8,
					Padding: 2,
				},
			}),
		)
		btn.SetBorderSize(0)
		btn.Handle(Click, func(ed EventData) error {
			w.Event(cfg.Event, ed)
			return ErrStopPropagation // TODO: doesn't work :(
		})
		btn.Hide()
		w.titleButtons[i] = btn

		frame.Pack(btn, Pack{
			Side: E,
		})
	}

	return frame, label
}

// SetButtons sets the title bar buttons to show in the window.
//
// The value should be the OR of CloseButton, MaximizeButton and MinimizeButton
// that you want to be enabled.
//
// Window buttons only work if the window is managed by Supervisor and you have
// called the Supervise() method of the window.
func (w *Window) SetButtons(buttons int) {
	// Show/hide each button based on the value given.
	var toggle = []struct {
		Value int
		Index int
	}{
		{
			Value: CloseButton,
			Index: 0,
		},
		{
			Value: MaximizeButton,
			Index: 1,
		},
		{
			Value: MinimizeButton,
			Index: 2,
		},
	}

	for _, item := range toggle {
		if buttons&item.Value == item.Value {
			w.titleButtons[item.Index].Show()
		} else {
			w.titleButtons[item.Index].Hide()
		}
	}
}

// Supervise enables the window to be dragged around by its title bar by
// adding its relevant event hooks to your Supervisor.
func (w *Window) Supervise(s *Supervisor) {
	// Add a click handler to the title bar to enable dragging.
	w.titleBar.Handle(MouseDown, func(ed EventData) error {
		w.startDragAt = ed.Point
		w.dragOrigPoint = w.Point()

		s.DragStartWidget(w)
		return nil
	})

	// Clicking anywhere in the window focuses the window.
	w.Handle(MouseDown, func(ed EventData) error {
		s.FocusWindow(w)
		return nil
	})

	// Window as a whole receives DragMove events while being dragged.
	w.Handle(DragMove, func(ed EventData) error {
		// Get the delta of movement from where we began.
		delta := w.startDragAt.Compare(ed.Point)
		if delta != render.Origin {
			moveTo := w.dragOrigPoint
			moveTo.Add(delta)
			w.MoveTo(moveTo)
		}
		return nil
	})

	// Window button handlers.
	w.Handle(CloseWindow, func(ed EventData) error {
		w.Hide()
		return nil
	})
	w.Handle(MaximizeWindow, func(ed EventData) error {
		w.SetMaximized(!w.maximized)
		return nil
	})

	// Add the title bar to the supervisor.
	s.Add(w.titleBar)
	for _, btn := range w.titleButtons {
		s.Add(btn)
	}
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
		bg = w.style.ActiveTitleBackground
		fg = w.style.ActiveTitleForeground
	)
	if !w.focused {
		bg = w.style.InactiveTitleBackground
		fg = w.style.InactiveTitleForeground
	}
	w.titleBar.SetBackground(bg)
	w.titleLabel.Font.Color = fg
	w.titleLabel.Font.Stroke = bg.Darken(40)
}

// Maximized returns whether the window is maximized.
func (w *Window) Maximized() bool {
	return w.maximized
}

// SetMaximized sets the state of the maximized window.
// Must have called Compute() once before so the window can hang on to the
// render.Engine, to calculate the size of the parent window.
func (w *Window) SetMaximized(v bool) {
	w.maximized = v

	if v && w.engine != nil {
		w.origPoint = w.Point()
		w.origSize = w.Size()
		w.MoveTo(render.Origin)
		w.Resize(render.NewRect(w.engine.WindowSize()))
		w.Compute(w.engine)
	} else if w.engine != nil {
		w.MoveTo(w.origPoint)
		w.Resize(w.origSize)
		w.Compute(w.engine)
	}
}

// Size returns the window's size (the size of its underlying body frame,
// including its title bar and content frames).
func (w *Window) Size() render.Rect {
	return w.body.Size()
}

// Resize the window.
func (w *Window) Resize(size render.Rect) {
	w.BaseWidget.Resize(size)
	w.body.Resize(size)
}

// Center the window on screen by providing your screen (app window) size.
func (w *Window) Center(width, height int) {
	w.MoveTo(render.Point{
		X: (width / 2) - (w.Size().W / 2),
		Y: (height / 2) - (w.Size().H / 2),
	})
}

// Close the window, hiding it from display and calling its CloseWindow handler.
func (w *Window) Close() {
	w.Hide()
	w.Event(CloseWindow, EventData{})
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
	w.body.Place(child, config)
}

// TitleBar returns the title bar widgets.
func (w *Window) TitleBar() (*Frame, *Label) {
	return w.titleBar, w.titleLabel
}

// Configure the widget. Color and style changes are passed down to the inner
// content frame of the window.
func (w *Window) Configure(C Config) {
	w.BaseWidget.Configure(C)
	w.body.Configure(C)

	// Don't pass dimensions down any further than the body.
	// TODO: this causes the content frame to compute its size
	// dynamically based on Packed widgets, but if using Place on
	// your window, the content frame doesn't know a size by which
	// to place the child relative to (Frame has size 0x0).
	// Commenting out these two lines causes windows to render very
	// incorrectly (child frame content flying off the window bottom).
	// In the meantime, Window.Place intercepts it and draws it onto
	// the parent window directly so it works how you expect.
	C.Width = 0
	C.Height = 0
	w.content.Configure(C)
}

// ConfigureTitle configures the title bar widget.
func (w *Window) ConfigureTitle(C Config) {
	w.titleBar.Configure(C)
}

// ContentFrame returns the main content Frame of this window.
func (w *Window) ContentFrame() *Frame {
	return w.content
}

// Compute the window.
func (w *Window) Compute(e render.Engine) {
	w.engine = e // hang onto it in case of maximize
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

// Destroy hides the window.
func (w *Window) Destroy() {
	w.Hide()
}
