package ui

import (
	"errors"
	"sync"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/event"
)

// Event is a named event that the supervisor will send.
type Event int

// Events.
const (
	NullEvent Event = iota
	MouseOver
	MouseOut
	MouseDown
	MouseUp
	Click
	KeyDown
	KeyUp
	KeyPress

	// Drag/drop event handlers.
	DragStop // if a widget is being dragged and the drag is done
	DragMove // mouse movements sent to a widget being dragged.
	Drop     // a "drop site" widget under the cursor when a drag is done

	// Window Manager events.
	CloseWindow
	MaximizeWindow
	MinimizeWindow
	CloseModal

	// Lifecycle event handlers.
	Compute // fired whenever the widget runs Compute
	Present // fired whenever the widget runs Present
)

// EventData carries common data to event handlers.
type EventData struct {
	// Point is usually the cursor position on click and mouse events.
	Point render.Point

	// Engine is the render engine on Compute and Present events.
	Engine render.Engine

	// Supervisor is the reference to the supervisor who sent the event.
	Supervisor *Supervisor
}

// Supervisor keeps track of widgets of interest to notify them about
// interaction events such as mouse hovers and clicks in their general
// vicinity.
type Supervisor struct {
	lock     sync.RWMutex
	serial   int                 // ID number of each widget added in order
	widgets  map[int]WidgetSlot  // map of widget ID to WidgetSlot
	hovering map[int]interface{} // map of widgets under the cursor
	clicked  map[int]bool        // map of widgets being clicked
	dd       *DragDrop

	// Stack of modal widgets that have event priority.
	modals []Widget

	// List of window focus history for Window Manager.
	winFocus  *FocusedWindow
	winTop    *FocusedWindow // pointer to top-most window
	winBottom *FocusedWindow // pointer to bottom-most window
}

// WidgetSlot holds a widget with a unique ID number in a sorted list.
type WidgetSlot struct {
	id     int
	widget Widget
}

// NewSupervisor creates a supervisor.
func NewSupervisor() *Supervisor {
	return &Supervisor{
		widgets:  map[int]WidgetSlot{},
		hovering: map[int]interface{}{},
		clicked:  map[int]bool{},
		modals:   []Widget{},
		dd:       NewDragDrop(),
	}
}

// DragStart sets the drag state without a widget.
//
// An example where you'd use this is if you want a widget to respond to a
// Drop event (mouse released over a drop-site widget) but the 'thing' being
// dragged is not a ui.Widget, i.e., for custom app specific logic.
func (s *Supervisor) DragStart() {
	s.dd.Start()
}

// DragStartWidget sets the drag state to true with a target widget attached.
//
// The widget being dragged is given DragMove events while the drag is
// underway. When the mouse button is released, the widget is given a
// DragStop event and the widget below the cursor is given a Drop event.
func (s *Supervisor) DragStartWidget(w Widget) {
	s.dd.SetWidget(w)
	s.dd.Start()
}

// DragStop stops the drag state.
func (s *Supervisor) DragStop() {
	s.dd.Stop()
}

// IsDragging returns whether the drag state is enabled.
func (s *Supervisor) IsDragging() bool {
	return s.dd.IsDragging()
}

// Error messages that may be returned by Supervisor.Loop()
var (
	// The caller should STOP forwarding any mouse or keyboard events to any
	// other handles for the remainder of this tick.
	ErrStopPropagation = errors.New("stop all event propagation")
	ErrNoEventHandler  = errors.New("no event handler")
)

// Loop to check events and pass them to managed widgets.
//
// Useful errors returned by this may be:
// - ErrStopPropagation
func (s *Supervisor) Loop(ev *event.State) error {
	var (
		XY = render.Point{
			X: ev.CursorX,
			Y: ev.CursorY,
		}
	)

	// See if we are hovering over any widgets.
	hovering, outside := s.Hovering(XY)

	// If we are dragging something around, do not trigger any mouse events
	// to other widgets but DO notify any widget we dropped on top of!
	if s.dd.IsDragging() {
		if !ev.Button1 && !ev.Button3 {
			// The mouse has been released. TODO: make mouse button important?
			for _, child := range hovering {
				child.widget.Event(Drop, EventData{
					Point: XY,
				})
			}
			s.DragStop()
		} else {
			// If we have a target widget being dragged, send it mouse events.
			if target := s.dd.Widget(); target != nil {
				target.Event(DragMove, EventData{
					Point: XY,
				})
			}
		}
		return ErrStopPropagation
	}

	// Check if the top focused window has been closed and auto-focus the next.
	if s.winFocus != nil && s.winFocus.window.Hidden() {
		next := s.winFocus.next
		for next != nil {
			if !next.window.Hidden() {
				s.FocusWindow(next.window)
				break
			}
			next = next.next
		}
	}

	// Run events in managed windows first, from top to bottom.
	// Widgets in unmanaged windows will be handled next.
	// err := s.runWindowEvents(XY, ev, hovering, outside)
	// Only run if there is no active modal (modals have top priority)
	if len(s.modals) == 0 {
		handled, err := s.runWidgetEvents(XY, ev, hovering, outside, true)
		if err == ErrStopPropagation || handled {
			// A widget in the active window has accepted an event. Do not pass
			// the event also to lower widgets.
			return ErrStopPropagation
		}
	}

	// Run events for the other widgets not in a managed window.
	// (Modal event priority is handled in runWidgetEvents)
	s.runWidgetEvents(XY, ev, hovering, outside, false)

	return nil
}

