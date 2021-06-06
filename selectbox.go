package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/theme"
)

// SelectBox is a kind of MenuButton which allows choosing a value from a list.
type SelectBox struct {
	MenuButton

	name       string

	// Configurables after SelectBox creation.
	AlwaysChange bool // always call the Change event, even if selection not changed.
	
	// Child widgets specific to the SelectBox.
	frame *Frame
	label *Label
	arrow *Image

	// Data storage.
	textVariable string
	values []SelectValue
}

// SelectValue holds a mapping between a text label for a SelectBox and
// its underlying value (an arbitrary data type).
type SelectValue struct {
	Label string
	Value interface{}
}

// NewSelectBox creates a new SelectBox.
//
// The Label configuration passed in should be used to set font styles
// and padding; the Text, TextVariable and IntVariable of the Label will
// all be ignored, as SelectBox will handle the values internally.
func NewSelectBox(name string, withLabel Label) *SelectBox {
	w := &SelectBox{
		name: name,
		textVariable: "Choose one",
		values: []SelectValue{},
	}

	// Ensure the label has no text of its own.
	withLabel.Text = ""
	withLabel.TextVariable = &w.textVariable
	withLabel.IntVariable = nil

	w.frame = NewFrame(name + " Frame")
	w.Button.child = w.frame

	w.label = NewLabel(withLabel)
	w.frame.Pack(w.label, Pack{
		Side: W,
	})

	// arrow, _ := GetGlyph(GlyphDownArrow9x9)
	// w.image = ImageFromImage(arrow, )

	// Configure the button's appearance.
	w.Button.Configure(Config{
		BorderSize: 2,
		BorderStyle: BorderSunken,
		Background: render.White,
	})

	// Set sensible default padding on the label.
	if w.label.Font.Padding == 0 && w.label.Font.PadX == 0 && w.label.Font.PadY == 0 {
		w.label.Font.PadX = 4
		w.label.Font.PadY = 2
	}

	w.IDFunc(func() string {
		return fmt.Sprintf("SelectBox<%s>", name)
	})

	w.setup()
	return w
}

// AddItem adds a new option to the SelectBox's menu.
// The label is the text value to display.
// The value is the underlying value (string or int) for the TextVariable or IntVariable.
// The function callback runs when the option is picked.
func (w *SelectBox) AddItem(label string, value interface{}, f func()) {
	// Add this label and its value mapping to the SelectBox.
	w.values = append(w.values, SelectValue{
		Label: label,
		Value: value,
	})

	// Call the inherited MenuButton.AddItem.
	w.MenuButton.AddItem(label, func() {
		// Set the bound label.
		var changed = w.textVariable != label
		w.textVariable = label

		if changed || w.AlwaysChange {
			w.Event(Change, EventData{
				Supervisor: w.MenuButton.supervisor,
			})
		}
	})

	// If the current text label isn't in the options, pick
	// the first option.
	if _, ok := w.GetValue(); !ok {
		w.textVariable = w.values[0].Label
	}
}

// TODO: RemoveItem()

// Value returns the currently selected item in the SelectBox.
//
// Returns the SelectValue and true on success, and the Label or underlying Value
// can be read from the SelectValue struct. If no valid option is selected, the
// bool value returns false.
func (w *SelectBox) GetValue() (SelectValue, bool) {
	for _, row := range w.values {
		if w.textVariable == row.Label {
			return row, true
		}
	}

	return SelectValue{}, false
}

// Compute to re-evaluate the button state (in the case of radio buttons where
// a different button will affect the state of this one when clicked).
func (w *SelectBox) Compute(e render.Engine) {
	w.MenuButton.Compute(e)
}

// setup the UI components and event handlers.
func (w *SelectBox) setup() {
	w.Configure(Config{
		BorderSize:  1,
		BorderStyle: BorderSunken,
		Background:  theme.InputBackgroundColor,
	})

	w.Handle(MouseOver, func(ed EventData) error {
		w.hovering = true
		w.SetBackground(theme.ButtonHoverColor)
		return nil
	})
	w.Handle(MouseOut, func(ed EventData) error {
		w.hovering = false
		w.SetBackground(theme.InputBackgroundColor)
		return nil
	})

	w.Handle(MouseDown, func(ed EventData) error {
		w.clicked = true
		w.SetBackground(theme.ButtonBackgroundColor)
		return nil
	})
	w.Handle(MouseUp, func(ed EventData) error {
		w.clicked = false
		w.SetBackground(theme.InputBackgroundColor)
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

