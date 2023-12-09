package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/style"
)

// ListBox is a selectable list of values like a multi-line SelectBox.
type ListBox struct {
	*Frame
	name       string
	children   []*ListValue
	style      *style.ListBox
	supervisor *Supervisor

	list           *Frame
	scrollbar      *ScrollBar
	scrollFraction float64
	maxHeight      int

	// Variable bindings: give these pointers to your values.
	Variable interface{} // pointer to e.g. a string or int
	// TextVariable *string // string value
	// IntVariable  *int    // integer value
}

// ListValue is an item in the ListBox. It has an arbitrary widget as a
// "label" (usually a Label) and a value (string or int) when it's "selected"
type ListValue struct {
	Frame *Frame
	Label Widget
	Value interface{}
}

// NewListBox creates a new ListBox.
func NewListBox(name string, config ListBox) *ListBox {
	w := &ListBox{
		Frame:    NewFrame(name + " Frame"),
		list:     NewFrame(name + " List"),
		name:     name,
		children: []*ListValue{},
		Variable: config.Variable,
		// TextVariable: config.TextVariable,
		// IntVariable:  config.IntVariable,
		style: &style.DefaultListBox,
	}

	// if config.Width > 0 && config.Height > 0 {
	// 	w.Frame.Resize(render.NewRect(config.Width, config.Height))
	// }

	w.IDFunc(func() string {
		return fmt.Sprintf("ListBox<%s>", name)
	})

	w.SetStyle(Theme.ListBox)

	w.setup()
	return w
}

// SetStyle sets the listbox style.
func (w *ListBox) SetStyle(v *style.ListBox) {
	if v == nil {
		v = &style.DefaultListBox
	}

	w.style = v
	fmt.Printf("set style: %+v\n", v)
	w.Frame.Configure(Config{
		BorderSize:  w.style.BorderSize,
		BorderStyle: BorderStyle(w.style.BorderStyle),
		Background:  w.style.Background,
	})

	// If the child is a Label, apply the foreground color.
	// if label, ok := w.child.(*Label); ok {
	// 	label.Font.Color = w.style.Foreground
	// }
}

// GetStyle gets the listbox style.
func (w *ListBox) GetStyle() *style.ListBox {
	return w.style
}

// Supervise the ListBox. This is necessary for granting mouse-over events
// to the items in the list.
func (w *ListBox) Supervise(s *Supervisor) {
	w.supervisor = s
	w.scrollbar.Supervise(s)

	// Add all the list items to be supervised.
	for _, c := range w.children {
		w.supervisor.Add(c.Frame)
	}
}

// AddLabel adds a simple text-based label to the Listbox.
// The label is the text value to display.
// The value is the underlying value (string or int) for the TextVariable or IntVariable.
// The function callback runs when the option is picked.
func (w *ListBox) AddLabel(label string, value interface{}, f func()) {
	row := NewFrame(label + " Frame")

	child := NewLabel(Label{
		Text: label,
		Font: render.Text{
			Color:   w.style.Foreground,
			Size:    11,
			Padding: 2,
		},
	})
	row.Pack(child, Pack{
		Side:  W,
		FillX: true,
	})

	// Add this label and its value mapping to the ListBox.
	w.children = append(w.children, &ListValue{
		Frame: row,
		Label: child,
		Value: value,
	})

	// Event handlers for the item row.
	// row.Handle(MouseOver, func(ed EventData) error {
	// 	if ed.Point.Inside(AbsoluteRect(w.scrollbar)) {
	// 		return nil // ignore if over scrollbar
	// 	}

	// 	row.SetBackground(w.style.HoverBackground)
	// 	child.Font.Color = w.style.HoverForeground
	// 	return nil
	// })
	row.Handle(MouseMove, func(ed EventData) error {
		if ed.Point.Inside(AbsoluteRect(w.scrollbar)) {
			// we wandered onto the scrollbar, cancel mouseover
			return row.Event(MouseOut, ed)
		}
		row.SetBackground(w.style.HoverBackground)
		child.Font.Color = w.style.HoverForeground
		return nil
	})
	row.Handle(MouseOut, func(ed EventData) error {
		if cur, ok := w.GetValue(); ok && cur == value {
			row.SetBackground(w.style.SelectedBackground)
			child.Font.Color = w.style.SelectedForeground
		} else {
			fmt.Printf("couldn't get value? %+v %+v\n", cur, ok)
			row.SetBackground(w.style.Background)
			child.Font.Color = w.style.Foreground
		}
		return nil
	})
	row.Handle(MouseUp, func(ed EventData) error {
		if cur, ok := w.GetValue(); ok && cur == value {
			row.SetBackground(w.style.SelectedBackground)
			child.Font.Color = w.style.SelectedForeground
		} else {
			row.SetBackground(w.style.Background)
			child.Font.Color = w.style.Foreground
		}
		return nil
	})

	row.Handle(Click, func(ed EventData) error {
		// Trigger if we are not hovering over the (overlapping) scrollbar.
		if !ed.Point.Inside(AbsoluteRect(w.scrollbar)) {
			w.Event(Change, EventData{
				Supervisor: w.supervisor,
				Value:      value,
			})
		}
		return nil
	})

	// Append the item into the ListBox frame.
	w.Frame.Pack(row, Pack{
		Side: N,
		PadY: 1,
		Fill: true,
	})

	// If the current text label isn't in the options, pick
	// the first option.
	if _, ok := w.GetValue(); !ok {
		w.Variable = w.children[0].Value
		row.SetBackground(w.style.SelectedBackground)
	}
}

