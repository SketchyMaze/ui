// Example script for using the Place strategy of ui.Frame.
package main

import (
	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui"
)

func main() {
	mw, err := ui.NewMainWindow("Frame Placement Demo | Click a Button", 800, 600)
	if err != nil {
		panic(err)
	}

	mw.SetBackground(render.White)

	// Create a sub-frame with its own buttons packed within.
	frame := ui.NewFrame("Blue Frame")
	frame.Configure(ui.Config{
		Width:       300,
		Height:      150,
		Background:  render.DarkBlue,
		BorderSize:  1,
		BorderStyle: ui.BorderSunken,
	})
	mw.Place(frame, ui.Place{
		Point: render.NewPoint(80, 80),
	})

	// Create another frame that attaches itself to the bottom right
	// of the window.
	frame2 := ui.NewFrame("Red Frame")
	frame2.Configure(ui.Config{
		Width:      300,
		Height:     150,
		Background: render.DarkRed,
	})
	mw.Place(frame2, ui.Place{
		Right:  80,
		Bottom: 80,
	})

	// Draw rings of buttons around various widgets. The buttons say things
	// like "Top Left", "Top Center", "Left Middle", "Center" etc. encompassing
	// all 9 side placement options.
	CreateButtons(mw, frame)
	CreateButtons(mw, frame2)
	CreateButtons(mw, mw.Frame())

	mw.MainLoop()
}

// CreateButtons creates a set of Placed buttons around all the edges and
// center of the parent frame.
func CreateButtons(window *ui.MainWindow, parent *ui.Frame) {
	// Draw buttons around the edges of the window.
	buttons := []struct {
		Label string
		Place ui.Place
	}{
		{
			Label: "Top Left",
			Place: ui.Place{
				Point: render.NewPoint(12, 12),
			},
		},
		{
			Label: "Top Middle",
			Place: ui.Place{
				Top:    12,
				Center: true,
			},
		},
		{
			Label: "Top Right",
			Place: ui.Place{
				Top:   12,
				Right: 12,
			},
		},
		{
			Label: "Left Middle",
			Place: ui.Place{
				Left:   12,
				Middle: true,
			},
		},
		{
			Label: "Center",
			Place: ui.Place{
				Center: true,
				Middle: true,
			},
		},
		{
			Label: "Right Middle",
			Place: ui.Place{
				Right:  12,
				Middle: true,
			},
		},
		{
			Label: "Bottom Left",
			Place: ui.Place{
				Left:   12,
				Bottom: 12,
			},
		},
		{
			Label: "Bottom Center",
			Place: ui.Place{
				Bottom: 12,
				Center: true,
			},
		},
		{
			Label: "Bottom Right",
			Place: ui.Place{
				Bottom: 12,
				Right:  12,
			},
		},
	}
	for _, setting := range buttons {
		setting := setting

		button := ui.NewButton(setting.Label, ui.NewLabel(ui.Label{
			Text: setting.Label,
			Font: render.Text{
				FontFilename: "../DejaVuSans.ttf",
				Size:         12,
				Color:        render.Black,
			},
		}))

		// When clicked, change the window title to ID this button.
		button.Handle(ui.Click, func(p render.Point) {
			window.SetTitle(parent.Name + ": " + setting.Label)
		})

		parent.Place(button, setting.Place)
		window.Add(button)
	}
}
