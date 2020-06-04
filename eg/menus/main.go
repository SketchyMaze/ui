package main

import (
	"fmt"
	"os"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/event"
	"git.kirsle.net/go/render/sdl"
	"git.kirsle.net/go/ui"
)

// Program globals.
var (
	// Size of the MainWindow.
	Width  = 640
	Height = 480

	BGColor = render.White
)

func init() {
	sdl.DefaultFontFilename = "../DejaVuSans.ttf"
}

func main() {
	mw, err := ui.NewMainWindow("Menu Demo", Width, Height)
	if err != nil {
		panic(err)
	}

	setupMainMenu(mw)

	// Menu button in middle of window.
	{
		btn := ui.NewMenuButton("MenuBtn", ui.NewLabel(ui.Label{
			Text: "Click me!",
		}))
		btn.Supervise(mw.Supervisor())
		mw.Place(btn, ui.Place{
			Center: true,
			Middle: true,
		})

		for _, label := range []string{
			"MenuButtons open menus",
			"when clicked. MenuBar is ",
			"a Frame of MenuButtons",
			"attached to the top of the",
			"window.",
			"",
			"They all provide a nice API",
			"to insert menus and items.",
		} {
			label := label
			if label == "" {
				btn.AddSeparator()
				continue
			}
			btn.AddItem(label, func() {
				fmt.Printf("Button '%s' clicked!\n", label)
			})
		}

	}

	// Menu button on the bottom right of screen.
	{
		btn := ui.NewMenuButton("BrBtn", ui.NewLabel(ui.Label{
			Text: "Fruits",
		}))
		btn.Supervise(mw.Supervisor())
		mw.Place(btn, ui.Place{
			Right:  20,
			Bottom: 20,
		})

		btn.AddItem("Apples", func() {})
		btn.AddItem("Oranges", func() {})
		btn.AddItem("Bananas", func() {})
		btn.AddItem("Pears", func() {})
	}

	// Menu button on the bottom left of screen.
	{
		btn := ui.NewMenuButton("BlBtn", ui.NewLabel(ui.Label{
			Text: "Set Window Color",
		}))
		btn.Supervise(mw.Supervisor())
		mw.Place(btn, ui.Place{
			Left:   20,
			Bottom: 20,
		})

		setBg := func(color render.Color) func() {
			return func() {
				BGColor = color
			}
		}

		// Really fancy buttons.
		var colors = []struct {
			label string
			hex   string
			color render.Color
		}{
			{"Black", "#000", render.Black},
			{"Red", "#F00", render.Red},
			{"Yellow", "#FF0", render.Yellow},
			{"Green", "#0F0", render.Green},
			{"Cyan", "#0FF", render.Cyan},
			{"Blue", "#00F", render.Blue},
			{"Magenta", "#F0F", render.Magenta},
			{"White", "#FFF", render.White},
		}
		for _, opt := range colors {
			item := btn.AddItemAccel(opt.label, opt.hex, setBg(opt.color))
			item.SetBackground(opt.color.Lighten(128))
		}

		// btn.AddItemAccel("Black", "#000", setBg(render.White))
		// btn.AddItemAccel("Red", "#F00", setBg(render.Red))
		// btn.AddItemAccel("Yellow", "#FF0", setBg(render.Yellow))
		// btn.AddItemAccel("Green", "#0F0", setBg(render.Green))
		// btn.AddItemAccel("Cyan", "#0FF", setBg(render.Cyan))
		// btn.AddItemAccel("Blue", "#00F", setBg(render.Blue))
		// btn.AddItemAccel("Magenta", "#F0F", setBg(render.Magenta))
		// btn.AddItemAccel("White", "#FFF", setBg(render.White))
	}

	// The "Long Menu" on the middle left side
	{
		btn := ui.NewMenuButton("BlBtn", ui.NewLabel(ui.Label{
			Text: "Tall Growing Menu",
		}))
		btn.Supervise(mw.Supervisor())
		mw.Place(btn, ui.Place{
			Left:   20,
			Middle: true,
		})

		var id int
		btn.AddItem("Add New Option", func() {
			id++
			id := id
			btn.AddItem(fmt.Sprintf("Menu Item #%d", id), func() {
				fmt.Printf("Chosen menu item %d\n", id)
			})
		})

		btn.AddSeparator()
	}

	mw.OnLoop(func(e *event.State) {
		mw.SetBackground(BGColor)
		if e.Up {
			fmt.Println("Supervised widgets:")
			for widg := range mw.Supervisor().Widgets() {
				fmt.Printf("%+v\n", widg)
			}
		}
		if e.Escape {
			os.Exit(0)
		}
	})

	mw.MainLoop()
}

func setupMainMenu(mw *ui.MainWindow) {
	bar := ui.NewMenuBar("Main Menu")

	fileMenu := bar.AddMenu("File")
	fileMenu.AddItemAccel("New", "Ctrl-N", func() {
		fmt.Println("Chose File->New")
	})
	fileMenu.AddItemAccel("Open", "Ctrl-O", func() {
		fmt.Println("Chose File->Open")
	})
	fileMenu.AddSeparator()
	fileMenu.AddItemAccel("Exit", "Alt-F4", func() {
		fmt.Println("Chose File->Exit")
		os.Exit(0)
	})

	editMenu := bar.AddMenu("Edit")
	editMenu.AddItemAccel("Undo", "Ctrl-Z", func() {})
	editMenu.AddItemAccel("Redo", "Shift-Ctrl-Z", func() {})
	editMenu.AddSeparator()
	editMenu.AddItemAccel("Cut", "Ctrl-X", func() {})
	editMenu.AddItemAccel("Copy", "Ctrl-C", func() {})
	editMenu.AddItemAccel("Paste", "Ctrl-V", func() {})
	editMenu.AddSeparator()
	editMenu.AddItem("Settings...", func() {})

	viewMenu := bar.AddMenu("View")
	viewMenu.AddItemAccel("Toggle Full Screen", "F11", func() {})

	helpMenu := bar.AddMenu("Help")
	helpMenu.AddItemAccel("Contents", "F1", func() {})
	helpMenu.AddItem("About", func() {})

	bar.Supervise(mw.Supervisor())
	bar.Compute(mw.Engine)
	mw.Pack(bar, bar.PackTop())

	fmt.Printf("Setup MenuBar: %s\n", bar.Size())
}
