package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
)

// Menu is a rectangle that holds menu items.
type Menu struct {
	BaseWidget
	Name string

	body *Frame
}

// NewMenu creates a new Menu. It is hidden by default. Usually you'll
// use it with a MenuButton or in a right-click handler.
func NewMenu(name string) *Menu {
	w := &Menu{
		Name: name,
		body: NewFrame(name + ":Body"),
	}
	w.body.Configure(Config{
		Width:       150,
		BorderSize:  12,
		BorderStyle: BorderRaised,
		Background:  render.Grey,
	})
	w.IDFunc(func() string {
		return fmt.Sprintf("Menu<%s>", w.Name)
	})
	return w
}

// Compute the menu
func (w *Menu) Compute(e render.Engine) {
	w.body.Compute(e)

	// Call the BaseWidget Compute in case we have subscribers.
	w.BaseWidget.Compute(e)
}

// Present the menu
func (w *Menu) Present(e render.Engine, p render.Point) {
	w.body.Present(e, p)

	// Call the BaseWidget Present in case we have subscribers.
	w.BaseWidget.Present(e, p)
}

// AddItem quickly adds an item to a menu.
func (w *Menu) AddItem(label string, command func()) *MenuItem {
	menu := NewMenuItem(label, command)
	w.Pack(menu)
	return menu
}

// Pack a menu item onto the menu.
func (w *Menu) Pack(item *MenuItem) {
	w.body.Pack(item, Pack{
		Side: NE,
		// Expand: true,
		// Padding: 8,
		FillX: true,
	})
}

// MenuItem is an item in a Menu.
type MenuItem struct {
	Button
	Label       string
	Accelerator string
	Command     func()
	button      *Button
}

// NewMenuItem creates a new menu item.
func NewMenuItem(label string, command func()) *MenuItem {
	w := &MenuItem{
		Label:   label,
		Command: command,
	}
	w.IDFunc(func() string {
		return fmt.Sprintf("MenuItem<%s>", w.Label)
	})

	font := DefaultFont
	font.Color = render.White
	font.PadX = 12
	w.Button.child = NewLabel(Label{
		Text: label,
		Font: font,
	})
	w.Button.Configure(Config{
		Background: render.Blue,
	})

	w.Button.Handle(Click, func(ed EventData) {
		w.Command()
	})

	// Assign the button
	return w
}
