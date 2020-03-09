package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
)

// Frame is a widget that contains other widgets.
type Frame struct {
	Name string
	BaseWidget

	// Widget placement settings.
	packs   map[Side][]packedWidget // Packed widgets
	placed  []placedWidget          // Placed widgets
	widgets []Widget
}

// NewFrame creates a new Frame.
func NewFrame(name string) *Frame {
	w := &Frame{
		Name:    name,
		packs:   map[Side][]packedWidget{},
		widgets: []Widget{},
	}
	w.SetBackground(render.RGBA(1, 0, 0, 0)) // invisible default BG
	w.IDFunc(func() string {
		return fmt.Sprintf("Frame<%s>",
			name,
		)
	})
	return w
}

// Setup ensures all the Frame's data is initialized and not null.
func (w *Frame) Setup() {
	if w.packs == nil {
		w.packs = map[Side][]packedWidget{}
	}
	if w.widgets == nil {
		w.widgets = []Widget{}
	}
}

// Children returns all of the child widgets.
func (w *Frame) Children() []Widget {
	return w.widgets
}

// Compute the size of the Frame.
func (w *Frame) Compute(e render.Engine) {
	w.computePacked(e)
	w.computePlaced(e)
}

// Present the Frame.
func (w *Frame) Present(e render.Engine, P render.Point) {
	if w.Hidden() {
		return
	}

	var (
		S = w.Size()
	)

	// Draw the widget's border and everything.
	w.DrawBox(e, P)

	// Draw the background color.
	e.DrawBox(w.Background(), render.Rect{
		X: P.X + w.BoxThickness(1),
		Y: P.Y + w.BoxThickness(1),
		W: S.W - w.BoxThickness(2),
		H: S.H - w.BoxThickness(2),
	})

	// Draw the widgets.
	for _, child := range w.widgets {
		// child.Compute(e)
		p := child.Point()
		moveTo := render.NewPoint(
			P.X+p.X+w.BoxThickness(1),
			P.Y+p.Y+w.BoxThickness(1),
		)
		// if child.ID() == "Canvas" {
		// 	log.Debug("Frame X=%d  Child X=%d  Box=%d  Point=%s", P.X, p.X, w.BoxThickness(1), p)
		// 	log.Debug("Frame Y=%d  Child Y=%d  Box=%d  MoveTo=%s", P.Y, p.Y, w.BoxThickness(1), moveTo)
		// }
		// child.MoveTo(moveTo) // TODO: if uncommented the child will creep down the parent each tick
		// if child.ID() == "Canvas" {
		// 	log.Debug("New Point: %s", child.Point())
		// }
		child.Present(e, moveTo)
	}
}
