package main

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/sdl"
	"git.kirsle.net/go/ui"
	"git.kirsle.net/go/ui/magicform"
	"git.kirsle.net/go/ui/style"
)

func init() {
	sdl.DefaultFontFilename = "../DejaVuSans.ttf"
}

var (
	MenuFont = render.Text{
		Size: 12,
		PadX: 4,
		PadY: 2,
	}
	TabFont = render.Text{
		Size: 12,
		PadX: 4,
		PadY: 2,
	}
)

var ButtonStylePrimary = &style.Button{
	Background:      render.RGBA(0, 60, 153, 255),
	Foreground:      render.White,
	HoverBackground: render.RGBA(0, 153, 255, 255),
	HoverForeground: render.White,
	OutlineColor:    render.DarkGrey,
	OutlineSize:     1,
	BorderStyle:     style.BorderRaised,
	BorderSize:      2,
}

func main() {
	mw, err := ui.NewMainWindow("Forms Test", 500, 375)
	if err != nil {
		panic(err)
	}

	// Tabbed UI.
	tabFrame := ui.NewTabFrame("Tabs")
	makeAppFrame(mw, tabFrame)
	makeAboutFrame(mw, tabFrame)

	tabFrame.Supervise(mw.Supervisor())
	mw.Pack(tabFrame, ui.Pack{
		Side:    ui.N,
		Expand:  true,
		Padding: 10,
	})

	mw.SetBackground(render.Grey)
	mw.MainLoop()
}

