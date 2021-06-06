package main

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/sdl"
	"git.kirsle.net/go/ui"
)

func init() {
	sdl.DefaultFontFilename = "../DejaVuSans.ttf"
}

func main() {
	mw, err := ui.NewMainWindow("Forms Test")
	if err != nil {
		panic(err)
	}

	mw.SetBackground(render.White)

	// Buttons row.
	{
		frame := ui.NewFrame("Frame 1")
		mw.Pack(frame, ui.Pack{
			Side: ui.N,
			FillX: true,
			Padding: 4,
		})
		label := ui.NewLabel(ui.Label{
			Text: "Buttons:",
		})
		frame.Pack(label, ui.Pack{
			Side: ui.W,
		})

		// Buttons.
		btn := ui.NewButton("Button 1", ui.NewLabel(ui.Label{
			Text: "Click me!",
		}))
		btn.Handle(ui.Click, func(ed ui.EventData) error {
			fmt.Println("Clicked!")
			return nil
		})
		frame.Pack(btn, ui.Pack{
			Side: ui.W,
			PadX: 4,
		})

		mw.Supervisor().Add(btn)
	}

	// Selectbox row.
	{
		frame := ui.NewFrame("Frame 2")
		mw.Pack(frame, ui.Pack{
			Side: ui.N,
			FillX: true,
			Padding: 4,
		})
		label := ui.NewLabel(ui.Label{
			Text: "Set window color:",
		})
		frame.Pack(label, ui.Pack{
			Side: ui.W,
		})

		var colors = []struct{
			Label string
			Value render.Color
		}{
			{"White", render.White},
			{"Yellow", render.Yellow},
			{"Cyan", render.Cyan},
			{"Green", render.Green},
			{"Blue", render.RGBA(0, 153, 255, 255)},
			{"Pink", render.Pink},
		}

		// Create the SelectBox and populate its options.
		sel := ui.NewSelectBox("Select 1", ui.Label{})
		for _, option := range colors {
			sel.AddItem(option.Label, option.Value, func() {
				fmt.Printf("Picked option: %s\n", option.Value)
			})
		}

		// On change: set the window BG color.
		sel.Handle(ui.Change, func(ed ui.EventData) error {
			if val, ok := sel.GetValue(); ok {
				if color, ok := val.Value.(render.Color); ok {
					fmt.Printf("Set background to: %s\n", val.Label)
					mw.SetBackground(color)
				} else {
					fmt.Println("Not a valid color!")
				}
			} else {
				fmt.Println("Not a valid SelectBox value!")
			}
			return nil
		})

		frame.Pack(sel, ui.Pack{
			Side: ui.W,
			PadX: 4,
		})
		sel.Supervise(mw.Supervisor())
		mw.Supervisor().Add(sel) // TODO: ideally Supervise() is all that's needed,
		                         // but w/o this extra Add() the Button doesn't react.
	}

	mw.MainLoop()
}
