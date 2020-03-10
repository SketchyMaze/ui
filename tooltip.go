package ui

import (
	"fmt"
	"strings"

	"git.kirsle.net/go/render"
)

func init() {
	precomputeArrows()
}

// Tooltip attaches a mouse-over popup to another widget.
type Tooltip struct {
	BaseWidget

	// Configurable attributes.
	Text         string  // Text to show in the tooltip.
	TextVariable *string // String pointer instead of text.
	Edge         Edge    // side to display tooltip on

	target     Widget
	lineHeight int
	font       render.Text
}

// Constants for tooltips.
const (
	tooltipArrowSize = 5
)

// NewTooltip creates a new tooltip attached to a widget.
func NewTooltip(target Widget, tt Tooltip) *Tooltip {
	w := &Tooltip{
		Text:         tt.Text,
		TextVariable: tt.TextVariable,
		Edge:         tt.Edge,
		target:       target,
	}

	// Default style.
	w.Hide()
	w.SetBackground(render.RGBA(0, 0, 0, 230))
	w.font = render.Text{
		Size:    10,
		Color:   render.White,
		Padding: 4,
	}

	// Add event bindings to the target widget.
	// - Show the tooltip on MouseOver
	// - Hide it on MouseOut
	// - Compute the tooltip when the parent widget Computes
	// - Present the tooltip when the parent widget Presents
	target.Handle(MouseOver, func(ed EventData) {
		w.Show()
	})
	target.Handle(MouseOut, func(ed EventData) {
		w.Hide()
	})
	target.Handle(Compute, func(ed EventData) {
		w.Compute(ed.Engine)
	})
	target.Handle(Present, func(ed EventData) {
		w.Present(ed.Engine, w.Point())
	})

	w.IDFunc(func() string {
		return fmt.Sprintf(`Tooltip<"%s">`, w.Value())
	})

	return w
}

// Value returns the current text displayed in the tooltop, whether from the
// configured Text or the TextVariable pointer.
func (w *Tooltip) Value() string {
	return w.text().Text
}

// text returns the raw render.Text holding the current value to be displayed
// in the tooltip, either from Text or TextVariable.
func (w *Tooltip) text() render.Text {
	if w.TextVariable != nil {
		w.font.Text = *w.TextVariable
	} else {
		w.font.Text = w.Text
	}
	return w.font
}

// Compute the size of the tooltip.
func (w *Tooltip) Compute(e render.Engine) {
	// Compute the size based on the text.
	w.computeText(e)

	// Compute the position based on the Edge and the target widget.
	var (
		size = w.Size()

		target = w.target
		tSize  = target.Size()
		tPoint = AbsolutePosition(target)

		moveTo render.Point
	)

	switch w.Edge {
	case Top:
		moveTo.Y = tPoint.Y - size.H - tooltipArrowSize
		moveTo.X = tPoint.X + (tSize.W / 2) - (size.W / 2)
	case Left:
		moveTo.X = tPoint.X - size.W - tooltipArrowSize
		moveTo.Y = tPoint.Y + (tSize.H / 2) - (size.H / 2)
	case Right:
		moveTo.X = tPoint.X + tSize.W + tooltipArrowSize
		moveTo.Y = tPoint.Y + (tSize.H / 2) - (size.H / 2)
	case Bottom:
		moveTo.Y = tPoint.Y + tSize.H + tooltipArrowSize
		moveTo.X = tPoint.X + (tSize.W / 2) - (size.W / 2)
	}

	// Adjust to keep the tooltip from clipping outside the window boundaries.
	{
		width, height := e.WindowSize()
		if moveTo.X < 0 {
			moveTo.X = 0
		} else if moveTo.X+size.W > width {
			moveTo.X = width - size.W
		}

		if moveTo.Y < 0 {
			moveTo.Y = 0
		} else if moveTo.Y+size.H > height {
			moveTo.Y = height - size.H
		}
	}

	w.MoveTo(moveTo)
}

// computeText handles the text compute, very similar to Label.Compute.
func (w *Tooltip) computeText(e render.Engine) {
	text := w.text()
	lines := strings.Split(text.Text, "\n")

	// Max rect to encompass all lines of text.
	var maxRect = render.Rect{}
	for _, line := range lines {
		if line == "" {
			line = "<empty>"
		}

		text.Text = line // only this line at this time.
		rect, err := e.ComputeTextRect(text)
		if err != nil {
			panic(fmt.Sprintf("%s: failed to compute text rect: %s", w, err)) // TODO return an error
		}

		if rect.W > maxRect.W {
			maxRect.W = rect.W
		}
		maxRect.H += rect.H
		w.lineHeight = int(rect.H)
	}

	var (
		padX = w.font.Padding + w.font.PadX
		padY = w.font.Padding + w.font.PadY
	)

	w.Resize(render.Rect{
		W: maxRect.W + (padX * 2),
		H: maxRect.H + (padY * 2),
	})
}

