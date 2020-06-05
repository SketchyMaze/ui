package main

import (
	"fmt"
	"os"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/event"
	"git.kirsle.net/go/render/sdl"
	"git.kirsle.net/go/ui"
)

// Program globals.
var (
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

func init() {
	sdl.DefaultFontFilename = "../DejaVuSans.ttf"
}

func main() {
	mw, err := ui.NewMainWindow("Hello World", Width, Height)
	if err != nil {
		panic(err)
	}

	// Menu bar.
	menu := ui.NewMenuBar("Main Menu")
	file := menu.AddMenu("Options")
	file.AddItem("New window", func() {
		addWindow(mw)
	})
	file.AddItem("Close all windows", func() {
		OpenWindows -= mw.Supervisor().CloseAllWindows()
	})

	menu.Supervise(mw.Supervisor())
	menu.Compute(mw.Engine)
	mw.Pack(menu, menu.PackTop())

	// Add some windows to play with.
	addWindow(mw)
	addWindow(mw)

	mw.SetBackground(render.White)

	mw.OnLoop(func(e *event.State) {
		if e.Escape {
			os.Exit(0)
		}
	})

	mw.MainLoop()
}

// Add a new child window.
func addWindow(mw *ui.MainWindow) {
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
	win1.Compute(mw.Engine)
	win1.Supervise(mw.Supervisor())

	// Re-open a window when the last one is closed.
	OpenWindows++
	win1.Handle(ui.CloseWindow, func(ed ui.EventData) error {
		OpenWindows--
		if OpenWindows <= 0 {
			addWindow(mw)
		}
		return nil
	})

	// Default placement via cascade.
	win1.MoveTo(Cascade)
	Cascade.Add(CascadeStep)
	if Cascade.Y > Height-240-64 {
		CascadeLoops++
		Cascade.Y = 32
		Cascade.X = 24 * CascadeLoops
	}

	// Add a window duplicator button.
	btn2 := ui.NewButton(title+":Button2", ui.NewLabel(ui.Label{
		Text: "New Window",
	}))
	btn2.Handle(ui.Click, func(ed ui.EventData) error {
		addWindow(mw)
		return nil
	})
	mw.Add(btn2)
	win1.Place(btn2, ui.Place{
		Top:   10,
		Right: 10,
	})
}
