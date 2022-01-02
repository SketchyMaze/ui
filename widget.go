package ui

import (
	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui/theme"
)

// BorderStyle options for widget.SetBorderStyle()
type BorderStyle string

// Styles for a widget border.
const (
	BorderNone   BorderStyle = ""
	BorderSolid  BorderStyle = "solid"
	BorderRaised             = "raised"
	BorderSunken             = "sunken"
)

// Widget is a user interface element.
type Widget interface {
	ID() string           // Get the widget's string ID.
	IDFunc(func() string) // Set a function that returns the widget's ID.
	String() string
	Point() render.Point
	MoveTo(render.Point)
	MoveBy(render.Point)
	Size() render.Rect    // Return the Width and Height of the widget.
	FixedSize() bool      // Return whether the size is fixed (true) or automatic (false)
	BoxSize() render.Rect // Return the full size including the border and outline.
	Resize(render.Rect)
	ResizeBy(render.Rect)
	ResizeAuto(render.Rect)
	Rect() render.Rect // Return the full absolute rect combining the Size() and Point()

	Handle(Event, func(EventData) error)
	Event(Event, EventData) error // called internally to trigger an event

	// Thickness of the padding + border + outline.
	BoxThickness(multiplier int) int
	DrawBox(render.Engine, render.Point)

	// Widget configuration getters.
	Margin() int                  // Margin away from other widgets
	SetMargin(int)                //
	Background() render.Color     // Background color
	SetBackground(render.Color)   //
	Foreground() render.Color     // Foreground color
	SetForeground(render.Color)   //
	BorderStyle() BorderStyle     // Border style: none, raised, sunken
	SetBorderStyle(BorderStyle)   //
	BorderColor() render.Color    // Border color (default is Background)
	SetBorderColor(render.Color)  //
	BorderSize() int              // Border size (default 0)
	SetBorderSize(int)            //
	OutlineColor() render.Color   // Outline color (default Invisible)
	SetOutlineColor(render.Color) //
	OutlineSize() int             // Outline size (default 0)
	SetOutlineSize(int)           //

	// Visibility
	Hide()
	Show()
	Hidden() bool

	// Container widgets like Frames can wire up associations between the
	// child widgets and the parent.
	Parent() (parent Widget, ok bool)
	SetParent(parent Widget) // for the container to assign itself the parent
	Children() []Widget      // for containers to return their children

	// Run any render computations; by the end the widget must know its
	// Width and Height. For example the Label widget will render itself onto
	// an SDL Surface and then it will know its bounding box, but not before.
	Compute(render.Engine)

	// Render the final widget onto the drawing engine.
	Present(render.Engine, render.Point)

	// Destroy: implement this if you have resources to free up on teardown.
	Destroy()
}

// Config holds common base widget configs for quick configuration.
type Config struct {
	// Size management. If you provide a non-zero value for Width and Height,
	// the widget will be resized and the "fixedSize" flag is set, meaning it
	// will not re-compute its size dynamically. To set the size while also
	// keeping the auto-resize property, pass AutoResize=true too. This is
	// mainly used internally when widgets are calculating their automatic sizes.
	AutoResize   bool
	Width        int
	Height       int
	Margin       int
	MarginX      int
	MarginY      int
	Background   render.Color
	Foreground   render.Color
	BorderSize   int
	BorderStyle  BorderStyle
	BorderColor  render.Color
	OutlineSize  int
	OutlineColor render.Color
}

// BaseWidget holds common functionality for all widgets, such as managing
// their widths and heights.
type BaseWidget struct {
	id           string
	idFunc       func() string
	fixedSize    bool
	hidden       bool
	width        int
	height       int
	point        render.Point
	margin       int
	background   render.Color
	foreground   render.Color
	borderStyle  BorderStyle
	borderColor  render.Color
	borderSize   int
	outlineColor render.Color
	outlineSize  int
	handlers     map[Event][]func(EventData) error
	hasParent    bool
	parent       Widget
}

// SetID sets a string name for your widget, helpful for debugging purposes.
func (w *BaseWidget) SetID(id string) {
	w.id = id
}

// ID returns the ID that the widget calls itself by.
func (w *BaseWidget) ID() string {
	if w.idFunc == nil {
		w.IDFunc(func() string {
			return "Widget<Untitled>"
		})
	}
	return w.idFunc()
}

// IDFunc sets an ID function.
func (w *BaseWidget) IDFunc(fn func() string) {
	w.idFunc = fn
}

func (w *BaseWidget) String() string {
	return w.ID()
}

