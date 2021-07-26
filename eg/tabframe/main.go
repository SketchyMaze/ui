package main

// See the MakeTabFrame() function just below for the meat of this example.

import (
	"fmt"
	"os"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/event"
	"git.kirsle.net/go/render/sdl"
	"git.kirsle.net/go/ui"
	"git.kirsle.net/go/ui/theme"
)

// Program globals.
var (
	// Size of the MainWindow.
	Width  = 1024
	Height = 768

	// Cascade offset for creating multiple windows.
	Cascade      = render.NewPoint(10, 32)
	CascadeStep  = render.NewPoint(24, 24)
	CascadeLoops = 1

	// Colors for each window created.
	WindowColors = []render.Color{
		render.Blue,
		render.Red,
		render.DarkYellow,
		render.DarkGreen,
		render.DarkCyan,
		render.DarkBlue,
		render.DarkRed,
	}
	WindowID    int
	OpenWindows int

	TabFont = render.Text{
		Size:    10,
		Color:   render.Black,
		Padding: 4,
	}
)

func init() {
	sdl.DefaultFontFilename = "../DejaVuSans.ttf"
}

// MakeTabFrame is the example use of the TabFrame widget.
// The rest of this file is basically a copy of the eg/windows
// demo, except each window embeds the TabFrame.
func MakeTabFrame(mw *ui.MainWindow) *ui.TabFrame {
	notebook := ui.NewTabFrame("Example")

	// AddTab gives you the Frame to populate for that tab.

	// First Tab contents.
	tab1 := notebook.AddTab("Tab 1", ui.NewLabel(ui.Label{
		Text: "First Tab",
		Font: TabFont,
	}))
	{
		label := ui.NewLabel(ui.Label{
			Text: "Hello world",
			Font: render.Text{
				Size:  24,
				Color: render.SkyBlue,
			},
		})
		tab1.Pack(label, ui.Pack{
			Side: ui.N,
		})

		label2 := ui.NewLabel(ui.Label{
			Text: "This is the text content of the first\n" +
				"of the three tab frames.",
			Font: render.Text{
				Size:  10,
				Color: render.SkyBlue.Darken(40),
			},
		})
		tab1.Pack(label2, ui.Pack{
			Side: ui.N,
			PadY: 8,
		})
	}

	// Second Tab.
	tab2 := notebook.AddTab("Tab 2", ui.NewLabel(ui.Label{
		Text: "Second",
		Font: TabFont,
	}))
	{
		label := ui.NewLabel(ui.Label{
			Text: "Goodbye Mars",
			Font: render.Text{
				Size:  24,
				Color: render.Orange,
			},
		})
		tab2.Pack(label, ui.Pack{
			Side: ui.N,
		})

		label2 := ui.NewLabel(ui.Label{
			Text: "This is the text content of the second\n" +
				"of the three tab frames.\n\nIt has longer text\nin it!",
			Font: render.Text{
				Size:  10,
				Color: render.Orange.Darken(20),
			},
		})
		tab2.Pack(label2, ui.Pack{
			Side: ui.N,
			PadY: 8,
		})
	}

	// Third Tab.
	tab3 := notebook.AddTab("Tab 3", ui.NewLabel(ui.Label{
		Text: "Third",
		Font: TabFont,
	}))
	{
		label := ui.NewLabel(ui.Label{
			Text: "The Third Tab",
			Font: render.Text{
				Size:  24,
				Color: render.Pink,
			},
		})
		tab3.Pack(label, ui.Pack{
			Side: ui.N,
		})

		label2 := ui.NewLabel(ui.Label{
			Text: "This is the text content of the third\n" +
				"of the tab frames.",
			Font: render.Text{
				Size:  10,
				Color: render.Pink.Darken(40),
			},
		})
		tab3.Pack(label2, ui.Pack{
			Side: ui.N,
			PadY: 8,
		})
	}

	notebook.Supervise(mw.Supervisor())

	// notebook.SetBackground(render.DarkGrey)
	return notebook
}

func main() {
	mw, err := ui.NewMainWindow("TabFrame Example", Width, Height)
	if err != nil {
		panic(err)
	}

	// Dark theme.
	// ui.Theme = theme.DefaultDark

	// Menu bar.
	menu := ui.NewMenuBar("Main Menu")
	file := menu.AddMenu("UI Theme")
	file.AddItem("Default", func() {
		ui.Theme = theme.Default
		addWindow(mw)
	})
	file.AddItem("DefaultFlat", func() {
		ui.Theme = theme.DefaultFlat
		addWindow(mw)
	})
	file.AddItem("DefaultDark", func() {
		ui.Theme = theme.DefaultDark
		addWindow(mw)
	})
	file.AddSeparator()
	file.AddItem("Close all windows", func() {
		OpenWindows -= mw.Supervisor().CloseAllWindows()
	})

	menu.Supervise(mw.Supervisor())
	menu.Compute(mw.Engine)
	mw.Pack(menu, menu.PackTop())

	// Add some windows to play with.
	addWindow(mw)
	addWindow(mw)

	mw.SetBackground(render.White)

	mw.OnLoop(func(e *event.State) {
		if e.Escape {
			os.Exit(0)
		}
	})

	mw.MainLoop()
}

// Add a new child window.
func addWindow(mw *ui.MainWindow) {
	var (
		color = WindowColors[WindowID%len(WindowColors)]
		title = fmt.Sprintf("Window %d (%s)", WindowID+1, ui.Theme.Name)
	)
	WindowID++

	win1 := ui.NewWindow(title)
	win1.SetButtons(ui.CloseButton)
	win1.ActiveTitleBackground = color
	win1.InactiveTitleBackground = color.Darken(60)
	win1.InactiveTitleForeground = render.Grey
	win1.Configure(ui.Config{
		Width:  320,
		Height: 240,
	})
	win1.Compute(mw.Engine)
	win1.Supervise(mw.Supervisor())

	// Re-open a window when the last one is closed.
	OpenWindows++
	win1.Handle(ui.CloseWindow, func(ed ui.EventData) error {
		OpenWindows--
		if OpenWindows <= 0 {
			addWindow(mw)
		}
		return nil
	})

	// Default placement via cascade.
	win1.MoveTo(Cascade)
	Cascade.Add(CascadeStep)
	if Cascade.Y > Height-240-64 {
		CascadeLoops++
		Cascade.Y = 32
		Cascade.X = 24 * CascadeLoops
	}

	// Add the TabFrame.
	tabframe := MakeTabFrame(mw)
	win1.Pack(tabframe, ui.Pack{
		Side:   ui.W,
		Expand: true,
	})

	// Add a window duplicator button.
	btn2 := ui.NewButton(title+":Button2", ui.NewLabel(ui.Label{
		Text: "New Window",
	}))
	btn2.Handle(ui.Click, func(ed ui.EventData) error {
		addWindow(mw)
		return nil
	})
	btn2.Compute(mw.Engine)
	mw.Add(btn2)
	win1.Compute(mw.Engine)
	win1.Place(btn2, ui.Place{
		Bottom: 12,
		Right:  12,
	})
}
