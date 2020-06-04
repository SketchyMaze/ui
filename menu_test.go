package ui_test

import (
	"git.kirsle.net/go/ui"
)

// Example of using the menu widgets.
func ExampleMenu() {
	mw, err := ui.NewMainWindow("Menu Bar Example", 800, 600)
	if err != nil {
		panic(err)
	}

	// Create a main menu for your window.
	menu := ui.NewMenuBar("Main Menu")

	// File menu. Some items with accelerators, some without.
	// NOTE: key bindings are up to you, the accelerators are
	// purely decorative.
	file := menu.AddMenu("File")
	file.AddItemAccel("New", "Ctrl-N", func() {})
	file.AddItemAccel("Open", "Ctrl-O", func() {})
	file.AddItemAccel("Save", "Ctrl-S", func() {})
	file.AddItem("Save as...", func() {})
	file.AddSeparator()
	file.AddItem("Close window", func() {})
	file.AddItemAccel("Exit", "Alt-F4", func() {})

	// Help menu.
	help := menu.AddMenu("Help")
	help.AddItemAccel("Contents", "F1", func() {})
	help.AddItem("About", func() {})

	// Give the menu bar your Supervisor so it can wire all
	// events up and make the menus work.
	menu.Supervise(mw.Supervisor())

	// Compute and pack the menu bar against the top of
	// the main window (or other parent container)
	menu.Compute(mw.Engine)
	mw.Pack(menu, menu.PackTop()) // Side: N, FillX: true

	// Each loop you must then:
	// - Call Supervisor.Loop() as normal to handle events.
	// - Call Supervisor.Present() to draw the modal popup menus.
	// MainLoop() of the MainWindow does this for you.
	mw.MainLoop()
}

// Example of using the MenuButton.
func ExampleMenuButton() {
	mw, err := ui.NewMainWindow("Menu Button", 800, 600)
	if err != nil {
		panic(err)
	}

	// Create a MenuButton much as you would a normal Button.
	btn := ui.NewMenuButton("Button1", ui.NewLabel(ui.Label{
		Text: "File",
	}))
	mw.Place(btn, ui.Place{ // place it in the center
		Center: true,
		Middle: true,
	})

	// Add menu items to it.
	btn.AddItemAccel("New", "Ctrl-N", func() {})
	btn.AddItemAccel("Open", "Ctrl-O", func() {})
	btn.AddItemAccel("Save", "Ctrl-S", func() {})
	btn.AddItem("Save as...", func() {})
	btn.AddSeparator()
	btn.AddItem("Close window", func() {})
	btn.AddItemAccel("Exit", "Alt-F4", func() {})

	// Add the button to Supervisor for events to work.
	btn.Supervise(mw.Supervisor())

	// Each loop you must then:
	// - Call Supervisor.Loop() as normal to handle events.
	// - Call Supervisor.Present() to draw the modal popup menus.
	// MainLoop() of the MainWindow does this for you.
	mw.MainLoop()
}
