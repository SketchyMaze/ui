package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/style"
)

// Scrollbar dimensions, TODO: make configurable.
var (
	scrollWidth     = 20
	scrollbarHeight = 40
)

// ScrollBar is a classic scrolling widget.
type ScrollBar struct {
	*Frame
	style      *style.Button
	supervisor *Supervisor

	trough *Frame
	slider *Frame

	// Configurable scroll ranges.
	Min   int
	Max   int
	Step  int
	value int

	// Variable bindings: give these pointers to your values.
	Variable interface{} // pointer to e.g. a string or int
	// TextVariable *string // string value
	// IntVariable  *int    // integer value

	// Drag/drop state.
	dragging    bool         // mouse down on slider
	scrollPx    int          // px from the top where the slider is placed
	dragStart   render.Point // where the mouse was on click
	wasScrollPx int

	everyTick func()
}

// NewScrollBar creates a new ScrollBar.
func NewScrollBar(config ScrollBar) *ScrollBar {
	w := &ScrollBar{
		Frame:    NewFrame("Scrollbar Frame"),
		Variable: config.Variable,
		style:    &style.DefaultButton,
		Min:      config.Min,
		Max:      config.Max,
		Step:     config.Step,
	}

	if w.Max == 0 {
		w.Max = 100
	}
	if w.Step == 0 {
		w.Step = 1
	}

	w.IDFunc(func() string {
		return "ScrollBar"
	})

	w.SetStyle(Theme.Button)

	w.setup()
	return w
}

// SetStyle sets the ScrollBar style.
func (w *ScrollBar) SetStyle(v *style.Button) {
	if v == nil {
		v = &style.DefaultButton
	}

	w.style = v
	fmt.Printf("set style: %+v\n", v)
	w.Frame.Configure(Config{
		BorderSize:  w.style.BorderSize,
		BorderStyle: BorderSunken,
		Background:  w.style.Background.Darken(40),
	})
}

// GetStyle gets the ScrollBar style.
func (w *ScrollBar) GetStyle() *style.Button {
	return w.style
}

// Supervise the ScrollBar. This is necessary for granting mouse-over events
// to the items in the list.
func (w *ScrollBar) Supervise(s *Supervisor) {
	w.supervisor = s

	// Add all the list items to be supervised.
	w.supervisor.Add(w.slider)
	for _, c := range w.Frame.Children() {
		w.supervisor.Add(c)
	}
}

// Compute to re-evaluate the button state (in the case of radio buttons where
// a different button will affect the state of this one when clicked).
func (w *ScrollBar) Compute(e render.Engine) {
	w.Frame.Compute(e)
	if w.everyTick != nil {
		w.everyTick()
	}
}

// setup the UI components and event handlers.
func (w *ScrollBar) setup() {
	w.Configure(Config{
		Width: scrollWidth,
	})

	// The trough that holds the slider.
	w.trough = NewFrame("Trough")

	// Up button
	upBtn := NewButton("Up", NewLabel(Label{
		Text: "^",
	}))
	upBtn.Handle(MouseDown, func(ed EventData) error {
		w.everyTick = func() {
			w.scrollPx -= w.Step
			if w.scrollPx < 0 {
				w.scrollPx = 0
			}
			w.trough.Place(w.slider, Place{
				Top: w.scrollPx,
			})
			w.sendScrollEvent()
		}
		return nil
	})
	upBtn.Handle(MouseUp, func(ed EventData) error {
		w.everyTick = nil
		return nil
	})

	// The slider
	w.slider = NewFrame("Slider")
	w.slider.Configure(Config{
		BorderSize:  w.style.BorderSize,
		BorderStyle: BorderStyle(w.style.BorderStyle),
		Background:  w.style.Background,
		Width:       scrollWidth - w.BoxThickness(w.style.BorderSize),
		Height:      scrollbarHeight,
	})

	// Slider events
	w.slider.Handle(MouseOver, func(ed EventData) error {
		w.slider.SetBackground(w.style.HoverBackground)
		return nil
	})
	w.slider.Handle(MouseOut, func(ed EventData) error {
		w.slider.SetBackground(w.style.Background)
		return nil
	})
	w.slider.Handle(MouseDown, func(ed EventData) error {
		w.dragging = true
		w.dragStart = ed.Point
		w.wasScrollPx = w.scrollPx
		fmt.Printf("begin drag from %s\n", ed.Point)
		return nil
	})
	w.slider.Handle(MouseUp, func(ed EventData) error {
		fmt.Println("mouse released")
		w.dragging = false
		return nil
	})
	w.slider.Handle(MouseMove, func(ed EventData) error {
		if w.dragging {
			var (
				delta  = w.dragStart.Compare(ed.Point)
				moveTo = w.wasScrollPx + delta.Y
			)

			if moveTo < 0 {
				moveTo = 0
			} else if moveTo > w.trough.height-w.slider.height {
				moveTo = w.trough.height - w.slider.height
			}

			fmt.Printf("delta drag: %s\n", delta)
			w.scrollPx = moveTo
			w.trough.Place(w.slider, Place{
				Top: w.scrollPx,
			})
			w.sendScrollEvent()
		}
		return nil
	})

	downBtn := NewButton("Down", NewLabel(Label{
		Text: "v",
	}))
	downBtn.Handle(MouseDown, func(ed EventData) error {
		w.everyTick = func() {
			w.scrollPx += w.Step
			if w.scrollPx > w.trough.height-w.slider.height {
				w.scrollPx = w.trough.height - w.slider.height
			}
			w.trough.Place(w.slider, Place{
				Top: w.scrollPx,
			})
			w.sendScrollEvent()
		}
		return nil
	})
	downBtn.Handle(MouseUp, func(ed EventData) error {
		w.everyTick = nil
		return nil
	})

	w.Frame.Pack(upBtn, Pack{
		Side:  N,
		FillX: true,
	})
	w.Frame.Pack(w.trough, Pack{
		Side:   N,
		Fill:   true,
		Expand: true,
	})
	w.trough.Place(w.slider, Place{
		Top:  w.scrollPx,
		Left: 0,
	})
	w.Frame.Pack(downBtn, Pack{
		Side:  N,
		FillX: true,
	})
}

// Present the scrollbar.
func (w *ScrollBar) Present(e render.Engine, p render.Point) {
	w.Frame.Present(e, p)
}

func (w *ScrollBar) sendScrollEvent() {
	var fraction float64
	if w.scrollPx > 0 {
		fraction = float64(w.scrollPx) / (float64(w.trough.height) - float64(w.slider.height))
	}
	w.Event(Scroll, EventData{
		ScrollFraction: fraction,
		ScrollUnits:    int(fraction * float64(w.Max)),
		ScrollPages:    0,
	})
}