// Configure the base widget with all the common properties at once. Any
// property left as the zero value will not update the widget.
func (w *BaseWidget) Configure(c Config) {
	if c.Width != 0 || c.Height != 0 {
		w.fixedSize = !c.AutoResize
		if c.Width != 0 {
			w.width = c.Width
		}
		if c.Height != 0 {
			w.height = c.Height
		}
	}

	if c.Margin != 0 {
		w.margin = c.Margin
	}
	if c.Background != render.Invisible {
		w.background = c.Background
	}
	if c.Foreground != render.Invisible {
		w.foreground = c.Foreground
	}
	if c.BorderColor != render.Invisible {
		w.borderColor = c.BorderColor
	}
	if c.OutlineColor != render.Invisible {
		w.outlineColor = c.OutlineColor
	}

	if c.BorderSize != 0 {
		w.borderSize = c.BorderSize
	}
	if c.BorderStyle != BorderNone {
		w.borderStyle = c.BorderStyle
	}
	if c.OutlineSize != 0 {
		w.outlineSize = c.OutlineSize
	}
}

// Rect returns the widget's absolute rectangle, the combined Size and Point.
func (w *BaseWidget) Rect() render.Rect {
	return render.Rect{
		X: w.point.X,
		Y: w.point.Y,
		W: w.width,
		H: w.height,
	}
}

// Point returns the X,Y position of the widget on the window.
func (w *BaseWidget) Point() render.Point {
	return w.point
}

// MoveTo updates the X,Y position to the new point.
func (w *BaseWidget) MoveTo(v render.Point) {
	w.point = v
}

// MoveBy adds the X,Y values to the widget's current position.
func (w *BaseWidget) MoveBy(v render.Point) {
	w.point.X += v.X
	w.point.Y += v.Y
}

// Size returns the box with W and H attributes containing the size of the
// widget. The X,Y attributes of the box are ignored and zero.
func (w *BaseWidget) Size() render.Rect {
	return render.Rect{
		W: w.width,
		H: w.height,
	}
}

// BoxSize returns the full rendered size of the widget including its box
// thickness (border, padding and outline).
func (w *BaseWidget) BoxSize() render.Rect {
	return render.Rect{
		W: w.width + w.BoxThickness(2),
		H: w.height + w.BoxThickness(2),
	}
}

// FixedSize returns whether the widget's size has been hard-coded by the user
// (true) or if it automatically resizes based on its contents (false).
func (w *BaseWidget) FixedSize() bool {
	return w.fixedSize
}

// Resize sets the size of the widget to the .W and .H attributes of a rect.
func (w *BaseWidget) Resize(v render.Rect) {
	w.fixedSize = true
	w.width = v.W
	w.height = v.H
}

// ResizeBy resizes by a relative amount.
func (w *BaseWidget) ResizeBy(v render.Rect) {
	w.fixedSize = true
	w.width += v.W
	w.height += v.H
}

// ResizeAuto sets the size of the widget but doesn't set the fixedSize flag.
func (w *BaseWidget) ResizeAuto(v render.Rect) {
	w.width = v.W
	w.height = v.H
}

// BoxThickness returns the full sum of the padding, border and outline.
// m = multiplier, i.e., 1 or 2. If m=1 this returns the box thickness of one
// edge of the widget, if m=2 it would account for both edges of the widget.
func (w *BaseWidget) BoxThickness(m int) int {
	if m == 0 {
		m = 1
	}
	return (w.Margin() * m) + (w.BorderSize() * m) + (w.OutlineSize() * m)
}

// Parent returns the parent widget, like a Frame, and a boolean indicating
// whether the widget had a parent.
func (w *BaseWidget) Parent() (Widget, bool) {
	return w.parent, w.hasParent
}

// SetParent sets the widget's parent. This function is called by container
// widgets like Frame when they add a child widget to their care.
// Pass a nil parent to unset the parent.
func (w *BaseWidget) SetParent(parent Widget) {
	if parent == nil {
		w.hasParent = false
		w.parent = nil
	} else {
		w.hasParent = true
		w.parent = parent
	}
}

// Children returns the widget's children, to be implemented by containers.
// The default implementation returns an empty slice.
func (w *BaseWidget) Children() []Widget {
	return []Widget{}
}

// Hide the widget from being rendered.
func (w *BaseWidget) Hide() {
	w.hidden = true
}

// Show the widget.
func (w *BaseWidget) Show() {
	w.hidden = false
}

// Hidden returns whether the widget is hidden. If this widget is not hidden,
// but it has a parent, this will recursively crawl the parents to see if any
// of them are hidden.
func (w *BaseWidget) Hidden() bool {
	if w.hidden {
		return true
	}

	// Return if any parents are hidden.
	parent, ok := w.Parent()
	for ok {
		if parent.Hidden() {
			return true
		}
		parent, ok = parent.Parent()
	}

	return false
}

