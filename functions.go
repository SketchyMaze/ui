package ui

import (
	"git.kirsle.net/go/render"
)

// AbsolutePosition computes a widget's absolute X,Y position on the
// window on screen by crawling its parent widget tree.
func AbsolutePosition(w Widget) render.Point {
	abs := w.Point()

	var (
		node = w
		ok   bool
	)

	for {
		node, ok = node.Parent()
		if !ok { // reached the top of the tree
			return abs
		}

		abs.Add(node.Point())
	}
}

// AbsoluteRect returns a Rect() offset with the absolute position.
// X and Y are the AbsolutePosition of the widget.
// W and H are the widget's width and height. (X,Y not added to them)
func AbsoluteRect(w Widget) render.Rect {
	var (
		P = AbsolutePosition(w)
		R = w.Rect()
	)
	return render.Rect{
		X: P.X,
		Y: P.Y,
		W: R.W,
		H: R.H, // TODO: the Canvas in EditMode lets you draw pixels
		// below the status bar if we do `+ R.Y` here.
	}
}

// widgetInFocusedWindow returns whether a widget (like a Button) is a
// descendant of a Window that is being Window Managed by Supervisor, and
// said window is in a Focused state.
//
// This is used by Supervisor to decide whether the widget should be given
// events or not: a widget in a non-focused window ignores events, so that a
// button in a "lower" window could not be clicked through a "higher" window
// that overlaps it.
func widgetInFocusedWindow(w Widget) (isManaged, isFocused bool) {
	var node = w

	for {
		// Is the node a Window?
		if window, ok := node.(*Window); ok {
			return window.managed, window.Focused()
		}

		node, _ = node.Parent()
		if node == nil {
			return false, true // reached the root
		}
	}
}
