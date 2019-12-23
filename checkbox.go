package ui

import "git.kirsle.net/go/render"

// Checkbox combines a CheckButton with a widget like a Label.
type Checkbox struct {
	Frame
	button *CheckButton
	child  Widget
}

// NewCheckbox creates a new Checkbox.
func NewCheckbox(name string, boolVar *bool, child Widget) *Checkbox {
	return makeCheckbox(name, boolVar, nil, "", child)
}

// NewRadiobox creates a new Checkbox in radio mode.
func NewRadiobox(name string, stringVar *string, value string, child Widget) *Checkbox {
	return makeCheckbox(name, nil, stringVar, value, child)
}

// makeCheckbox constructs an appropriate type of checkbox.
func makeCheckbox(name string, boolVar *bool, stringVar *string, value string, child Widget) *Checkbox {
	// Our custom checkbutton widget.
	mark := NewFrame(name + "_mark")

	w := &Checkbox{
		child: child,
	}
	if boolVar != nil {
		w.button = NewCheckButton(name+"_button", boolVar, mark)
	} else if stringVar != nil {
		w.button = NewRadioButton(name+"_button", stringVar, value, mark)
	}
	w.Frame.Setup()

	// Forward clicks on the child widget to the CheckButton.
	for _, e := range []Event{MouseOver, MouseOut, MouseUp, MouseDown} {
		func(e Event) {
			w.child.Handle(e, func(p render.Point) {
				w.button.Event(e, p)
			})
		}(e)
	}

	w.Pack(w.button, Pack{
		Anchor: W,
	})
	w.Pack(w.child, Pack{
		Anchor: W,
	})

	return w
}

// Child returns the child widget.
func (w *Checkbox) Child() Widget {
	return w.child
}

// Supervise the checkbutton inside the widget.
func (w *Checkbox) Supervise(s *Supervisor) {
	s.Add(w.button)
	s.Add(w.child)
}
