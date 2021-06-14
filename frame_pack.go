package ui

import (
	"fmt"

	"git.kirsle.net/go/render"
)

// Pack provides configuration fields for Frame.Pack().
type Pack struct {
	// Side of the parent to anchor the position to, like N, SE, W. Default
	// is Center.
	Side Side

	// If the widget is smaller than its allocated space, grow the widget
	// to fill its space in the Frame.
	Fill  bool
	FillX bool
	FillY bool

	Padding int // Equal padding on X and Y.
	PadX    int
	PadY    int
	Expand  bool // Widget should grow its allocated space to better fill the parent.
}

// Pack a widget along a side of the frame.
func (w *Frame) Pack(child Widget, config ...Pack) {
	var C Pack
	if len(config) > 0 {
		C = config[0]
	}

	// Initialize the pack list for this side?
	if _, ok := w.packs[C.Side]; !ok {
		w.packs[C.Side] = []packedWidget{}
	}

	// Padding: if the user only provided Padding add it to both
	// the X and Y value. If the user additionally provided the X
	// and Y value, it will add to the base padding as you'd expect.
	C.PadX += C.Padding
	C.PadY += C.Padding

	// Fill: true implies both directions.
	if C.Fill {
		C.FillX = true
		C.FillY = true
	}

	// Adopt the child widget so it can access the Frame.
	child.SetParent(w)

	w.packs[C.Side] = append(w.packs[C.Side], packedWidget{
		widget: child,
		pack:   C,
	})
	w.Add(child)
}

// Unpack removes the widget from the packed lists.
func (w *Frame) Unpack(child Widget) bool {
	var any = false
	for side, widgets := range w.packs {
		var (
			replace = []packedWidget{}
			found   = false
		)

		fmt.Printf("unpack:%s side:%s\n", child, side)

		for _, widget := range widgets {
			if widget.widget == child {
				fmt.Printf("found!\n")
				found = true
				any = true
				continue
			}
			replace = append(replace, widget)
		}

		if found {
			w.packs[side] = replace
		}
	}
	return any
}

