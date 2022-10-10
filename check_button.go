package ui

import (
	"fmt"
	"strconv"

	"git.kirsle.net/go/render"
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

	w.SetStyle(Theme.Button)

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

	w.SetStyle(Theme.Button)

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
	} else if w.BoolVar != nil {
		// Checkbutton, always re-assign the border style in case the caller
		// has flipped the boolean behind our back.
		if *w.BoolVar {
			w.SetBorderStyle(BorderSunken)
		} else {
			w.SetBorderStyle(BorderRaised)
		}
	}
	w.Button.Compute(e)
}

// setup the common things between checkboxes and radioboxes.
func (w *CheckButton) setup() {
	var (
		borderStyle BorderStyle = BorderRaised
		background              = w.style.Background
	)
	if w.BoolVar != nil {
		if *w.BoolVar == true {
			borderStyle = BorderSunken
			background = w.style.Background.Darken(40)
		}
	}

	w.Configure(Config{
		BorderSize:   2,
		BorderStyle:  borderStyle,
		OutlineSize:  1,
		OutlineColor: w.style.OutlineColor,
		Background:   background,
	})

	w.Handle(MouseOver, func(ed EventData) error {
		w.hovering = true
		w.SetBackground(w.style.HoverBackground)
		return nil
	})
	w.Handle(MouseOut, func(ed EventData) error {
		w.hovering = false

		var sunken bool
		if w.BoolVar != nil {
			sunken = *w.BoolVar == true
		} else if w.StringVar != nil {
			sunken = *w.StringVar == w.Value
		}

		if sunken {
			w.SetBackground(w.style.Background.Darken(40))
		} else {
			w.SetBackground(w.style.Background)
		}

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
			w.SetBackground(w.style.Background.Darken(40))
		} else {
			w.SetBorderStyle(BorderRaised)
			w.SetBackground(w.style.Background)
		}

		return nil
	})
}
