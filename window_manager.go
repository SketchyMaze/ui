package ui

import (
	"errors"
	"fmt"

	"git.kirsle.net/go/render"
)

/*
window_manager.go holds data types and Supervisor methods related to the
management of ui.Window widgets.
*/

// Window button options. OR these together in a call to Window.SetButtons().
const (
	CloseButton = 0x01

	// NOTICE: MaximizeButton behavior is currently buggy, window doesn't
	// redraw itself at the new size properly.
	MaximizeButton = 0x02

	// Minimize button has no default behavior attached; you can bind it with
	// window.Handle(MinimizeWindow) to set your own event handler.
	MinimizeButton = 0x04
)

// FocusedWindow is a doubly-linked list of recently focused Windows, with
// the current and most-recently focused on top. TODO make not exported.
type FocusedWindow struct {
	window *Window
	prev   *FocusedWindow
	next   *FocusedWindow
}

// String of the FocusedWindow returns the underlying Window's String().
func (fw FocusedWindow) String() string {
	return fw.window.String()
}

// Print the structure of the linked list from top to bottom.
func (fw *FocusedWindow) Print() {
	var (
		node = fw
		i    = 0
	)
	for node != nil {
		fmt.Printf("[%d] window=%s  prev=%s  next=%s\n",
			i, node.window, node.prev, node.next,
		)
		node = node.next
		i++
	}
}

// addWindow installs a Window into the supervisor to be managed. It is called
// by ui.Window.Supervise() and the newly added window becomes the focused
// one by default at the top of the linked list.
func (s *Supervisor) addWindow(win *Window) {
	// Record in the window that it is managed by Supervisor, useful to control
	// event propagation to non-focused windows.
	win.managed = true

	if s.winFocus == nil {
		// First window added.
		s.winFocus = &FocusedWindow{
			window: win,
		}
		s.winTop = s.winFocus
		s.winBottom = s.winFocus
		win.SetFocus(true)
	} else {
		// New window, make it the top one.
		oldTop := s.winFocus
		s.winFocus = &FocusedWindow{
			window: win,
			next:   oldTop,
		}
		oldTop.prev = s.winFocus
		oldTop.window.SetFocus(false)
		win.SetFocus(true)
	}
}

// FocusWindow brings the given window to the top of the supervisor's focus.
//
// The window must have previously been added to the supervisor's Window Manager
// by calling the Supervise() method of the window.
func (s *Supervisor) FocusWindow(win *Window) error {
	if s.winFocus == nil {
		return errors.New("no windows managed by supervisor")
	}

	// If the top window is already the target, return.
	if s.winFocus.window == win {
		return nil
	}

	// Find the window in the linked list.
	var (
		item      = s.winFocus   // item as we iterate the list
		oldTop    = s.winFocus   // original first item in the list
		target    *FocusedWindow // identified target window to raise
		newBottom *FocusedWindow // if the target was the bottom, this is new bottom
		i         = 0
	)
	for item != nil {
		if item.window == win {
			// Found it!
			target = item

			// Is it the last window in the list? Record the new bottom node.
			if item.next == nil && item.prev != nil {
				newBottom = item.prev
			}

			// Remove it from its position in the linked list. Join its
			// previous and next nodes to bridge the gap.
			if item.next != nil {
				item.next.prev = item.prev
			}
			if item.prev != nil {
				item.prev.next = item.next
			}

			break
		}
		item = item.next
		i++
	}

	// Found it?
	if target != nil {
		// Put the target at the top of the list, pointing to the old top.
		target.next = oldTop
		target.prev = nil
		oldTop.prev = target
		s.winFocus = target

		// Fix the top and bottom pointers.
		s.winTop = s.winFocus
		if newBottom != nil {
			s.winBottom = newBottom
		}

		// Toggle the focus states.
		oldTop.window.SetFocus(false)
		target.window.SetFocus(true)
	}

	return nil
}

// IsPointInWindow returns whether the given Point overlaps with a window managed
// by the Supervisor.
func (s *Supervisor) IsPointInWindow(point render.Point) bool {
	node := s.winFocus
	for node != nil {
		if point.Inside(AbsoluteRect(node.window)) && !node.window.hidden {
			return true
		}
		node = node.next
	}
	return false
}

// CloseAllWindows closes all open windows being managed by supervisor.
// Returns the number of windows closed.
func (s *Supervisor) CloseAllWindows() int {
	var (
		node = s.winFocus
		i    = 0
	)
	for node != nil {
		i++
		node.window.Hide()
		node = node.next
	}
	return i
}

// presentWindows draws the windows from bottom to top.
func (s *Supervisor) presentWindows(e render.Engine) {
	item := s.winBottom
	for item != nil {
		item.window.Compute(e)
		item.window.Present(e, item.window.Point())
		item = item.prev
	}
}
