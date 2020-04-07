package main

import (
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
	Cascade     = render.NewPoint(10, 10)
	CascadeStep = render.NewPoint(24, 24)
)

func init() {
	sdl.DefaultFontFilename = "../DejaVuSans.ttf"
}

func main() {
	mw, err := ui.NewMainWindow("Hello World", Width, Height)
	if err != nil {
		panic(err)
	}

	// Add some windows to play with.
	addWindow(mw, "First window")
	addWindow(mw, "Second window")

	mw.SetBackground(render.White)

	mw.OnLoop(func(e *event.State) {
		if e.Escape {
			os.Exit(0)
		}
	})

	mw.MainLoop()
}

// Add a new child window.
func addWindow(mw *ui.MainWindow, title string) {
	win1 := ui.NewWindow(title)
	win1.Configure(ui.Config{
		Width:  640,
		Height: 480,
	})
	win1.Compute(mw.Engine)
	win1.Supervise(mw.Supervisor())

	// Attach it to the MainWindow with no placement management, i.e.
	// instead of Pack() or Place(). Since draggable windows set their own
	// position, a position manager would only interfere and "snap" the
	// window back into place as soon as you drop the title bar!
	// mw.Attach(win1)

	// Default placement via cascade.
	win1.MoveTo(Cascade)
	Cascade.Add(CascadeStep)

	// Add a button to the window.
	// btn := ui.NewButton("Button1", ui.NewLabel(ui.Label{
	// 	Text: "Click me!",
	// }))
	// btn.Handle(ui.Click, func(ed ui.EventData) {
	// 	fmt.Printf("Window '%s' button clicked!\n", title)
	// })
	// mw.Add(btn)
	// win1.Place(btn, ui.Place{
	// 	Top:  10,
	// 	Left: 10,
	// })

	// Add a window duplicator button.
	btn2 := ui.NewButton(title+":Button2", ui.NewLabel(ui.Label{
		Text: "New Window",
	}))
	btn2.Handle(ui.Click, func(ed ui.EventData) error {
		addWindow(mw, "New Window")
		return nil
	})
	mw.Add(btn2)
	win1.Place(btn2, ui.Place{
		Top:   10,
		Right: 10,
	})
}