// Present the tooltip.
func (w *Tooltip) Present(e render.Engine, P render.Point) {
	if w.Hidden() {
		return
	}

	// Draw the text.
	w.presentText(e, P)

	// Draw the arrow.
	w.presentArrow(e, P)
}

// presentText draws the text similar to Label.
func (w *Tooltip) presentText(e render.Engine, P render.Point) {
	var (
		text = w.text()
		padX = w.font.Padding + w.font.PadX
		padY = w.font.Padding + w.font.PadY
	)

	w.DrawBox(e, P)
	for i, line := range strings.Split(text.Text, "\n") {
		text.Text = line
		e.DrawText(text, render.Point{
			X: P.X + padX,
			Y: P.Y + padY + (i * w.lineHeight),
		})
	}
}

// presentArrow draws the arrow between the tooltip and its target widget.
func (w *Tooltip) presentArrow(e render.Engine, P render.Point) {
	var (
		// size = w.Size()

		target = w.target
		tSize  = target.Size()
		tPoint = AbsolutePosition(target)

		drawAt render.Point
		arrow  [][]render.Point
	)

	switch w.Edge {
	case Top:
		arrow = arrowDown
		drawAt = render.Point{
			X: tPoint.X + (tSize.W / 2) - tooltipArrowSize,
			Y: tPoint.Y - tooltipArrowSize,
		}
	case Bottom:
		arrow = arrowUp
		drawAt = render.Point{
			X: tPoint.X + (tSize.W / 2) - tooltipArrowSize,
			Y: tPoint.Y + tSize.H,
		}
	case Left:
		arrow = arrowRight
		drawAt = render.Point{
			X: tPoint.X - tooltipArrowSize,
			Y: tPoint.Y + (tSize.H / 2) - tooltipArrowSize,
		}
	case Right:
		arrow = arrowLeft
		drawAt = render.Point{
			X: tPoint.X + tSize.W,
			Y: tPoint.Y + (tSize.H / 2) - tooltipArrowSize,
		}
	}
	drawArrow(e, w.Background(), drawAt, arrow)
}

// Draw an arrow at a given top/left coordinate.
func drawArrow(e render.Engine, color render.Color, p render.Point, arrow [][]render.Point) {
	for _, row := range arrow {
		if len(row) == 1 {
			point := render.NewPoint(row[0].X, row[0].Y)
			point.Add(p)
			e.DrawPoint(color, point)
		} else {
			start := render.NewPoint(row[0].X, row[0].Y)
			end := render.NewPoint(row[1].X, row[1].Y)
			start.Add(p)
			end.Add(p)
			e.DrawLine(color, start, end)
		}
	}
}

// Arrows for the tooltip widget.
var (
	arrowDown  [][]render.Point
	arrowUp    [][]render.Point
	arrowLeft  [][]render.Point
	arrowRight [][]render.Point
)

func precomputeArrows() {
	arrowDown = [][]render.Point{
		{render.NewPoint(0, 0), render.NewPoint(10, 0)},
		{render.NewPoint(1, 1), render.NewPoint(9, 1)},
		{render.NewPoint(2, 2), render.NewPoint(8, 2)},
		{render.NewPoint(3, 3), render.NewPoint(7, 3)},
		{render.NewPoint(4, 4), render.NewPoint(6, 4)},
		{render.NewPoint(5, 5)},
	}
	arrowUp = [][]render.Point{
		{render.NewPoint(5, 0)},
		{render.NewPoint(4, 1), render.NewPoint(6, 1)},
		{render.NewPoint(3, 2), render.NewPoint(7, 2)},
		{render.NewPoint(2, 3), render.NewPoint(8, 3)},
		{render.NewPoint(1, 4), render.NewPoint(9, 4)},
		// {render.NewPoint(0, 5), render.NewPoint(10, 5)},
	}
	arrowLeft = [][]render.Point{
		{render.NewPoint(0, 5)},
		{render.NewPoint(1, 4), render.NewPoint(1, 6)},
		{render.NewPoint(2, 3), render.NewPoint(2, 7)},
		{render.NewPoint(3, 2), render.NewPoint(3, 8)},
		{render.NewPoint(4, 1), render.NewPoint(4, 9)},
		// {render.NewPoint(5, 0), render.NewPoint(5, 10)},
	}
	arrowRight = [][]render.Point{
		{render.NewPoint(0, 0), render.NewPoint(0, 10)},
		{render.NewPoint(1, 1), render.NewPoint(1, 9)},
		{render.NewPoint(2, 2), render.NewPoint(2, 8)},
		{render.NewPoint(3, 3), render.NewPoint(3, 7)},
		{render.NewPoint(4, 4), render.NewPoint(4, 6)},
		{render.NewPoint(5, 5)},
	}
}