// Hovering returns all of the widgets managed by Supervisor that are under
// the mouse cursor. Returns the set of widgets below the cursor and the set
// of widgets not below the cursor.
func (s *Supervisor) Hovering(cursor render.Point) (hovering, outside []WidgetSlot) {
	var XY = cursor // for shorthand
	hovering = []WidgetSlot{}
	outside = []WidgetSlot{}

	// Check all the widgets under our care.
	for child := range s.Widgets() {
		var (
			w  = child.widget
			P  = AbsolutePosition(w)
			S  = w.Size()
			P2 = render.Point{
				X: P.X + S.W,
				Y: P.Y + S.H,
			}
		)

		if XY.X >= P.X && XY.X < P2.X && XY.Y >= P.Y && XY.Y < P2.Y {
			// Cursor intersects the widget.
			hovering = append(hovering, child)
		} else {
			outside = append(outside, child)
		}
	}

	return hovering, outside
}

// runWindowEvents is a subroutine of Supervisor.Loop().
//
// After determining the widgets below the cursor (hovering) and outside the
// cursor, transmit mouse events to the widgets.
//
// This function has two use cases:
// - In runWindowEvents where we run events for the top-most focused window of
//   the window manager.
// - In Supervisor.Loop() for the widgets that are NOT owned by a managed
//   window, so that these widgets always get events.
//
// Parameters:
//    XY (Point): mouse cursor position as calculated in Loop()
//    ev, hovering, outside: values from Loop(), self explanatory.
//    behavior: indicates how this method is being used.
//
// behavior options:
//    0: widgets NOT part of a managed window. On this pass, if a widget IS
//       a part of a window, it gets no events triggered.
//    1: widgets are part of the active focused window.
func (s *Supervisor) runWidgetEvents(XY render.Point, ev *event.State,
	hovering, outside []WidgetSlot, toFocusedWindow bool) (bool, error) {
	// Do we run any events?
	var (
		stopPropagation bool
		ranEvents       bool
	)

	// Do we have active modals? Modal widgets have top event priority given
	// only to the top-most modal.
	var modal Widget
	if len(s.modals) > 0 {
		modal = s.modals[len(s.modals)-1]
	}

	// If we're running this method in "Phase 2" (to widgets NOT in the focused
	// window), only send mouse events to widgets if the cursor is NOT inside
	// the bounding box of the active focused window. Prevents clicking "thru"
	// the window and activating widgets/other windows behind it.
	var cursorInsideFocusedWindow bool
	if !toFocusedWindow && s.winFocus != nil && !s.winFocus.window.Hidden() {
		// Get the bounding box of the focused window.
		if XY.Inside(AbsoluteRect(s.winFocus.window)) {
			cursorInsideFocusedWindow = true
		}
	}

	// Handler for an Event response errors.
	handle := func(err error) {
		// Did any event handler run?
		if err != ErrNoEventHandler {
			ranEvents = true
		}

		// Are we stopping propagation?
		if err == ErrStopPropagation {
			stopPropagation = true
		}
	}

	for _, child := range hovering {
		if stopPropagation {
			break
		}

		// If the cursor is inside the box of the focused window, don't trigger
		// active (hovering) mouse events. MouseOut type events, below, can still
		// trigger.
		// Does not apply when a modal widget is active.
		if cursorInsideFocusedWindow && modal == nil {
			break
		}

		var (
			id = child.id
			w  = child.widget
		)
		if w.Hidden() {
			// TODO: somehow the Supervisor wasn't triggering hidden widgets
			// anyway, but I don't know why. Adding this check for safety.
			continue
		}

		// If we have a modal active, validate this widget is a child of
		// the modal widget.
		if modal != nil {
			if !HasParent(w, modal) {
				continue
			}
		}

		// Check if the widget is part of a Window managed by Supervisor.
		isManaged, isFocused := widgetInFocusedWindow(w)

		// Are we sending events to it?
		if toFocusedWindow {
			// Only sending events to widgets owned by the focused window.
			if !(isManaged && isFocused) {
				continue
			}
		} else {
			// Sending only to widgets NOT managed by a window. This can include
			// Window widgets themselves, so lower unfocused windows may be
			// brought to foreground.
			window, isWindow := w.(*Window)
			if isManaged && !isWindow {
				continue
			}

			// It is a window, but can only be the non-focused window.
			if isWindow && window.focused {
				continue
			}
		}

		// Cursor has intersected the widget.
		if _, ok := s.hovering[id]; !ok {
			handle(w.Event(MouseOver, EventData{
				Point: XY,
			}))
			s.hovering[id] = nil
		}

		isClicked, _ := s.clicked[id]
		if ev.Button1 {
			if !isClicked {
				err := w.Event(MouseDown, EventData{
					Point: XY,
				})
				handle(err)
				s.clicked[id] = true
			}
		} else if isClicked {
			handle(w.Event(MouseUp, EventData{
				Point: XY,
			}))
			handle(w.Event(Click, EventData{
				Point: XY,
			}))
			delete(s.clicked, id)
		}
	}
	for _, child := range outside {
		var (
			id = child.id
			w  = child.widget
		)

		// If we have a modal active, validate this widget is a child of
		// the modal widget.
		if modal != nil {
			if !HasParent(w, modal) {
				continue
			}
		}

		// Cursor is not intersecting the widget.
		if _, ok := s.hovering[id]; ok {
			handle(w.Event(MouseOut, EventData{
				Point: XY,
			}))
			delete(s.hovering, id)
		}

		if _, ok := s.clicked[id]; ok {
			handle(w.Event(MouseUp, EventData{
				Point: XY,
			}))
			delete(s.clicked, id)
		}
	}

	// If a modal is active and a click was registered outside the modal's
	// bounding box, send the CloseModal event.
	if modal != nil && !XY.Inside(AbsoluteRect(modal)) {
		if ev.Button1 {
			modal.Event(CloseModal, EventData{
				Supervisor: s,
			})
		}
	}

	// If there was a modal, return stopPropagation (so callers that manage
	// events externally of go/ui can see that a modal intercepted events)
	if modal != nil {
		return ranEvents, ErrStopPropagation
	}

	// If a stopPropagation was called, return it up the stack.
	if stopPropagation {
		return ranEvents, ErrStopPropagation
	}

	// If ANY event handler was called, return nil to signal
	return ranEvents, nil
}

