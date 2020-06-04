package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/theme"
)

// MenuFont is the default font settings for MenuBar buttons.
var MenuFont = render.Text{
	Size:  12,
	Color: render.Black,
	PadX:  4,
	PadY:  2,
}

// MenuBar is a frame that holds several MenuButtons, such as for the main
// menu at the top of a window.
type MenuBar struct {
	Frame
	name string

	supervisor *Supervisor
	buttons    []*MenuButton
}

// NewMenuBar creates a new menu bar frame.
func NewMenuBar(name string) *MenuBar {
	w := &MenuBar{
		name:    name,
		buttons: []*MenuButton{},
	}
	w.SetBackground(theme.ButtonBackgroundColor)
	w.Frame.Setup()
	w.IDFunc(func() string {
		return fmt.Sprintf("MenuBar<%s>", w.name)
	})
	return w
}

// Supervise the menu bar, making its child menu buttons work correctly.
func (w *MenuBar) Supervise(s *Supervisor) {
	w.supervisor = s

	// Supervise the existing buttons.
	for _, btn := range w.buttons {
		s.Add(btn)
		btn.Supervise(s)
	}
}

// AddMenu adds a new menu button to the bar. Returns the MenuButton
// object so that you can add items to it.
func (w *MenuBar) AddMenu(label string) *MenuButton {
	btn := NewMenuButton(label, NewLabel(Label{
		Text: label,
		Font: MenuFont,
	}))
	w.buttons = append(w.buttons, btn)

	// Pack and supervise it.
	w.Pack(btn, Pack{
		Side: W,
	})
	if w.supervisor != nil {
		w.supervisor.Add(btn)
		btn.Supervise(w.supervisor)
	}
	return btn
}

// PackTop returns the default Frame Pack settings to place the menu
// at the top of the parent widget.
func (w *MenuBar) PackTop() Pack {
	return Pack{
		Side:  N,
		FillX: true,
	}
}
