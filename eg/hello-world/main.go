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
		Side: ui.N,
		PadY: 12,
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
	button.Handle(ui.Click, func(ed ui.EventData) {
		fmt.Println("I've been clicked!")
	})
	mw.Pack(button, ui.Pack{
		Side: ui.N,
	})

	mw.MainLoop()
}
