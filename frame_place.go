package ui

import (
	"git.kirsle.net/go/render"
)

// Place provides configuration fields for Frame.Place().
type Place struct {
	// X and Y coordinates for explicit location of widget within its parent.
	// This placement option trumps all others.
	Point render.Point

	// Place relative to an edge of the window. The widget will stick to the
	// edge of the window even as it resizes. Options are ignored if Point
	// is set.
	Top    int
	Left   int
	Right  int
	Bottom int
	Center bool
	Middle bool
}

// Strategy returns the placement strategy for a Place config struct.
// Returns 'Point' if a render.Point is used (even if zero, zero)
// Returns 'Side' if the side values are set.
func (p Place) Strategy() string {
	if p.Top != 0 || p.Left != 0 || p.Right != 0 || p.Bottom != 0 || p.Center || p.Middle {
		return "Side"
	}
	return "Point"
}

// placedWidget holds the data for a widget placed in a frame.
type placedWidget struct {
	widget Widget
	place  Place
}

// Place a widget into the frame.
func (w *Frame) Place(child Widget, config Place) {
	w.placed = append(w.placed, placedWidget{
		widget: child,
		place:  config,
	})
	w.Add(child)

	// Adopt the child widget so it can access the Frame.
	child.SetParent(w)
}

// computePlaced processes all the Place layout widgets in the Frame,
// determining their X,Y location and whether they need to change.
func (w *Frame) computePlaced(e render.Engine) {
	var (
		frameSize = w.BoxSize()
		// maxWidth int
		// maxHeight int
	)

	for _, row := range w.placed {
		// X,Y placement takes priority.
		switch row.place.Strategy() {
		case "Point":
			row.widget.MoveTo(row.place.Point)
			row.widget.Compute(e)
		case "Side":
			var moveTo render.Point

			// Compute the initial X,Y based on Top, Left, Right, Bottom.
			if row.place.Left > 0 {
				moveTo.X = row.place.Left
			}
			if row.place.Top > 0 {
				moveTo.Y = row.place.Top
			}
			if row.place.Right > 0 {
				moveTo.X = frameSize.W - row.widget.Size().W - row.place.Right
			}
			if row.place.Bottom > 0 {
				moveTo.Y = frameSize.H - row.widget.Size().H - row.place.Bottom
			}

			// Center and Middle aligned values override Left/Right, Top/Bottom
			// settings respectively.
			if row.place.Center {
				moveTo.X = frameSize.W - (w.Size().W / 2) - (row.widget.Size().W / 2)
			}
			if row.place.Middle {
				moveTo.Y = frameSize.H - (w.Size().H / 2) - (row.widget.Size().H / 2)
			}
			row.widget.MoveTo(moveTo)
			row.widget.Compute(e)
		}

		// If this widget itself has placed widgets, call its function too.
		if frame, ok := row.widget.(*Frame); ok {
			frame.computePlaced(e)
		}
	}
}
