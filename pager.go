package ui

import (
	"fmt"
	"strconv"

	"git.kirsle.net/go/render"
)

// Pager is a frame with Pagers for paginated UI.
type Pager struct {
	BaseWidget

	// Config settings.
	Page     int // default 1
	Pages    int
	PerPage  int // default 20
	Font     render.Text
	OnChange func(page, perPage int)

	supervisor *Supervisor
	child      Widget
	buttons    []Widget
	page       string // radio button value of Page

	// Private options.
	hovering bool
	clicked  bool
}

// NewPager creates a new Pager.
func NewPager(config Pager) *Pager {
	w := &Pager{
		Page:     config.Page,
		Pages:    config.Pages,
		PerPage:  config.PerPage,
		Font:     config.Font,
		OnChange: config.OnChange,
		buttons:  []Widget{},
	}

	// default settings
	if w.Page == 0 {
		w.Page = 1
	}
	if w.PerPage == 0 {
		w.PerPage = 20
	}

	w.IDFunc(func() string {
		return fmt.Sprintf("Pager<%d of %d>", w.Page, w.PerPage)
	})

	w.child = w.setup()

	return w
}

// Supervise the pager to make its buttons work.
func (w *Pager) Supervise(s *Supervisor) {
	w.supervisor = s
	for _, btn := range w.buttons {
		w.supervisor.Add(btn)
	}
}

// setup the frame
func (w *Pager) setup() *Frame {
	frame := NewFrame("Pager Frame")
	frame.SetParent(w)

	if w.Pages == 0 {
		return frame
	}

	w.buttons = []Widget{}
	w.page = fmt.Sprintf("%d", w.Page)

	// Previous Page Button
	prev := NewButton("Previous", NewLabel(Label{
		Text: "<",
		Font: w.Font,
	}))
	w.buttons = append(w.buttons, prev)
	prev.Handle(Click, func(ed EventData) error {
		return w.next(-1)
	})
	frame.Pack(prev, Pack{
		Side: W,
	})

	// Draw the numbered buttons.
	for i := 1; i <= w.Pages; i++ {
		page := fmt.Sprintf("%d", i)

		btn := NewRadioButton(
			"Page "+page,
			&w.page,
			page,
			NewLabel(Label{
				Text: page,
				Font: w.Font,
			}))
		w.buttons = append(w.buttons, btn)

		btn.Handle(Click, func(ed EventData) error {
			if w.OnChange != nil {
				page, _ := strconv.Atoi(w.page)
				w.OnChange(page, w.PerPage)
			}
			return nil
		})

		if w.supervisor != nil {
			w.supervisor.Add(btn)
		}

		frame.Pack(btn, Pack{
			Side: W,
		})
	}

	// Next Page Button
	next := NewButton("Next", NewLabel(Label{
		Text: ">",
		Font: w.Font,
	}))
	w.buttons = append(w.buttons, next)
	next.Handle(Click, func(ed EventData) error {
		return w.next(1)
	})
	frame.Pack(next, Pack{
		Side: W,
	})

	return frame
}

// next (1) or previous (-1) button
func (w *Pager) next(value int) error {
	fmt.Printf("next(%d)\n", value)
	intvalue, _ := strconv.Atoi(w.page)
	intvalue += value

	if intvalue < 1 {
		intvalue = 1
	} else if intvalue > w.Pages {
		intvalue = w.Pages
	}

	w.page = fmt.Sprintf("%d", intvalue)

	if w.OnChange != nil {
		w.OnChange(intvalue, w.PerPage)
	}

	return nil
}

// Compute the size of the Pager.
func (w *Pager) Compute(e render.Engine) {
	// Compute the size of the inner widget first.
	w.child.Compute(e)

	// Auto-resize only if we haven't been given a fixed size.
	if !w.FixedSize() {
		size := w.child.Size()
		w.Resize(render.Rect{
			W: size.W + w.BoxThickness(2),
			H: size.H + w.BoxThickness(2),
		})
	}

	w.BaseWidget.Compute(e)
}

// Present the Pager.
func (w *Pager) Present(e render.Engine, P render.Point) {
	if w.Hidden() {
		return
	}

	w.Compute(e)
	var (
		S         = w.Size()
		ChildSize = w.child.Size()
	)

	// Draw the widget's border and everything.
	w.DrawBox(e, P)

	// Offset further if we are currently sunken.
	var clickOffset int
	if w.clicked {
		clickOffset++
	}

	// Where to place the child widget.
	moveTo := render.Point{
		X: P.X + w.BoxThickness(1) + clickOffset,
		Y: P.Y + w.BoxThickness(1) + clickOffset,
	}

	// If we're bigger than we need to be, center the child widget.
	if S.Bigger(ChildSize) {
		moveTo.X = P.X + (S.W / 2) - (ChildSize.W / 2)
	}

	// Draw the text label inside.
	w.child.Present(e, moveTo)

	w.BaseWidget.Present(e, P)
}
