// +build !js

package ui

import (
	"fmt"
	"time"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/event"
	"git.kirsle.net/go/render/sdl"
)

// Target frames per second for the MainWindow to render at.
var (
	FPS = 60
)

// Default width and height for MainWindow.
var (
	DefaultWidth  = 640
	DefaultHeight = 480
)

// MainWindow is the parent window of a UI application.
type MainWindow struct {
	Engine        render.Engine
	supervisor    *Supervisor
	frame         *Frame
	loopCallbacks []func(*event.State)
	w             int
	h             int
}

// NewMainWindow initializes the MainWindow. You should probably only have one
// of these per application. Dimensions are the width and height of the window.
//
// Example: NewMainWindow("Title Bar")  // default 640x480 window
// NewMainWindow("Title", 800, 600)     // both required
func NewMainWindow(title string, dimensions ...int) (*MainWindow, error) {
	var (
		width  = DefaultWidth
		height = DefaultHeight
	)

	if len(dimensions) > 0 {
		if len(dimensions) != 2 {
			return nil, fmt.Errorf("provide width and height dimensions, like NewMainWindow(title, 800, 600)")
		}
		width, height = dimensions[0], dimensions[1]
	}

	mw := &MainWindow{
		w:             width,
		h:             height,
		supervisor:    NewSupervisor(),
		loopCallbacks: []func(*event.State){},
	}

	mw.Engine = sdl.New(
		title,
		mw.w,
		mw.h,
	)
	if err := mw.Engine.Setup(); err != nil {
		return nil, err
	}

	// Add a default frame to the window.
	mw.frame = NewFrame("MainWindow Body")
	mw.frame.SetBackground(render.RGBA(0, 153, 255, 100))

	// Compute initial window size.
	mw.resized()

	return mw, nil
}

// SetTitle changes the title of the window.
func (mw *MainWindow) SetTitle(title string) {
	mw.Engine.SetTitle(title)
}

// Add a child widget to the window's supervisor. This alone does not make the
// child widget render each frame; use Pack, Place or Attach for that.
func (mw *MainWindow) Add(w Widget) {
	mw.supervisor.Add(w)
}

// Pack a child widget into the window's default frame.
func (mw *MainWindow) Pack(w Widget, pack Pack) {
	mw.Add(w)
	mw.frame.Pack(w, pack)
}

// Place a child widget into the window's default frame.
func (mw *MainWindow) Place(w Widget, config Place) {
	mw.Add(w)
	mw.frame.Place(w, config)
}

// Attach a child widget to the window without its position managed. The
// widget's Present() method will be called each time the window Presents, but
// the positioning of the child widget must be handled manually by the caller.
//
// Pack and Place are usually the methods you want to use to put a child widget
// into the window. One example use case for Attach is when you want to create
// child Window widgets which can be dragged by their title bars; their dynamic
// drag-drop positioning is best managed manually, and Pack or Place would
// interfere with their positioning otherwise.
//
// This also calls .Add() to add the widget to the MainWindow's Supervisor.
//
// Implementation details:
// - Adds the widget to the MainWindow's Supervisor.
// - Calls Frame.Add(w) so it will Present each time the main frame Presents.
// - Calls w.Compute() on your widget so it can calculate its initial size.
func (mw *MainWindow) Attach(w Widget) {
	mw.Add(w)
	mw.frame.Add(w)
	w.Compute(mw.Engine)
}

// Frame returns the window's main frame, if needed.
func (mw *MainWindow) Frame() *Frame {
	return mw.frame
}

// Supervisor returns the window's Supervisor instance.
func (mw *MainWindow) Supervisor() *Supervisor {
	return mw.supervisor
}

// resized handles the window being resized.
func (mw *MainWindow) resized() {
	mw.frame.Resize(render.Rect{
		W: mw.w,
		H: mw.h,
	})
}

// SetBackground changes the window's frame's background color.
func (mw *MainWindow) SetBackground(color render.Color) {
	mw.frame.SetBackground(color)
}

// OnLoop registers a function to be called on every loop of the main window.
// This enables your application to register global event handlers or whatnot.
// The function is called between the event polling and the updating of any UI
// elements.
func (mw *MainWindow) OnLoop(callback func(*event.State)) {
	mw.loopCallbacks = append(mw.loopCallbacks, callback)
}

// MainLoop starts the main event loop and blocks until there's an error.
func (mw *MainWindow) MainLoop() error {
	for true {
		if err := mw.Loop(); err != nil {
			return err
		}
	}
	return nil
}

// Loop does one loop of the UI.
func (mw *MainWindow) Loop() error {
	mw.Engine.Clear(render.White)

	// Record how long this loop took.
	start := time.Now()

	// Poll for events.
	ev, err := mw.Engine.Poll()
	if err != nil {
		return fmt.Errorf("event poll error: %s", err)
	}

	if ev.WindowResized {
		w, h := mw.Engine.WindowSize()
		if w != mw.w || h != mw.h {
			mw.w = w
			mw.h = h
			mw.resized()
		}
	}

	// Ping any loop callbacks.
	for _, cb := range mw.loopCallbacks {
		cb(ev)
	}

	mw.frame.Compute(mw.Engine)

	// Render the child widgets.
	mw.supervisor.Loop(ev)
	mw.frame.Present(mw.Engine, mw.frame.Point())
	mw.supervisor.Present(mw.Engine)
	mw.Engine.Present()

	// Delay to maintain target frames per second.
	var delay uint32
	var targetFPS = 1000 / FPS
	elapsed := time.Now().Sub(start) / time.Millisecond
	if targetFPS-int(elapsed) > 0 {
		delay = uint32(targetFPS - int(elapsed))
	}
	mw.Engine.Delay(delay)

	return nil
}