func makeAppFrame(mw *ui.MainWindow, tf *ui.TabFrame) *ui.Frame {
	frame := tf.AddTab("Index", ui.NewLabel(ui.Label{
		Text: "Form Controls",
		Font: TabFont,
	}))

	// Form variables
	var (
		bgcolor    = render.Grey
		letter     string
		checkBool1 bool
		checkBool2 = true
		pagerLabel = "Page 1 of 20"
	)

	// Magic Form is a handy module for easily laying out forms of widgets.
	form := magicform.Form{
		Supervisor: mw.Supervisor(),
		Engine:     mw.Engine,
		Vertical:   true,
		LabelWidth: 120,
		PadY:       2,
		PadX:       8,
	}

	// You add to it a list of fields which support all sorts of different
	// form control types.
	fields := []magicform.Field{
		// Simple text sections - you can write paragraphs or use a bold font
		// to make section labels that span the full width of your frame.
		{
			Label: "Checkbox controls bound to bool values:",
			Font:  MenuFont,
		},

		// Checkbox widgets: just bind a BoolVariable and this row will draw
		// with a checkbox next to a label.
		{
			Label:        "Check this box to toggle a boolean",
			Font:         MenuFont,
			BoolVariable: &checkBool1,
			OnClick: func() {
				fmt.Printf("The checkbox was clicked! Value is now: %+v\n", checkBool1)
			},
		},
		{
			Label:        "Uncheck this one",
			Font:         MenuFont,
			BoolVariable: &checkBool2,
			OnClick: func() {
				fmt.Printf("The checkbox was clicked! Value is now: %+v\n", checkBool1)
			},
		},

		// SelectBox widgets: just bind a SelectValue and provide Options and
		// it will draw with a label (LabelWidth wide) next to a SelectBox button.
		{
			Label:       "Window color:",
			Font:        MenuFont,
			SelectValue: &bgcolor,
			Options: []magicform.Option{
				{
					Label: "Grey",
					Value: render.Grey,
				},
				{
					Label: "White",
					Value: render.White,
				},
				{
					Label: "Yellow",
					Value: render.Yellow,
				},
				{
					Label: "Cyan",
					Value: render.Cyan,
				},
				{
					Label: "Green",
					Value: render.Green,
				},
				{
					Label: "Blue",
					Value: render.RGBA(0, 153, 255, 255),
				},
				{
					Label: "Pink",
					Value: render.Pink,
				},
			},
			OnSelect: func(v interface{}) {
				value, _ := v.(render.Color)
				mw.SetBackground(value)
			},
		},

		// ListBox widgets
		{
			Type:        magicform.Listbox,
			Label:       "Favorite letter:",
			Font:        MenuFont,
			SelectValue: &letter,
			Options: []magicform.Option{
				{Label: "A is for apple", Value: "A"},
				{Label: "B is for boy", Value: "B"},
				{Label: "C is for cat", Value: "C"},
				{Label: "D is for dog", Value: "D"},
				{Label: "E is for elephant", Value: "E"},
				{Label: "F is for far", Value: "F"},
				{Label: "G is for ghost", Value: "G"},
				{Label: "H is for high", Value: "H"},
				{Label: "I is for inside", Value: "I"},
				{Label: "J is for joker", Value: "J"},
				{Label: "K is for kangaroo", Value: "K"},
				{Label: "L is for lion", Value: "L"},
				{Label: "M is for mouse", Value: "M"},
				{Label: "N is for night", Value: "N"},
				{Label: "O is for over", Value: "O"},
				{Label: "P is for parry", Value: "P"},
				{Label: "Q is for quarry", Value: "Q"},
				{Label: "R is for reality", Value: "R"},
				{Label: "S is for sunshine", Value: "S"},
				{Label: "T is for tree", Value: "T"},
				{Label: "U is for under", Value: "U"},
				{Label: "V is for vehicle", Value: "V"},
				{Label: "W is for watermelon", Value: "W"},
				{Label: "X is for xylophone", Value: "X"},
				{Label: "Y is for yellow", Value: "Y"},
				{Label: "Z is for zebra", Value: "Z"},
			},
			OnSelect: func(v interface{}) {
				value, _ := v.(string)
				fmt.Printf("You clicked on: %s\n", value)
			},
		},

		// Pager rows to show an easy paginated UI.
		// TODO: this is currently broken and Supervisor doesn't pick it up
		{
			Label: "A paginator when you need one. You can limit MaxPageButtons\n" +
				"and the right arrow can keep selecting past the last page.",
		},
		{
			LabelVariable: &pagerLabel,
			Label:         "Page:",
			Pager: ui.NewPager(ui.Pager{
				Page:           1,
				Pages:          20,
				PerPage:        10,
				MaxPageButtons: 8,
				Font:           MenuFont,
				OnChange: func(page, perPage int) {
					fmt.Printf("Pager clicked: page=%d perPage=%d\n", page, perPage)
					pagerLabel = fmt.Sprintf("Page %d of %d", page, 20)
				},
			}),
		},

		// Simple variable bindings.
		{
			Type:         magicform.Value,
			Label:        "The first bool var:",
			TextVariable: &pagerLabel,
		},

		// Buttons for the bottom of your form.
		{
			Buttons: []magicform.Field{
				{
					Label:       "Save",
					ButtonStyle: ButtonStylePrimary,
					Font:        MenuFont,
					OnClick: func() {
						fmt.Println("Primary button clicked")
					},
				},
				{
					Label: "Cancel",
					Font:  MenuFont,
					OnClick: func() {
						fmt.Println("Secondary button clicked")
					},
				},
			},
		},
	}

	form.Create(frame, fields)
	return frame
}

func makeAboutFrame(mw *ui.MainWindow, tf *ui.TabFrame) *ui.Frame {
	frame := tf.AddTab("About", ui.NewLabel(ui.Label{
		Text: "About",
		Font: TabFont,
	}))

	form := magicform.Form{
		Supervisor: mw.Supervisor(),
		Engine:     mw.Engine,
		Vertical:   true,
		LabelWidth: 120,
		PadY:       2,
		PadX:       8,
	}

	fields := []magicform.Field{
		{
			Label: "About",
			Font:  MenuFont,
		},

		{
			Label: "This example shows off the UI toolkit's use for form controls,\n" +
				"and how the magicform helper module can make simple forms\n" +
				"easy to compose quickly.",
			Font: MenuFont,
		},
	}

	form.Create(frame, fields)
	return frame
}
