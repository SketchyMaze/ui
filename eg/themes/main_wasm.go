// +build js,wasm

// WebAssembly version of the window manager demo.
// To build: make wasm
// To test: make wasm-serve

package main

import (
	"fmt"
	"time"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/canvas"
	"git.kirsle.net/go/ui"
)

// Program globals.
var (
	ThrottleFPS = 1000 / 60

	// Size of the MainWindow.
	Width  = 1024
	Height = 768

	// Cascade offset for creating multiple windows.
	Cascade      = render.NewPoint(10, 32)
	CascadeStep  = render.NewPoint(24, 24)
	CascadeLoops = 1

	// Colors for each window created.
	WindowColors = []render.Color{
		render.Blue,
		render.Red,
		render.DarkYellow,
		render.DarkGreen,
		render.DarkCyan,
		render.DarkBlue,
		render.DarkRed,
	}
	WindowID    int
	OpenWindows int
)

func main() {
	mw, err := canvas.New("canvas")
	if err != nil {
		panic(err)
	}

	// Bind DOM event handlers.
	mw.AddEventListeners()

	supervisor := ui.NewSupervisor()

	frame := ui.NewFrame("Main Frame")
	frame.Resize(render.NewRect(mw.WindowSize()))
	frame.Compute(mw)

	_, height := mw.WindowSize()
	lbl := ui.NewLabel(ui.Label{
		Text: "Window Manager Demo",
		Font: render.Text{
			FontFilename: "DejaVuSans.ttf",
			Size:         32,
			Color:        render.SkyBlue,
			Shadow:       render.SkyBlue.Darken(60),
		},
	})
	lbl.Compute(mw)
	lbl.MoveTo(render.NewPoint(
		20,
		height-lbl.Size().H-20,
	))

	// Menu bar.
	menu := ui.NewMenuBar("Main Menu")
	file := menu.AddMenu("Options")
	file.AddItem("New window", func() {
		addWindow(mw, frame, supervisor)
	})
	file.AddItem("Close all windows", func() {
		OpenWindows -= supervisor.CloseAllWindows()
	})

	menu.Supervise(supervisor)
	menu.Compute(mw)
	frame.Pack(menu, menu.PackTop())

	// Add some windows to play with.
	addWindow(mw, frame, supervisor)
	addWindow(mw, frame, supervisor)

	for {
		mw.Clear(render.RGBA(255, 255, 200, 255))
		start := time.Now()
		ev, err := mw.Poll()
		if err != nil {
			panic(err)
		}

		frame.Present(mw, frame.Point())
		lbl.Present(mw, lbl.Point())
		supervisor.Loop(ev)
		supervisor.Present(mw)

		var delay uint32
		elapsed := time.Now().Sub(start)
		tmp := elapsed / time.Millisecond
		if ThrottleFPS-int(tmp) > 0 {
			delay = uint32(ThrottleFPS - int(tmp))
		}
		mw.Delay(delay)
	}
}

// Add a new child window.
func addWindow(engine render.Engine, parent *ui.Frame, sup *ui.Supervisor) {
	var (
		color = WindowColors[WindowID%len(WindowColors)]
		title = fmt.Sprintf("Window %d", WindowID+1)
	)
	WindowID++

	win1 := ui.NewWindow(title)
	win1.SetButtons(ui.CloseButton)
	win1.ActiveTitleBackground = color
	win1.InactiveTitleBackground = color.Darken(60)
	win1.InactiveTitleForeground = render.Grey
	win1.Configure(ui.Config{
		Width:  320,
		Height: 240,
	})
	win1.Compute(engine)
	win1.Supervise(sup)

	// Re-open a window when the last one is closed.
	OpenWindows++
	win1.Handle(ui.CloseWindow, func(ed ui.EventData) error {
		OpenWindows--
		if OpenWindows <= 0 {
			addWindow(engine, parent, sup)
		}
		return nil
	})

	// Default placement via cascade.
	win1.MoveTo(Cascade)
	Cascade.Add(CascadeStep)
	if Cascade.Y > Height-240-64 {
		CascadeLoops++
		Cascade.Y = 24
		Cascade.X = 24 * CascadeLoops
	}

	// Add a window duplicator button.
	btn2 := ui.NewButton(title+":Button2", ui.NewLabel(ui.Label{
		Text: "New Window",
	}))
	btn2.Handle(ui.Click, func(ed ui.EventData) error {
		addWindow(engine, parent, sup)
		return nil
	})
	sup.Add(btn2)
	win1.Place(btn2, ui.Place{
		Top:   10,
		Right: 10,
	})
}
