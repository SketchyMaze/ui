package ui

import (
	"fmt"
	"strings"

	"git.kirsle.net/apps/doodle/lib/render"
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
	Font         render.Text

	width      int32
	height     int32
	lineHeight int
}

// NewLabel creates a new label.
func NewLabel(c Label) *Label {
	w := &Label{
		Text:         c.Text,
		TextVariable: c.TextVariable,
		Font:         DefaultFont,
	}
	if !c.Font.IsZero() {
		w.Font = c.Font
	}
	w.IDFunc(func() string {
		return fmt.Sprintf(`Label<"%s">`, w.text().Text)
	})
	return w
}

// text returns the label's displayed text, coming from the TextVariable if
// available or else the Text attribute instead.
func (w *Label) text() render.Text {
	if w.TextVariable != nil {
		w.Font.Text = *w.TextVariable
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
		text.Text = line // only this line at this time.
		rect, err := e.ComputeTextRect(text)
		if err != nil {
			panic(fmt.Sprintf("%s: failed to compute text rect: %s", w, err)) // TODO return an error
			return
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
		w.resizeAuto(render.Rect{
			W: maxRect.W + (padX * 2),
			H: maxRect.H + (padY * 2),
		})
	}

	w.MoveTo(render.Point{
		X: maxRect.X + w.BoxThickness(1),
		Y: maxRect.Y + w.BoxThickness(1),
	})
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
			Y: P.Y + border + padY + int32(i*w.lineHeight),
		})
	}
}
