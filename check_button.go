package ui

import (
	"fmt"
	"strconv"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/theme"
)

// CheckButton implements a checkbox and radiobox widget. It's based on a
// Button and holds a boolean or string pointer (boolean for checkbox,
// string for radio).
type CheckButton struct {
	Button
	BoolVar   *bool
	StringVar *string
	Value     string
}

// NewCheckButton creates a new CheckButton.
func NewCheckButton(name string, boolVar *bool, child Widget) *CheckButton {
	w := &CheckButton{
		BoolVar: boolVar,
	}
	w.Button.child = child
	w.IDFunc(func() string {
		return fmt.Sprintf("CheckButton<%s %+v>", name, w.BoolVar)
	})

	w.setup()
	return w
}

// NewRadioButton creates a CheckButton bound to a string variable.
func NewRadioButton(name string, stringVar *string, value string, child Widget) *CheckButton {
	w := &CheckButton{
		StringVar: stringVar,
		Value:     value,
	}
	w.Button.child = child
	w.IDFunc(func() string {
		return fmt.Sprintf(`RadioButton<%s "%s" %s>`, name, w.Value, strconv.FormatBool(*w.StringVar == w.Value))
	})
	w.setup()
	return w
}

// Compute to re-evaluate the button state (in the case of radio buttons where
// a different button will affect the state of this one when clicked).
func (w *CheckButton) Compute(e render.Engine) {
	if w.StringVar != nil {
		// Radio button, always re-assign the border style in case a sister
		// radio button has changed the value.
		if *w.StringVar == w.Value {
			w.SetBorderStyle(BorderSunken)
		} else {
			w.SetBorderStyle(BorderRaised)
		}
	}
	w.Button.Compute(e)
}

// setup the common things between checkboxes and radioboxes.
func (w *CheckButton) setup() {
	var borderStyle BorderStyle = BorderRaised
	if w.BoolVar != nil {
		if *w.BoolVar == true {
			borderStyle = BorderSunken
		}
	}

	w.Configure(Config{
		BorderSize:   2,
		BorderStyle:  borderStyle,
		OutlineSize:  1,
		OutlineColor: theme.ButtonOutlineColor,
		Background:   theme.ButtonBackgroundColor,
	})

	w.Handle(MouseOver, func(ed EventData) error {
		w.hovering = true
		w.SetBackground(theme.ButtonHoverColor)
		return nil
	})
	w.Handle(MouseOut, func(ed EventData) error {
		w.hovering = false
		w.SetBackground(theme.ButtonBackgroundColor)
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
		var sunken bool
		if w.BoolVar != nil {
			if *w.BoolVar {
				*w.BoolVar = false
			} else {
				*w.BoolVar = true
				sunken = true
			}
		} else if w.StringVar != nil {
			*w.StringVar = w.Value
			sunken = true
		}

		if sunken {
			w.SetBorderStyle(BorderSunken)
		} else {
			w.SetBorderStyle(BorderRaised)
		}
		return nil
	})
}
