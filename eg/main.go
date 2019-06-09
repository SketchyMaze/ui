package main

import (
	"git.kirsle.net/apps/doodle/lib/render"
	"git.kirsle.net/apps/doodle/lib/ui"
)

func main() {
	mw, err := ui.NewMainWindow("UI Toolkit Demo")
	if err != nil {
		panic(err)
	}

	leftFrame := ui.NewFrame("Left Frame")
	leftFrame.Configure(ui.Config{
		Width:       200,
		BorderSize:  1,
		BorderStyle: ui.BorderRaised,
		Background:  render.Grey,
	})
	mw.Pack(leftFrame, ui.Pack{
		Anchor: ui.W,
		FillY:  true,
	})

	mainFrame := ui.NewFrame("Main Frame")
	mainFrame.Configure(ui.Config{
		Background: render.RGBA(255, 255, 255, 180),
	})
	mw.Pack(mainFrame, ui.Pack{
		Anchor: ui.W,
		Expand: true,
		PadX:   10,
	})

	label := ui.NewLabel(ui.Label{
		Text: "Hello world",
	})
	leftFrame.Pack(label, ui.Pack{
		Anchor: ui.SE,
	})

	err = mw.MainLoop()
	if err != nil {
		panic("MainLoop:" + err.Error())
	}
}
