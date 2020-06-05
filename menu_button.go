package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/theme"
)

// MenuButton is a button that opens a menu when clicked.
//
// After creating a MenuButton, call AddItem() to add options and callback
// functions to fill out the menu. When the MenuButton is clicked, its menu
// will be drawn and take modal priority in the Supervisor.
type MenuButton struct {
	Button

	name       string
	supervisor *Supervisor
	menu       *Menu
}

// NewMenuButton creates a new MenuButton (labels recommended).
//
// If the child is a Label, this function will set some sensible padding on
// its font if the Label does not already have non-zero padding set.
func NewMenuButton(name string, child Widget) *MenuButton {
	w := &MenuButton{
		name: name,
	}
	w.Button.child = child

	// If it's a Label (most common), set sensible default padding.
	if label, ok := child.(*Label); ok {
		if label.Font.Padding == 0 && label.Font.PadX == 0 && label.Font.PadY == 0 {
			label.Font.PadX = 8
			label.Font.PadY = 4
		}
	}

	w.IDFunc(func() string {
		return fmt.Sprintf("MenuButton<%s>", name)
	})

	w.setup()
	return w
}

// Supervise the MenuButton. This is necessary for the pop-up menu to work
// when the button is clicked.
func (w *MenuButton) Supervise(s *Supervisor) {
	w.initMenu()
	w.supervisor = s
	w.menu.Supervise(s)
}

// AddItem adds a new option to the MenuButton's menu.
func (w *MenuButton) AddItem(label string, f func()) {
	w.initMenu()
	w.menu.AddItem(label, f)
}

// AddItemAccel adds a new menu option with hotkey text.
func (w *MenuButton) AddItemAccel(label string, accelerator string, f func()) *MenuItem {
	w.initMenu()
	return w.menu.AddItemAccel(label, accelerator, f)
}

// AddSeparator adds a separator to the menu.
func (w *MenuButton) AddSeparator() {
	w.initMenu()
	w.menu.AddSeparator()
}

// Compute to re-evaluate the button state (in the case of radio buttons where
// a different button will affect the state of this one when clicked).
func (w *MenuButton) Compute(e render.Engine) {
	if w.menu != nil {
		w.menu.Compute(e)
		w.positionMenu(e)
	}
}

// positionMenu sets the position where the pop-up menu will appear when
// the button is clicked. Usually, the menu appears below and to the right of
// the button. But if the menu will hit a window boundary, its position will
// be adjusted to fit the window while trying not to overlap its own button.
func (w *MenuButton) positionMenu(e render.Engine) {
	var (
		// Position and size of the MenuButton button.
		buttonPoint = w.Point()
		buttonSize  = w.Size()

		// Size of the actual desktop window.
		Width, Height = e.WindowSize()
	)

	// Ideal location: below and to the right of the button.
	w.menu.MoveTo(render.Point{
		X: buttonPoint.X,
		Y: buttonPoint.Y + buttonSize.H + w.BoxThickness(2),
	})

	var (
		// Size of the menu.
		menuPoint = w.menu.Point()
		menuSize  = w.menu.Rect()
		margin    = 8  // keep away from directly touching window edges
		topMargin = 32 // keep room for standard Menu Bar
	)

	// Will we clip out the bottom of the window?
	if menuPoint.Y+menuSize.H+margin > Height {
		// Put us above the button instead, with the bottom of the
		// menu touching the top of the button.
		menuPoint = render.Point{
			X: buttonPoint.X,
			Y: buttonPoint.Y - menuSize.H - w.BoxThickness(2),
		}

		// If this would put us over the TOP edge of the window now,
		// cap the movement so the top of the menu is visible. We can't
		// avoid overlapping the button with the menu so might as well
		// start now.
		if menuPoint.Y < topMargin {
			menuPoint.Y = topMargin
		}

		w.menu.MoveTo(menuPoint)
	}

	// Will we clip out the right of the window?
	if menuPoint.X+menuSize.W > Width {
		// Move us in from the right side of the window.
		var delta = Width - menuSize.W - margin
		w.menu.MoveTo(render.Point{
			X: delta,
			Y: menuPoint.Y,
		})
	}
	_ = Width
}

// setup the common things between checkboxes and radioboxes.
func (w *MenuButton) setup() {
	w.Configure(Config{
		BorderSize:  1,
		BorderStyle: BorderSolid,
		Background:  theme.ButtonBackgroundColor,
	})

	w.Handle(MouseOver, func(ed EventData) error {
		w.hovering = true
		w.SetBorderStyle(BorderRaised)
		return nil
	})
	w.Handle(MouseOut, func(ed EventData) error {
		w.hovering = false
		w.SetBorderStyle(BorderSolid)
		return nil
	})

	w.Handle(MouseDown, func(ed EventData) error {
		w.clicked = true
		w.SetBorderStyle(BorderSunken)
		return nil
	})
	w.Handle(MouseUp, func(ed EventData) error {
		w.clicked = false
		return nil
	})

	w.Handle(Click, func(ed EventData) error {
		// Are we properly configured?
		if w.supervisor != nil && w.menu != nil {
			w.menu.Show()
			w.supervisor.PushModal(w.menu)
		}
		return nil
	})
}

// initialize the Menu widget.
func (w *MenuButton) initMenu() {
	if w.menu == nil {
		w.menu = NewMenu(w.name + ":Menu")
		w.menu.Hide()

		// Handle closing the menu when clicked outside.
		w.menu.Handle(CloseModal, func(ed EventData) error {
			ed.Supervisor.PopModal(w.menu)
			return nil
		})
	}
}