// computePacked processes all the Pack layout widgets in the Frame.
func (w *Frame) computePacked(e render.Engine) {
	var (
		frameSize = w.BoxSize()

		// maxWidth and maxHeight are always the computed minimum dimensions
		// that the Frame must be to contain all of its children. If the Frame
		// was configured with an explicit Size, the Frame will be that Size,
		// but we still calculate how much space the widgets _actually_ take
		// so we can expand them to fill remaining space in fixed size Frames.
		maxWidth  int
		maxHeight int
		visited   = []packedWidget{}
		expanded  = []packedWidget{}
	)

	// Iterate through all directions and compute how much space to
	// reserve to contain all of their widgets.
	for side := SideMin; side <= SideMax; side++ {
		if _, ok := w.packs[side]; !ok {
			continue
		}

		var (
			x          int
			y          int
			yDirection = 1
			xDirection = 1
		)

		if side.IsSouth() {
			y = frameSize.H - w.BoxThickness(4)
			yDirection = -1
		} else if side.IsEast() {
			x = frameSize.W - w.BoxThickness(4)
			xDirection = -1
		}

		for _, packedWidget := range w.packs[side] {

			child := packedWidget.widget
			pack := packedWidget.pack
			child.Compute(e)

			if child.Hidden() {
				continue
			}

			x += pack.PadX * xDirection
			y += pack.PadY * yDirection

			var (
				// point = child.Point()
				size  = child.Size()
				yStep = y * yDirection
				xStep = x * xDirection
			)

			if xStep+size.W+(pack.PadX*2) > maxWidth {
				maxWidth = xStep + size.W + (pack.PadX * 2)
			}
			if yStep+size.H+(pack.PadY*2) > maxHeight {
				maxHeight = yStep + size.H + (pack.PadY * 2)
			}

			if side.IsSouth() {
				y -= size.H - pack.PadY
			}
			if side.IsEast() {
				x -= size.W - pack.PadX
			}

			// NOTE: we place the child's position relative to the Frame's
			// position. So a child placed at the top/left of the Frame gets
			// an x,y near zero regardless of the Frame's position.
			child.MoveTo(render.NewPoint(x, y))

			if side.IsNorth() {
				y += size.H + pack.PadY
			}
			if side.IsWest() {
				x += size.W + pack.PadX
			}

			visited = append(visited, packedWidget)
			if pack.Expand { // TODO: don't fuck with children of fixed size
				expanded = append(expanded, packedWidget)
			}
		}
	}

	// If we have extra space in the Frame and any expanding widgets, let the
	// expanding widgets grow and share the remaining space.
	computedSize := render.NewRect(maxWidth, maxHeight)
	if len(expanded) > 0 && !frameSize.IsZero() { // && frameSize.Bigger(computedSize) {
		// Divy up the size available.
		growBy := render.Rect{
			W: ((frameSize.W - computedSize.W) / len(expanded)) - w.BoxThickness(4),
			H: ((frameSize.H - computedSize.H) / len(expanded)) - w.BoxThickness(4),
		}
		for _, pw := range expanded {
			// Grow the widget but maintain its auto-size flag, in case the widget
			// was not given an explicit size before.
			size := pw.widget.Size()
			pw.widget.ResizeAuto(render.Rect{
				W: size.W + growBy.W,
				H: size.H + growBy.H,
			})
			pw.widget.Compute(e)
		}
	}

	// If we're not using a fixed Frame size, use the dynamically computed one.
	if !w.FixedSize() {
		frameSize = render.NewRect(maxWidth, maxHeight)
	} else {
		// If either of the sizes were left zero, use the dynamically computed one.
		if frameSize.W == 0 {
			frameSize.W = maxWidth
		}
		if frameSize.H == 0 {
			frameSize.H = maxHeight
		}
	}

	// Rescan all the widgets in this side to re-center them
	// in their space.
	innerFrameSize := render.NewRect(
		frameSize.W-w.BoxThickness(2),
		frameSize.H-w.BoxThickness(2),
	)
	for _, pw := range visited {
		var (
			child   = pw.widget
			pack    = pw.pack
			point   = child.Point()
			size    = child.Size()
			resize  = size
			resized bool
			moved   bool
		)

		if pack.Side.IsNorth() || pack.Side.IsSouth() {
			// Aligned to the top or bottom. If the widget Fills horizontally,
			// resize it so its Width matches the frame's Width.
			if pack.FillX && resize.W < innerFrameSize.W {
				resize.W = innerFrameSize.W - w.BoxThickness(2) // TODO: child.BoxThickness instead??
				resized = true
			}

			// If it does not Fill horizontally and there is extra horizontal
			// space, center the widget inside the space. TODO: Anchor option
			// could align the widget to the left or right instead of center.
			if resize.W < innerFrameSize.W-w.BoxThickness(4) {
				if pack.Side.IsCenter() {
					point.X = (innerFrameSize.W / 2) - (resize.W / 2)
				} else if pack.Side.IsWest() {
					point.X = pack.PadX
				} else if pack.Side.IsEast() {
					point.X = innerFrameSize.W - resize.W - pack.PadX
				}

				moved = true
			}
		} else if pack.Side.IsWest() || pack.Side.IsEast() {
			// Similar logic to the above, but widget is packed against the
			// left or right edge. Handle vertical Fill to grow the widget.
			if pack.FillY && resize.H < innerFrameSize.H {
				resize.H = innerFrameSize.H - w.BoxThickness(2) // TODO: child.BoxThickness instead??
				resized = true
			}

			// Vertically align the widgets.
			if resize.H < innerFrameSize.H {
				if pack.Side.IsMiddle() {
					point.Y = (innerFrameSize.H / 2) - (resize.H / 2) // - w.BoxThickness(1)
				} else if pack.Side.IsNorth() {
					point.Y = pack.PadY // - w.BoxThickness(4)
				} else if pack.Side.IsSouth() {
					point.Y = innerFrameSize.H - resize.H - pack.PadY
				}
				moved = true
			}
		} else {
			panic("unsupported pack.Side")
		}

		if resized && size != resize {
			child.ResizeAuto(resize)
			child.Compute(e)
		}
		if moved {
			child.MoveTo(point)
		}
	}

	// TODO: the Frame should ResizeAuto so it doesn't mark fixedSize=true.
	// Currently there's a bug where frames will grow when the window grows but
	// never shrink again when the window shrinks.
	// if !w.FixedSize() {
	w.Resize(render.NewRect(
		frameSize.W-w.BoxThickness(2),
		frameSize.H-w.BoxThickness(2),
	))
	// }
}

// Side is a cardinal direction.
type Side uint8

// Side values.
const (
	Center Side = iota
	N
	NE
	E
	SE
	S
	SW
	W
	NW
)

// Range of Side values.
const (
	SideMin = Center
	SideMax = NW
)

// IsNorth returns if the side is N, NE or NW.
func (a Side) IsNorth() bool {
	return a == N || a == NE || a == NW
}

// IsSouth returns if the side is S, SE or SW.
func (a Side) IsSouth() bool {
	return a == S || a == SE || a == SW
}

// IsEast returns if the side is E, NE or SE.
func (a Side) IsEast() bool {
	return a == E || a == NE || a == SE
}

// IsWest returns if the side is W, NW or SW.
func (a Side) IsWest() bool {
	return a == W || a == NW || a == SW
}

// IsCenter returns if the side is Center, N or S, to determine
// whether to align text as centered for North/South sides.
func (a Side) IsCenter() bool {
	return a == Center || a == N || a == S
}

// IsMiddle returns if the side is Center, E or W, to determine
// whether to align text as middled for East/West sides.
func (a Side) IsMiddle() bool {
	return a == Center || a == W || a == E
}

type packLayout struct {
	widgets []packedWidget
}

type packedWidget struct {
	widget Widget
	pack   Pack
	fill   uint8
}

// packedWidget.fill values
const (
	fillNone uint8 = iota
	fillX
	fillY
	fillBoth
)
