package ui_test

import (
	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui"
)

// Example of using the Supervisor Window Manager.
func ExampleWindow() {
	mw, err := ui.NewMainWindow("Window Manager Example", 800, 600)
	if err != nil {
		panic(err)
	}

	// Create a window as normal.
	window := ui.NewWindow("Hello world!")
	window.Configure(ui.Config{
		Width:  320,
		Height: 240,
	})

	// Configure its title bar colors (optional; these are the defaults)
	window.ActiveTitleBackground = render.Blue
	window.ActiveTitleForeground = render.White
	window.InactiveTitleBackground = render.DarkGrey
	window.InactiveTitleForeground = render.Grey

	// Configure its window buttons (optional); default has no window buttons.
	// Window buttons are only functional in managed windows.
	window.SetButtons(ui.CloseButton | ui.MaximizeButton | ui.MinimizeButton)

	// Add some widgets to the window.
	btn := ui.NewButton("My Button", ui.NewLabel(ui.Label{
		Text: "Hello world!",
	}))
	window.Place(btn, ui.Place{
		Center: true,
		Middle: true,
	})

	// To enable the window manager controls, the key step is to give it
	// the Supervisor so it can be managed:
	window.Compute(mw.Engine)
	window.Supervise(mw.Supervisor())

	// Each loop you must then:
	// - Call Supervisor.Loop() as normal to handle events.
	// - Call Supervisor.Present() to draw the managed windows.
	// MainLoop() of the MainWindow does this for you.
	mw.MainLoop()
}
