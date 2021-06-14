package ui

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
	mark.Configure(Config{
		Width:  6,
		Height: 6,
	})

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
	for _, e := range []Event{MouseOver, MouseOut, MouseUp, MouseDown, Click} {
		func(e Event) {
			w.child.Handle(e, func(ed EventData) error {
				return w.button.Event(e, ed)
			})
		}(e)
	}

	w.Pack(w.button, Pack{
		Side: W,
	})
	w.Pack(w.child, Pack{
		Side: W,
	})

	return w
}

// Child returns the child widget.
func (w *Checkbox) Child() Widget {
	return w.child
}

// Pass event handlers on to descendents.
func (w *Checkbox) Handle(e Event, fn func(EventData) error) {
	w.button.Handle(e, fn)
}

// Supervise the checkbutton inside the widget.
func (w *Checkbox) Supervise(s *Supervisor) {
	s.Add(w.button)
	s.Add(w.child)
}
