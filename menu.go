package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/theme"
)

// MenuWidth sets the width of all popup menus. TODO, widths should be automatic.
var MenuWidth = 180

// Menu is a frame that holds menu items. It is the
type Menu struct {
	BaseWidget
	Name string

	supervisor *Supervisor
	body       *Frame
	items      []*MenuItem
}

// NewMenu creates a new Menu. It is hidden by default. Usually you'll
// use it with a MenuButton or in a right-click handler.
func NewMenu(name string) *Menu {
	w := &Menu{
		Name:  name,
		body:  NewFrame(name + ":Body"),
		items: []*MenuItem{},
	}
	w.body.Configure(Config{
		Width:       MenuWidth,
		Height:      100,
		BorderSize:  0,
		BorderStyle: BorderRaised,
		Background:  theme.ButtonBackgroundColor,
	})
	w.body.SetParent(w)
	w.IDFunc(func() string {
		return fmt.Sprintf("Menu<%s>", w.Name)
	})
	return w
}

// Children returns the child frame of the menu.
func (w *Menu) Children() []Widget {
	return []Widget{
		w.body,
	}
}

// Supervise the Menu. This will add all current and future MenuItem widgets
// to the supervisor.
func (w *Menu) Supervise(s *Supervisor) {
	w.supervisor = s
	for _, item := range w.items {
		w.supervisor.Add(item)
	}
}

// Compute the menu
func (w *Menu) Compute(e render.Engine) {
	w.body.Compute(e)

	// TODO: ideally the Frame Pack Compute would fix the size of the body
	// for the height to match the height of the menu items... but for now
	// manually set the height.
	var maxWidth int
	var height int
	for _, child := range w.body.Children() {
		size := child.Size()
		if size.W > maxWidth {
			maxWidth = size.W
		}
		height += child.Size().H
	}
	w.body.Resize(render.NewRect(maxWidth, height))

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
	menu := NewMenuItem(label, "", command)

	// Add a Click handler that closes the menu when a selection is made.
	menu.Handle(Click, w.menuClickHandler)

	w.Pack(menu)
	return menu
}

// AddItemAccel quickly adds an item to a menu with a shortcut key label.
func (w *Menu) AddItemAccel(label string, accelerator string, command func()) *MenuItem {
	menu := NewMenuItem(label, accelerator, command)

	// Add a Click handler that closes the menu when a selection is made.
	menu.Handle(Click, w.menuClickHandler)

	w.Pack(menu)
	return menu
}

// Click handler for all menu items, to also close the menu behind them.
func (w *Menu) menuClickHandler(ed EventData) error {
	if w.supervisor != nil {
		w.supervisor.PopModal(w)
	}
	return nil
}

// AddSeparator adds a separator bar to the menu to delineate items.
func (w *Menu) AddSeparator() *MenuItem {
	sep := NewMenuSeparator()
	w.Pack(sep)
	return sep
}

// Pack a menu item onto the menu.
func (w *Menu) Pack(item *MenuItem) {
	w.items = append(w.items, item)
	w.body.Pack(item, Pack{
		Side:  N,
		FillX: true,
	})
	if w.supervisor != nil {
		w.supervisor.Add(item)
	}
}

// Size returns the size of the menu's body.
func (w *Menu) Size() render.Rect {
	return w.body.Size()
}

// Rect returns the rect of the menu's body.
func (w *Menu) Rect() render.Rect {
	// TODO: the height reports wrong (0), manually add up the MenuItem sizes.
	// This manifests in Supervisor.runWidgetEvents when checking if the cursor
	// clicked outside the rect of the active menu modal.
	rect := w.body.Rect()
	rect.H = 0
	for _, child := range w.body.Children() {
		rect.H += child.Size().H
	}
	return rect
}

// MenuItem is an item in a Menu.
type MenuItem struct {
	Button
	Label       string
	Accelerator string
	Command     func()
	separator   bool
	button      *Button

	// store of most recent bg color set on a menu item
	cacheBg render.Color
	cacheFg render.Color
}

// NewMenuItem creates a new menu item.
func NewMenuItem(label, accelerator string, command func()) *MenuItem {
	w := &MenuItem{
		Label:       label,
		Accelerator: accelerator,
		Command:     command,
	}
	w.IDFunc(func() string {
		return fmt.Sprintf("MenuItem<%s>", w.Label)
	})

	font := DefaultFont
	font.Color = render.Black
	font.PadX = 12
	font.PadY = 2

	// The button child will be a Frame so we can have a left-aligned label
	// and a right-aligned accelerator.
	frame := NewFrame(label + ":Frame")
	frame.Configure(Config{
		Width: MenuWidth,
	})
	{
		// Left of frame: menu item label
		lbl := NewLabel(Label{
			Text: label,
			Font: font,
		})
		frame.Pack(lbl, Pack{
			Side: W,
		})

		// On the right: accelerator shortcut key
		if accelerator != "" {
			accel := NewLabel(Label{
				Text: accelerator,
				Font: font,
			})
			frame.Pack(accel, Pack{
				Side: E,
			})
		}
	}

	w.Button.child = frame
	w.Button.Configure(Config{
		BorderSize: 0,
		Background: theme.ButtonBackgroundColor,
	})

	w.Button.Handle(MouseOver, func(ed EventData) error {
		w.setHoverStyle(true)
		return nil
	})
	w.Button.Handle(MouseOut, func(ed EventData) error {
		w.setHoverStyle(false)
		return nil
	})

	w.Button.Handle(Click, func(ed EventData) error {
		w.Command()
		return nil
	})

	// Assign the button
	return w
}

// NewMenuSeparator creates a separator menu item.
func NewMenuSeparator() *MenuItem {
	w := &MenuItem{
		separator: true,
	}
	w.IDFunc(func() string {
		return "MenuItem<separator>"
	})
	w.Button.child = NewFrame("Menu Separator")
	w.Button.Configure(Config{
		Width:       MenuWidth,
		Height:      2,
		BorderSize:  1,
		BorderStyle: BorderSunken,
		BorderColor: render.Grey,
	})
	return w
}

// Set the hover styling (text/bg color)
func (w *MenuItem) setHoverStyle(hovering bool) {
	// Note: this only works if the MenuItem is using the standard
	// Frame and Labels layout created by AddItem(). If not, this function
	// does nothing.

	// BG color.
	if hovering {
		w.cacheBg = w.Background()
		w.SetBackground(render.SkyBlue)
	} else {
		w.SetBackground(w.cacheBg)
	}

	frame, ok := w.Button.child.(*Frame)
	if !ok {
		return
	}

	for _, widget := range frame.Children() {
		if label, ok := widget.(*Label); ok {
			if hovering {
				w.cacheFg = label.Font.Color
				label.Font.Color = render.White
			} else {
				label.Font.Color = w.cacheFg
			}
		}
	}
}
