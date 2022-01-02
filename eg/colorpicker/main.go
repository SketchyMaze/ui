package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/sdl"
	"git.kirsle.net/go/ui"
)

var WindowColor = render.SkyBlue

func init() {
	sdl.DefaultFontFilename = "../DejaVuSans.ttf"
}

func main() {
	mw, err := ui.NewMainWindow("ColorPicker Demo", 812, 375)
	if err != nil {
		panic(err)
	}

	mw.SetBackground(WindowColor)

	btn := ui.NewButton("Test", ui.NewLabel(ui.Label{
		Text: "Pick the background color",
		Font: render.Text{
			Size: 32,
		},
	}))
	btn.Handle(ui.Click, func(ed ui.EventData) error {
		colorpicker, err := ui.NewColorPicker(ui.ColorPicker{
			Title:      "Select a background color",
			Supervisor: mw.Supervisor(),
			Engine:     mw.Engine,
			Color:      WindowColor,

			// Until the UI toolkit has normal text entry controls, this work-around
			// allows your application to ask the user to enter a hex color code
			// themselves, using any means available. This is an asynchronous
			// procedure where you are given a callback function to send your answer
			// whenever you have it. For this example, look for the prompt
			// question in your terminal window!
			OnManualInput: func(callback func(render.Color)) {
				// Prompt the user to enter a hex color in the terminal.
				var s string
				fmt.Fprintf(os.Stderr, "Enter a hexadecimal color code> ")
				r := bufio.NewReader(os.Stdin)
				for {
					s, _ = r.ReadString('\n')
					if s != "" {
						break
					}
				}

				// Parse it as a color.
				fmt.Printf("Answer: %s\n", s)
				color, err := render.HexColor(strings.TrimSpace(s))
				if err != nil {
					fmt.Printf("%s\n", err)
					return
				}

				// Ping the callback function with our answer.
				callback(color)
			},
		})
		if err != nil {
			fmt.Printf("Error initializing colorpicker: %s\n", err)
			return err
		}

		colorpicker.Center(mw.Engine.WindowSize())
		colorpicker.Then(func(color render.Color) {
			WindowColor = color
			mw.SetBackground(WindowColor)
		})
		colorpicker.OnCancel(func() {
			fmt.Println("ColorPicker was dismissed by user")
		})

		fmt.Printf("Open ColorPicker: %+v\n", colorpicker)
		colorpicker.Show()
		return nil
	})
	mw.Place(btn, ui.Place{
		Center: true,
		Middle: true,
	})

	mw.MainLoop()
}
