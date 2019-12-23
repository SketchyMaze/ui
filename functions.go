package ui

import "git.kirsle.net/go/render"

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
func AbsoluteRect(w Widget) render.Rect {
	var (
		P = AbsolutePosition(w)
		R = w.Rect()
	)
	return render.Rect{
		X: P.X,
		Y: P.Y,
		W: R.W + P.X,
		H: R.H, // TODO: the Canvas in EditMode lets you draw pixels
		// below the status bar if we do `+ R.Y` here.
	}
}
