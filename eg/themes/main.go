package main

import (
	"os"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/event"
	"git.kirsle.net/go/render/sdl"
	"git.kirsle.net/go/ui"
	"git.kirsle.net/go/ui/theme"
)

// Program globals.
var (
	// Size of the MainWindow.
	Width  = 1024
	Height = 768
)

func init() {
	sdl.DefaultFontFilename = "../DejaVuSans.ttf"
}

func main() {
	mw, err := ui.NewMainWindow("Theme Demo", Width, Height)
	if err != nil {
		panic(err)
	}

	// Menu bar.
	menu := ui.NewMenuBar("Main Menu")
	file := menu.AddMenu("Select Theme")
	file.AddItem("Default", func() {
		addWindow(mw, theme.Default)
	})
	file.AddItem("DefaultFlat", func() {
		addWindow(mw, theme.DefaultFlat)
	})
	file.AddItem("DefaultDark", func() {
		addWindow(mw, theme.DefaultDark)
	})

	menu.Supervise(mw.Supervisor())
	menu.Compute(mw.Engine)
	mw.Pack(menu, menu.PackTop())

	mw.SetBackground(render.White)

	mw.OnLoop(func(e *event.State) {
		if e.Escape {
			os.Exit(0)
		}
	})

	mw.MainLoop()
}

// Add a new child window.
func addWindow(mw *ui.MainWindow, theme theme.Theme) {
	ui.Theme = theme

	win1 := ui.NewWindow(theme.Name)
	win1.SetButtons(ui.CloseButton)
	win1.Configure(ui.Config{
		Width:  320,
		Height: 240,
	})
	win1.Compute(mw.Engine)
	win1.Supervise(mw.Supervisor())

	// Draw a label.
	label := ui.NewLabel(ui.Label{
		Text: theme.Name,
	})
	win1.Place(label, ui.Place{
		Top:  10,
		Left: 10,
	})

	// Add a button with tooltip.
	btn2 := ui.NewButton(theme.Name+":Button2", ui.NewLabel(ui.Label{
		Text: "Button",
	}))
	btn2.Handle(ui.Click, func(ed ui.EventData) error {
		return nil
	})
	mw.Add(btn2)
	win1.Place(btn2, ui.Place{
		Top:   10,
		Right: 10,
	})
	ui.NewTooltip(btn2, ui.Tooltip{
		Text: "Hello world!",
		Edge: ui.Bottom,
	})

	// Add a checkbox.
	var b bool
	cb := ui.NewCheckbox("Checkbox", &b, ui.NewLabel(ui.Label{
		Text: "Check me!",
	}))
	mw.Add(cb)
	win1.Place(cb, ui.Place{
		Top:  30,
		Left: 10,
	})
}