// Widgets returns a channel of widgets managed by the supervisor in the order
// they were added.
func (s *Supervisor) Widgets() <-chan WidgetSlot {
	pipe := make(chan WidgetSlot)
	go func() {
		for i := 0; i < s.serial; i++ {
			if w, ok := s.widgets[i]; ok {
				pipe <- w
			}
		}
		close(pipe)
	}()
	return pipe
}

// Present all widgets managed by the supervisor.
//
// NOTE: only the Window Manager feature uses this method, and this method
// will render the windows from bottom to top with the focused window on top.
// For other widgets, they should be added to a parent Frame that will call
// Present on them each time the parent Presents, or otherwise you need to
// manage the presentation of widgets outside the Supervisor.
func (s *Supervisor) Present(e render.Engine) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Render the window manager windows from bottom to top.
	s.presentWindows(e)

	// Render the modals from bottom to top.
	if len(s.modals) > 0 {
		for _, modal := range s.modals {
			modal.Present(e, modal.Point())
		}
	}
}

// Add a widget to be supervised. Has no effect if the widget is already
// under the supervisor's care.
func (s *Supervisor) Add(w Widget) {
	s.lock.Lock()

	// Check it's not already there.
	for _, child := range s.widgets {
		if child.widget == w {
			return
		}
	}

	// Add it.
	s.widgets[s.serial] = WidgetSlot{
		id:     s.serial,
		widget: w,
	}
	s.serial++
	s.lock.Unlock()
}

// PushModal sets the widget to be a "modal" for the Supervisor.
//
// Modal widgets have top-most event priority: mouse and click events go ONLY
// to the modal and its descendants. Modals work as a stack: the most recently
// pushed widget is the active modal, and popping the modal will make the
// next most-recent widget be the active modal.
//
// If a Click event registers OUTSIDE the bounds of the modal widget, the
// widget receives a CloseModal event.
//
// Returns the length of the modal stack.
func (s *Supervisor) PushModal(w Widget) int {
	s.modals = append(s.modals, w)
	return len(s.modals)
}

// PopModal attempts to pop the modal from the stack, but only if the modal
// is at the top of the stack.
//
// A widget may safely attempt to PopModal itself on a CloseModal event to
// close themselves when the user clicks outside their box. If there were a
// newer modal on the stack, this PopModal action would do nothing.
func (s *Supervisor) PopModal(w Widget) bool {
	// only can pop if the topmost widget is the one being asked for
	if len(s.modals) > 0 && s.modals[len(s.modals)-1] == w {
		modal := s.modals[len(s.modals)-1]
		modal.Hide()

		// pop it off
		s.modals = s.modals[:len(s.modals)-1]

		return true
	}

	return false
}

// GetModal returns the modal on the top of the stack, or nil if there is
// no modal on top.
func (s *Supervisor) GetModal() Widget {
	if len(s.modals) == 0 {
		return nil
	}
	return s.modals[len(s.modals)-1]
}
