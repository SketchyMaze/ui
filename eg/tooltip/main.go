package main

import (
	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/sdl"
	"git.kirsle.net/go/ui"
)

func init() {
	sdl.DefaultFontFilename = "../DejaVuSans.ttf"
}

func main() {
	mw, err := ui.NewMainWindow("Tooltip Demo", 800, 600)
	if err != nil {
		panic(err)
	}

	mw.SetBackground(render.White)
	CreateButtons(mw, mw.Frame())

	btn := ui.NewButton("Test", ui.NewLabel(ui.Label{
		Text: "Click me",
		Font: render.Text{
			Size: 32,
		},
	}))
	mw.Place(btn, ui.Place{
		Center: true,
		Middle: true,
	})

	ui.NewTooltip(btn, ui.Tooltip{
		Text: "Hello world\nGoodbye mars!\nBlah blah blah...\nLOL",
		Edge: ui.Right,
	})

	mw.MainLoop()
}

// CreateButtons creates a set of Placed buttons around all the edges and
// center of the parent frame.
func CreateButtons(window *ui.MainWindow, parent *ui.Frame) {
	// Draw buttons around the edges of the window.
	buttons := []struct {
		Label string
		Edge  ui.Edge
		Place ui.Place
	}{
		{
			Label: "Top Left",
			Edge:  ui.Right,
			Place: ui.Place{
				Point: render.NewPoint(12, 12),
			},
		},
		{
			Label: "Top Middle",
			Edge:  ui.Bottom,
			Place: ui.Place{
				Top:    12,
				Center: true,
			},
		},
		{
			Label: "Top Right",
			Edge:  ui.Left,
			Place: ui.Place{
				Top:   12,
				Right: 12,
			},
		},
		{
			Label: "Left Middle",
			Edge:  ui.Right,
			Place: ui.Place{
				Left:   12,
				Middle: true,
			},
		},
		{
			Label: "Center",
			Edge:  ui.Bottom,
			Place: ui.Place{
				Center: true,
				Middle: true,
			},
		},
		{
			Label: "Right Middle",
			Edge:  ui.Left,
			Place: ui.Place{
				Right:  12,
				Middle: true,
			},
		},
		{
			Label: "Bottom Left",
			Edge:  ui.Right,
			Place: ui.Place{
				Left:   12,
				Bottom: 12,
			},
		},
		{
			Label: "Bottom Center",
			Edge:  ui.Top,
			Place: ui.Place{
				Bottom: 12,
				Center: true,
			},
		},
		{
			Label: "Bottom Right",
			Edge:  ui.Left,
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
		button.Handle(ui.Click, func(ed ui.EventData) error {
			window.SetTitle(parent.Name + ": " + setting.Label)
			return nil
		})

		// Tooltip for it.
		ui.NewTooltip(button, ui.Tooltip{
			Text: setting.Label + " Tooltip",
			Edge: setting.Edge,
		})

		parent.Place(button, setting.Place)
		window.Add(button)
	}
}