// DrawBox draws the border and outline.
func (w *BaseWidget) DrawBox(e render.Engine, P render.Point) {
	var (
		S           = w.Size()
		outline     = w.OutlineSize()
		border      = w.BorderSize()
		borderColor = w.BorderColor()
		highlight   = borderColor.Lighten(theme.BorderColorOffset)
		shadow      = borderColor.Darken(theme.BorderColorOffset)
		color       render.Color
		box         = render.Rect{
			X: P.X,
			Y: P.Y,
			W: S.W,
			H: S.H,
		}
	)

	if borderColor == render.Invisible {
		borderColor = render.Red
	}

	// Draw the outline layer as the full size of the widget.
	if outline > 0 && w.OutlineColor() != render.Invisible {
		e.DrawBox(w.OutlineColor(), render.Rect{
			X: P.X,
			Y: P.Y,
			W: S.W,
			H: S.H,
		})
	}
	box.X += outline
	box.Y += outline
	box.W -= outline * 2
	box.H -= outline * 2

	// Highlight on the top left edge.
	if border > 0 {
		if w.BorderStyle() == BorderRaised {
			color = highlight
		} else if w.BorderStyle() == BorderSunken {
			color = shadow
		} else {
			color = borderColor
		}
		e.DrawBox(color, box)
	}

	// Shadow on the bottom right edge.
	box.X += border
	box.Y += border
	box.W -= border
	box.H -= border
	if w.BorderSize() > 0 {
		if w.BorderStyle() == BorderRaised {
			color = shadow
		} else if w.BorderStyle() == BorderSunken {
			color = highlight
		} else {
			color = borderColor
		}
		e.DrawBox(color, box)
	}

	// Background color of the button.
	box.W -= border
	box.H -= border
	if w.Background() != render.Invisible {
		e.DrawBox(w.Background(), box)
	}
}

// Margin returns the margin width.
func (w *BaseWidget) Margin() int {
	return w.margin
}

// SetMargin sets the margin width.
func (w *BaseWidget) SetMargin(v int) {
	w.margin = v
}

// Background returns the background color.
func (w *BaseWidget) Background() render.Color {
	return w.background
}

// SetBackground sets the color.
func (w *BaseWidget) SetBackground(c render.Color) {
	w.background = c
}

// Foreground returns the foreground color.
func (w *BaseWidget) Foreground() render.Color {
	return w.foreground
}

// SetForeground sets the color.
func (w *BaseWidget) SetForeground(c render.Color) {
	w.foreground = c
}

// BorderStyle returns the border style.
func (w *BaseWidget) BorderStyle() BorderStyle {
	return w.borderStyle
}

// SetBorderStyle sets the border style.
func (w *BaseWidget) SetBorderStyle(v BorderStyle) {
	w.borderStyle = v
}

// BorderColor returns the border color, or defaults to the background color.
func (w *BaseWidget) BorderColor() render.Color {
	if w.borderColor == render.Invisible {
		return w.Background()
	}
	return w.borderColor
}

// SetBorderColor sets the border color.
func (w *BaseWidget) SetBorderColor(c render.Color) {
	w.borderColor = c
}

// BorderSize returns the border thickness.
func (w *BaseWidget) BorderSize() int {
	return w.borderSize
}

// SetBorderSize sets the border thickness.
func (w *BaseWidget) SetBorderSize(v int) {
	w.borderSize = v
}

// OutlineColor returns the background color.
func (w *BaseWidget) OutlineColor() render.Color {
	return w.outlineColor
}

// SetOutlineColor sets the color.
func (w *BaseWidget) SetOutlineColor(c render.Color) {
	w.outlineColor = c
}

// OutlineSize returns the outline thickness.
func (w *BaseWidget) OutlineSize() int {
	return w.outlineSize
}

// SetOutlineSize sets the outline thickness.
func (w *BaseWidget) SetOutlineSize(v int) {
	w.outlineSize = v
}

// Compute calls the base widget's Compute function, which just triggers
// events on widgets that want to be notified when the widget computes.
func (w *BaseWidget) Compute(e render.Engine) {
	w.Event(Compute, EventData{
		Engine: e,
	})
}

// Present calls the base widget's Present function, which just triggers
// events on widgets that want to be notified when the widget presents.
func (w *BaseWidget) Present(e render.Engine, p render.Point) {
	w.Event(Present, EventData{
		Point:  p,
		Engine: e,
	})
}

// Event is called internally by Doodle to trigger an event.
// Handlers can return ErrStopPropagation to prevent further widgets being
// notified of events.
func (w *BaseWidget) Event(event Event, e EventData) error {
	if handlers, ok := w.handlers[event]; ok {
		for _, fn := range handlers {
			res := fn(e)
			if res == ErrStopPropagation {
				return res
			}
		}
	}
	return ErrNoEventHandler
}

// Handle an event in the widget.
func (w *BaseWidget) Handle(event Event, fn func(EventData) error) {
	if w.handlers == nil {
		w.handlers = map[Event][]func(EventData) error{}
	}

	if _, ok := w.handlers[event]; !ok {
		w.handlers[event] = []func(EventData) error{}
	}

	w.handlers[event] = append(w.handlers[event], fn)
}

// OnMouseOut should be overridden on widgets who want this event.
func (w *BaseWidget) OnMouseOut(render.Point) {}

// Destroy does nothing on the base widget. Implement it for widgets which need it.
func (w *BaseWidget) Destroy() {}
