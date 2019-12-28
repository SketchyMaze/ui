package main

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui"
)

func main() {
	mw, err := ui.NewMainWindow("Hello World")
	if err != nil {
		panic(err)
	}

	mw.SetBackground(render.White)

	// Draw a label.
	label := ui.NewLabel(ui.Label{
		Text: "Hello, world!",
		Font: render.Text{
			FontFilename: "../DejaVuSans.ttf",
			Size:         32,
			Color:        render.SkyBlue,
			Shadow:       render.SkyBlue.Darken(40),
		},
	})
	mw.Pack(label, ui.Pack{
		Anchor: ui.N,
		PadY:   12,
	})

	// Draw a button.
	button := ui.NewButton("My Button", ui.NewLabel(ui.Label{
		Text: "Click me!",
		Font: render.Text{
			FontFilename: "../DejaVuSans.ttf",
			Size:         12,
			Color:        render.Red,
			Padding:      4,
		},
	}))
	button.Handle(ui.Click, func(p render.Point) {
		fmt.Println("I've been clicked!")
	})
	mw.Pack(button, ui.Pack{
		Anchor: ui.N,
	})

	// Add the button to the MainWindow's Supervisor so it can be
	// clicked on and interacted with.
	mw.Add(button)

	mw.MainLoop()
}
