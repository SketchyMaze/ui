package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/style"
	"git.kirsle.net/go/ui/theme"
)

// TabFrame is a tabbed notebook of multiple frames showing
// tab names along the top and clicking them reveals each
// named tab.
type TabFrame struct {
	Name string
	Frame

	supervisor *Supervisor
	style      *style.Button

	// Child widgets.
	header        *Frame
	content       *Frame
	tabButtons    []*Button
	tabFrames     []*Frame
	currentTabKey string
}

// NewTabFrame creates a new Frame.
func NewTabFrame(name string) *TabFrame {
	w := &TabFrame{
		Name:       name,
		style:      Theme.TabFrame,
		header:     NewFrame(name + " Header"),
		content:    NewFrame(name + " Content"),
		tabButtons: []*Button{},
		tabFrames:  []*Frame{},
	}

	// Initialize the root frame of this widget.
	w.Frame.Setup()

	// Pack the high-level layout into the root frame.
	// Only root needs to Present for this widget.
	w.Frame.Pack(w.header, Pack{
		Side:  N,
		FillX: true,
	})
	w.Frame.Pack(w.content, Pack{
		Side:   N,
		Fill:   true,
		Expand: true,
	})

	w.SetBackground(render.RGBA(1, 0, 0, 0)) // invisible default BG
	w.IDFunc(func() string {
		return fmt.Sprintf("TabFrame<%s>",
			name,
		)
	})

	return w
}

// AddTab creates a new content tab. The key is a unique identifier
// for the tab and is how the TabFrame knows which tab is selected.
//
// The child widget would probably be a Label or Image but could be
// any other kind of widget.
//
// The first tab added becomes the selected tab by default.
func (w *TabFrame) AddTab(key string, child Widget) *Frame {
	// Create the tab button for this tab.
	button := NewButton(key, child)
	button.SetStyle(w.style)
	button.FixedColor = true
	button.SetOutlineSize(0)
	button.SetBorderSize(1)
	w.setButtonStyle(button, len(w.tabButtons) == 0)
	w.header.Pack(button, Pack{
		Side: W,
	})

	button.Handle(MouseDown, func(ed EventData) error {
		w.SetTab(key)
		return nil
	})

	// Create the frame of the tab's body.
	frame := NewFrame(key)
	frame.Configure(Config{
		BorderSize:  w.style.BorderSize,
		Background:  w.style.Background,
		BorderStyle: BorderRaised,
	})
	if len(w.tabFrames) > 0 {
		frame.Hide()
	}

	// Pack this frame into the content part of the widget.
	w.content.Pack(frame, Pack{
		Side:  N,
		FillX: true,
	})

	w.tabButtons = append(w.tabButtons, button)
	w.tabFrames = append(w.tabFrames, frame)
	return frame
}

// SetTabsHidden can hide the tab buttons and reveal only their frames.
// It would be up to the caller to SetTab between the frames, using the
// TabFrame only for placement and tab handling.
func (w *TabFrame) SetTabsHidden(hidden bool) {
	if hidden {
		w.header.Hide()
	} else {
		w.header.Show()
	}
}

// Header returns access to the ui.Frame that holds the tab buttons. Use
// at your own risk -- the UI arrangement in this Frame is not guaranteed
// stable.
func (w *TabFrame) Header() *Frame {
	return w.header
}

// set the tab style between active and inactive
func (w *TabFrame) setButtonStyle(button *Button, active bool) {
	var style = button.GetStyle()
	if active {
		button.SetBackground(style.Background)
		button.SetBorderStyle(BorderRaised)
		if label, ok := button.child.(*Label); ok {
			label.Font.Color = style.Foreground
		}
	} else {
		button.SetBackground(style.Background.Darken(theme.BorderColorOffset))
		button.SetBorderStyle(BorderSolid)
		if label, ok := button.child.(*Label); ok {
			label.Font.Color = style.Foreground
		}
	}
}

// SetTab changes the selected tab to the new value. If the
// tab doesn't exist, the first tab is selected.
func (w *TabFrame) SetTab(key string) bool {
	var found bool
	for i, frame := range w.tabFrames {
		button := w.tabButtons[i]
		if frame.Name == key {
			frame.Show()
			w.setButtonStyle(button, true)
			w.currentTabKey = key
			found = true
		} else {
			frame.Hide()
			w.setButtonStyle(button, false)
		}
	}

	if !found && len(w.tabFrames) > 0 {
		w.tabFrames[0].Show()
		w.currentTabKey = w.tabFrames[0].Name
	}

	return found
}