// TODO: RemoveItem()

// GetValue returns the currently selected item in the ListBox.
//
// Returns the SelectValue and true on success, and the Label or underlying Value
// can be read from the SelectValue struct. If no valid option is selected, the
// bool value returns false.
func (w *ListBox) GetValue() (*ListValue, bool) {
	for _, row := range w.children {
		if w.Variable != nil && w.Variable == row.Value {
			return row, true
		}
	}
	return nil, false
}

// SetValueByLabel sets the currently selected option to the given label.
func (w *ListBox) SetValueByLabel(label string) bool {
	for _, option := range w.children {
		if child, ok := option.Label.(*Label); ok && child.Text == label {
			w.Variable = option.Value
			return true
		}
	}
	return false
}

// SetValue sets the currently selected option to the given value.
func (w *ListBox) SetValue(value interface{}) bool {
	w.Variable = value
	for _, option := range w.children {
		if option.Value == value {
			w.Variable = option.Value
			return true
		}
	}
	return false
}

// Compute to re-evaluate the button state (in the case of radio buttons where
// a different button will affect the state of this one when clicked).
func (w *ListBox) Compute(e render.Engine) {
	w.computeVisible()
	w.Frame.Compute(e)
}

// setup the UI components and event handlers.
func (w *ListBox) setup() {
	// w.Configure(Config{
	// 	BorderSize:  1,
	// 	BorderStyle: BorderSunken,
	// 	Background:  theme.InputBackgroundColor,
	// })
	w.scrollbar = NewScrollBar(ScrollBar{})
	w.scrollbar.Handle(Scroll, func(ed EventData) error {
		fmt.Printf("Scroll event: %f%% unit %d\n", ed.ScrollFraction*100, ed.ScrollUnits)
		w.scrollFraction = ed.ScrollFraction
		return nil
	})
	w.Frame.Pack(w.scrollbar, Pack{
		Side:    E,
		FillY:   true,
		Padding: 0,
	})
	// w.Frame.Pack(w.list, Pack{
	// 	Side:   E,
	// 	FillY:  true,
	// 	Expand: true,
	// })
}

// Compute which items of the list should be visible based on scroll position.
func (w *ListBox) computeVisible() {
	if len(w.children) == 0 {
		return
	}

	// Sample the first element's height.
	var (
		myHeight = w.height
		maxTop   = w.maxHeight - myHeight + w.children[len(w.children)-1].Frame.height
		top      = int(w.scrollFraction * float64(maxTop))
		// itemHeight = w.children[0].Label.Size().H
	)

	var (
		scan        int
		scrollFreed int
		totalHeight int
	)
	for _, c := range w.children {
		childHeight := c.Frame.Size().H + 2
		if top > 0 && scan+childHeight < top {
			scrollFreed += childHeight
			c.Frame.Hide()
		} else if scan+childHeight > myHeight+scrollFreed {
			c.Frame.Hide()
		} else {
			c.Frame.Show()
		}
		scan += childHeight // for padding
		totalHeight += childHeight
	}

	w.maxHeight = totalHeight
}

func (w *ListBox) Present(e render.Engine, p render.Point) {
	w.Frame.Present(e, p)

	// HACK to get the scrollbar to appear on top of the list frame :(
	pos := AbsolutePosition(w.scrollbar)
	// pos.X += w.BoxThickness(w.style.BorderSize / 2) // HACK
	w.scrollbar.Present(e, pos)
}
