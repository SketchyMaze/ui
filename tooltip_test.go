package ui_test

import "git.kirsle.net/go/ui"

// Tooltip usage example.
func ExampleTooltip() {
	mw, err := ui.NewMainWindow("Tooltip Example", 800, 600)
	if err != nil {
		panic(err)
	}

	// Add a widget that will have a tooltip attached, i.e. a button.
	btn := ui.NewButton("My Button", ui.NewLabel(ui.Label{
		Text: "Hello world!",
	}))
	mw.Place(btn, ui.Place{
		Center: true,
		Middle: true,
	})

	// Add a tooltip to it. The tooltip attaches itself to the button's
	// MouseOver, MouseOut, Compute and Present handlers -- you don't need to
	// place the tooltip inside the window or parent frame.
	ui.NewTooltip(btn, ui.Tooltip{
		Text: "This is a tooltip that pops up\non mouse hover!",
		Edge: ui.Right,
	})

	mw.MainLoop()
}
