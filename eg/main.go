package main

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui"
)

func main() {
	mw, err := ui.NewMainWindow("UI Toolkit Demo")
	if err != nil {
		panic(err)
	}

	leftFrame := ui.NewFrame("Left Frame")
	leftFrame.Configure(ui.Config{
		Width:       160,
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
		Anchor: ui.N,
		PadY:   12,
	})

	// Draw some buttons in the left frame.
	for i := 1; i <= 12; i++ {
		i := i

		btn := ui.NewButton(fmt.Sprintf("Button-%d", i), ui.NewLabel(ui.Label{
			Text: fmt.Sprintf("Button #%d", i),
		}))
		btn.Handle(ui.Click, func(p render.Point) {
			fmt.Printf("Button %d was clicked\n", i)
		})

		// Add the button to the MainWindow's event supervisor, so it may be
		// clicked and interacted with.
		mw.Add(btn)

		leftFrame.Pack(btn, ui.Pack{
			Anchor: ui.N,
			FillX:  true,
			PadY:   2,
		})
	}

	// Frame to show off check buttons.
	mainFrame.Pack(radioButtonFrame(mw), ui.Pack{
		Anchor: ui.N,
		FillX:  true,
		PadY:   8,
	})

	err = mw.MainLoop()
	if err != nil {
		panic("MainLoop:" + err.Error())
	}
}

// Frame that shows off radio buttons.
func radioButtonFrame(mw *ui.MainWindow) *ui.Frame {
	// The string variable that will be bound to the radio buttons.
	// This could also be a global variable at the package level.
	radioValue := "Red"

	// Main frame.
	frame := ui.NewFrame("radio button demo")
	frame.Configure(ui.Config{
		Background:  render.RGBA(153, 255, 153, 255),
		BorderSize:  1,
		BorderStyle: ui.BorderRaised,
	})

	// Top row to show the label and current radiobutton bound value.
	topFrame := ui.NewFrame("radio button label frame")
	frame.Pack(topFrame, ui.Pack{
		Anchor: ui.N,
		FillX:  true,
	})

	// Draw the labels.
	{
		label := ui.NewLabel(ui.Label{
			Text: "Radio buttons. Value:",
		})
		topFrame.Pack(label, ui.Pack{
			Anchor: ui.W,
		})

		valueLabel := ui.NewLabel(ui.Label{
			TextVariable: &radioValue,
		})
		topFrame.Pack(valueLabel, ui.Pack{
			Anchor: ui.W,
			PadX:   4,
		})
	}

	// The radio buttons themselves.
	btnFrame := ui.NewFrame("radio button frame")
	frame.Pack(btnFrame, ui.Pack{
		Anchor: ui.N,
		FillX:  true,
	})
	{
		colors := []string{"Red", "Green", "Blue", "Yellow"}
		for _, color := range colors {
			color := color

			btn := ui.NewRadioButton("color:"+color, &radioValue, color, ui.NewLabel(ui.Label{
				Text: color,
			}))
			mw.Add(btn)
			btnFrame.Pack(btn, ui.Pack{
				Anchor: ui.W,
				PadX:   2,
			})
		}
	}

	return frame
}
