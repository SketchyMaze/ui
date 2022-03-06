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
	tt := ui.NewTooltip(btn, ui.Tooltip{
		Text: "This is a tooltip that pops up\non mouse hover!",
		Edge: ui.Right,
	})

	// Notice: by default (with just the above code), the Tooltip will present
	// when its target widget presents. For densely packed UIs, the Tooltip may
	// be drawn "below" a neighboring widget, e.g. for horizontally packed buttons
	// where the Tooltip is on the Right: the tooltip for the left-most button
	// would present when the button does, but then the next button over will present
	// and overwrite the tooltip.
	//
	// For many simple UIs you can arrange your widgets and tooltip edge to
	// avoid this, but to guarantee the Tooltip always draws "on top", you
	// need to give it your Supervisor so it can register itself into its
	// Present stage (similar to window management). Be sure to call Supervisor.Present()
	// lastly in your main loop.
	tt.Supervise(mw.Supervisor())

	mw.MainLoop()
}