// Supervise activates the tab frame using your supervisor. If you
// don't call this, the tab buttons won't be clickable!
//
// Call this AFTER adding all tabs. This function calls Supervisor.Add
// on all tab buttons.
func (w *TabFrame) Supervise(supervisor *Supervisor) {
	for _, button := range w.tabButtons {
		supervisor.Add(button)
	}
}

// SetStyle controls the visual styling of the tab button bar.
func (w *TabFrame) SetStyle(style *style.Button) {
	w.style = style
	for _, button := range w.tabButtons {
		button.SetStyle(style)
		w.setButtonStyle(button, !button.Hidden())
	}
}

// Compute the size of the Frame.
func (w *TabFrame) Compute(e render.Engine) {
	// Compute all the child frames.
	w.Frame.Compute(e)

	// Call the BaseWidget Compute in case we have subscribers.
	w.BaseWidget.Compute(e)
}

// Present the Frame.
func (w *TabFrame) Present(e render.Engine, P render.Point) {
	if w.Hidden() {
		return
	}

	var (
		S = w.Size()
	)

	// Draw the widget's border and everything.
	w.DrawBox(e, P)

	// Draw the background color.
	e.DrawBox(w.Background(), render.Rect{
		X: P.X + w.BoxThickness(1),
		Y: P.Y + w.BoxThickness(1),
		W: S.W - w.BoxThickness(2),
		H: S.H - w.BoxThickness(2),
	})

	// Present the root frame.
	w.Frame.Present(e, P)

	// Draw the borders over the tabs.
	w.presentBorders(e, P)

	// Call the BaseWidget Present in case we have subscribers.
	w.BaseWidget.Present(e, P)
}

/*
presentBorders handles drawing the borders around tab buttons.

The tabs are simple Button widgets but drawn with no borders. Instead,
borders are painted on post-hoc in the Present function.
*/
func (w *TabFrame) presentBorders(e render.Engine, P render.Point) {
	if len(w.tabButtons) == 0 || w.header.Hidden() {
		return
	}

	// Prep some variables.
	var (
		// The 1st and last tab button widgets.
		first       = w.tabButtons[0]
		last        = w.tabButtons[len(w.tabButtons)-1]
		topLeft     = AbsolutePosition(first)
		bottomRight = AbsolutePosition(last)

		// The absolute bounding box of the tabs part of the UI,
		// from the top-left corner of Tab #1 to the bottom-right
		// corner of the final tab.
		bounding = render.Rect{
			X: P.X, //topLeft.X + first.BoxThickness(4),
			Y: P.Y, //topLeft.Y + first.BoxThickness(4),
			W: bottomRight.X + last.Size().W - topLeft.X,
			H: bottomRight.Y + last.Size().H - topLeft.Y,
		}

		// The very bottom edge of the whole tab bar,
		// to overlap the BorderSize=1 along their buttons.
		bottomLine = []render.Point{
			render.NewPoint(P.X+1, bounding.Y+bounding.H-1),
			render.NewPoint(bounding.X+bounding.W-1, bounding.Y+bounding.H-1),
		}
	)

	// Draw a shadow border on all the inactive tabs' right edges,
	// so they don't all blend together in solid grey.
	// Note: the active button has a BorderSize=1 and others are 0.
	for i, button := range w.tabButtons {
		if button.Name != w.currentTabKey {
			// If it immediately precedes the current tab, do not draw the line,
			// it would cover the highlight color of the current tab's button.
			if i+1 < len(w.tabButtons) && w.tabButtons[i+1].Name == w.currentTabKey {
				continue
			}

			var (
				abs    = AbsolutePosition(button)
				size   = button.BoxSize()
				points = []render.Point{
					render.NewPoint(abs.X+size.W-1, abs.Y+2),
					render.NewPoint(abs.X+size.W-1, abs.Y+size.H-2),
				}
			)
			e.DrawLine(button.Background().Darken(theme.BorderColorOffset), points[0], points[1])
		}
	}

	// Erase the button edge from all tabs.
	e.DrawLine(w.style.Background, bottomLine[0], bottomLine[1])
	e.DrawBox(w.style.Background, render.Rect{
		X: bottomLine[0].X + 1,
		Y: bottomLine[0].Y,
		W: bounding.W - 2,
		H: 4,
	})
}
