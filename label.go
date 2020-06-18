package ui

import (
	"fmt"
	"strings"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/style"
)

// DefaultFont is the default font settings used for a Label.
var DefaultFont = render.Text{
	Size:  12,
	Color: render.Black,
}

// Label is a simple text label widget.
type Label struct {
	BaseWidget

	// Configurable fields for the constructor.
	Text         string
	TextVariable *string
	IntVariable  *int
	Font         render.Text

	style      *style.Label
	width      int
	height     int
	lineHeight int
}

// NewLabel creates a new label.
func NewLabel(c Label) *Label {
	w := &Label{
		Text:         c.Text,
		TextVariable: c.TextVariable,
		IntVariable:  c.IntVariable,
		Font:         DefaultFont,
	}
	w.SetStyle(Theme.Label)
	if !c.Font.IsZero() {
		w.Font = c.Font
	}
	w.IDFunc(func() string {
		return fmt.Sprintf(`Label<"%s">`, w.text().Text)
	})
	return w
}

// SetStyle sets the label's default style.
func (w *Label) SetStyle(v *style.Label) {
	if v == nil {
		v = &style.DefaultLabel
	}

	w.style = v
	w.SetBackground(w.style.Background)
	w.Font.Color = w.style.Foreground
}

// text returns the label's displayed text, coming from the TextVariable if
// available or else the Text attribute instead.
func (w *Label) text() render.Text {
	if w.TextVariable != nil {
		w.Font.Text = *w.TextVariable
		return w.Font
	} else if w.IntVariable != nil {
		w.Font.Text = fmt.Sprintf("%d", *w.IntVariable)
		return w.Font
	}
	w.Font.Text = w.Text
	return w.Font
}

// Value returns the current text value displayed in the widget, whether it was
// the hardcoded value or a TextVariable.
func (w *Label) Value() string {
	return w.text().Text
}

// Compute the size of the label widget.
func (w *Label) Compute(e render.Engine) {
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
		padX = w.Font.Padding + w.Font.PadX
		padY = w.Font.Padding + w.Font.PadY
	)

	if !w.FixedSize() {
		w.ResizeAuto(render.Rect{
			W: maxRect.W + (padX * 2),
			H: maxRect.H + (padY * 2),
		})
	}

	// Call the BaseWidget Compute in case we have subscribers.
	w.BaseWidget.Compute(e)
}

// Present the label widget.
func (w *Label) Present(e render.Engine, P render.Point) {
	if w.Hidden() {
		return
	}

	border := w.BoxThickness(1)

	var (
		text = w.text()
		padX = w.Font.Padding + w.Font.PadX
		padY = w.Font.Padding + w.Font.PadY
	)

	w.DrawBox(e, P)
	for i, line := range strings.Split(text.Text, "\n") {
		text.Text = line
		e.DrawText(text, render.Point{
			X: P.X + border + padX,
			Y: P.Y + border + padY + (i * w.lineHeight),
		})
	}

	// Call the BaseWidget Present in case we have subscribers.
	w.BaseWidget.Present(e, P)
}
