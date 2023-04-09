// Package magicform helps create simple form layouts with go/ui.
package magicform

import (
	"fmt"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/ui"
	"git.kirsle.net/go/ui/style"
)

type Type int

const (
	Auto   Type = iota
	Text        // free, wide Label row
	Frame       // custom frame from the caller
	Button      // Single button with a label
	Value       // a Label & Value row (value not editable)
	Textbox
	Checkbox
	Radiobox
	Selectbox
	Listbox
	Color
	Pager
)

// Form configuration.
type Form struct {
	Supervisor *ui.Supervisor // Required for most useful forms
	Engine     render.Engine

	// For vertical forms.
	Vertical   bool
	LabelWidth int // size of left frame for labels.
	PadY       int // spacer between (vertical) forms
	PadX       int
}

/*
Field for your form (or form-aligned label sections, etc.)

The type of Form control to render is inferred based on bound
variables and other configuration.
*/
type Field struct {
	// Type may be inferred by presence of other params.
	Type Type

	// Set a text string and font for simple labels or paragraphs.
	Label         string
	LabelVariable *string // a TextVariable to drive the Label
	Font          render.Text

	// Easy button row: make Buttons an array of Button fields
	Buttons     []Field
	ButtonStyle *style.Button

	// Easy Paginator. DO NOT SUPERVISE, let the Create do so!
	Pager *ui.Pager

	// If you send a *ui.Frame to insert, the Type is inferred
	// to be Frame.
	Frame *ui.Frame

	// Variable bindings, the type may infer to be:
	BoolVariable *bool         // Checkbox
	TextVariable *string       // Textbox
	IntVariable  *int          // Textbox
	Options      []Option      // Selectbox
	SelectValue  interface{}   // Selectbox default choice
	Color        *render.Color // Color
	Readonly     bool          // draw the value as a flat label

	// For text-type fields, opt-in to let magicform prompt the
	// user using the game's developer shell.
	PromptUser func(answer string)

	// Tooltip to add to a form control.
	// Checkbox only for now.
	Tooltip ui.Tooltip // config for the tooltip only

	// Handlers you can configure
	OnSelect func(value interface{}) // Selectbox
	OnClick  func()                  // Button
}

// Option used in Selectbox or Radiobox fields.
type Option struct {
	Value     interface{}
	Label     string
	Separator bool
}

/*
Create the form field and populate it into the given Frame.

Renders the form vertically.
*/
func (form Form) Create(into *ui.Frame, fields []Field) {
	for n, row := range fields {
		row := row

		if row.Frame != nil {
			into.Pack(row.Frame, ui.Pack{
				Side:  ui.N,
				FillX: true,
			})
			continue
		}

		frame := ui.NewFrame(fmt.Sprintf("Line %d", n))
		into.Pack(frame, ui.Pack{
			Side:  ui.N,
			FillX: true,
			PadY:  form.PadY,
		})

		// Buttons row?
		if row.Buttons != nil && len(row.Buttons) > 0 {
			for _, row := range row.Buttons {
				row := row

				btn := ui.NewButton(row.Label, ui.NewLabel(ui.Label{
					Text: row.Label,
					Font: row.Font,
				}))
				if row.ButtonStyle != nil {
					btn.SetStyle(row.ButtonStyle)
				}

				btn.Handle(ui.Click, func(ed ui.EventData) error {
					if row.OnClick != nil {
						row.OnClick()
					} else {
						return fmt.Errorf("no OnClick handler for button %s", row.Label)
					}
					return nil
				})

				btn.Compute(form.Engine)
				form.Supervisor.Add(btn)

				// Tooltip? TODO - make nicer.
				if row.Tooltip.Text != "" || row.Tooltip.TextVariable != nil {
					tt := ui.NewTooltip(btn, row.Tooltip)
					tt.Supervise(form.Supervisor)
				}

				frame.Pack(btn, ui.Pack{
					Side: ui.W,
					PadX: 2,
					PadY: 2,
				})
			}

			continue
		}

		// Infer the type of the form field.
		if row.Type == Auto {
			row.Type = row.Infer()
			if row.Type == Auto {
				continue
			}
		}

		// Is there a label frame to the left?
		// - Checkbox gets a full row.
		fmt.Printf("Label=%+v  Var=%+v\n", row.Label, row.LabelVariable)
		if (row.Label != "" || row.LabelVariable != nil) && row.Type != Checkbox {
			labFrame := ui.NewFrame("Label Frame")
			labFrame.Configure(ui.Config{
				Width: form.LabelWidth,
			})
			frame.Pack(labFrame, ui.Pack{
				Side: ui.W,
			})

			// Draw the label text into it.
			label := ui.NewLabel(ui.Label{
				Text:         row.Label,
				TextVariable: row.LabelVariable,
				Font:         row.Font,
			})
			labFrame.Pack(label, ui.Pack{
				Side: ui.W,
			})
		}

		// Pager row?
		if row.Pager != nil {
			row.Pager.Supervise(form.Supervisor)
			frame.Pack(row.Pager, ui.Pack{
				Side: ui.W,
			})
		}

		// Simple "Value" row with a Label to its left.
		if row.Type == Value {
			lbl := ui.NewLabel(ui.Label{
				Text:         row.Label,
				Font:         row.Font,
				TextVariable: row.TextVariable,
				IntVariable:  row.IntVariable,
			})

			frame.Pack(lbl, ui.Pack{
				Side:   ui.W,
				FillX:  true,
				Expand: true,
			})

			// Tooltip? TODO - make nicer.
			if row.Tooltip.Text != "" || row.Tooltip.TextVariable != nil {
				tt := ui.NewTooltip(lbl, row.Tooltip)
				tt.Supervise(form.Supervisor)
			}
		}

		// Color picker button.
		if row.Type == Color && row.Color != nil {
			btn := ui.NewButton("ColorPicker", ui.NewLabel(ui.Label{
				Text: " ",
				Font: row.Font,
			}))
			style := style.DefaultButton
			style.Background = *row.Color
			style.HoverBackground = style.Background.Lighten(20)
			btn.SetStyle(&style)

			form.Supervisor.Add(btn)
			frame.Pack(btn, ui.Pack{
				Side:   ui.W,
				FillX:  true,
				Expand: true,
			})

			btn.Handle(ui.Click, func(ed ui.EventData) error {
				// Open a ColorPicker widget.
				picker, err := ui.NewColorPicker(ui.ColorPicker{
					Title:      "Select a color",
					Supervisor: form.Supervisor,
					Engine:     form.Engine,
					Color:      *row.Color,
					OnManualInput: func(callback func(render.Color)) {
						// TODO: prompt for color
					},
				})
				if err != nil {
					return err
				}

				picker.Then(func(color render.Color) {
					*row.Color = color
					style.Background = color
					style.HoverBackground = style.Background.Lighten(20)

					// call onClick to save change to disk now
					if row.OnClick != nil {
						row.OnClick()
					}
				})

				picker.Center(form.Engine.WindowSize())
				picker.Show()
				return nil
			})
		}

		// Buttons and Text fields (for now).
		if row.Type == Button || row.Type == Textbox {
			btn := ui.NewButton("Button", ui.NewLabel(ui.Label{
				Text:         row.Label,
				Font:         row.Font,
				TextVariable: row.TextVariable,
				IntVariable:  row.IntVariable,
			}))

			frame.Pack(btn, ui.Pack{
				Side:   ui.W,
				FillX:  true,
				Expand: true,
			})

			// Not clickable if Readonly.
			if !row.Readonly {
				form.Supervisor.Add(btn)
			}

			// Tooltip? TODO - make nicer.
			if row.Tooltip.Text != "" || row.Tooltip.TextVariable != nil {
				tt := ui.NewTooltip(btn, row.Tooltip)
				tt.Supervise(form.Supervisor)
			}

			// Handlers
			btn.Handle(ui.Click, func(ed ui.EventData) error {
				// Text boxes, we want to prompt the user to enter new value?
				if row.PromptUser != nil {
					var value string
					if row.TextVariable != nil {
						value = *row.TextVariable
					} else if row.IntVariable != nil {
						value = fmt.Sprintf("%d", *row.IntVariable)
					}

					// TODO: prompt user for new value
					_ = value
					// shmem.PromptPre("Enter new value: ", value, func(answer string) {
					// 	if answer != "" {
					// 		row.PromptUser(answer)
					// 	}
					// })
				}

				if row.OnClick != nil {
					row.OnClick()
				}
				return nil
			})
		}

		// Checkbox?
		if row.Type == Checkbox {
			cb := ui.NewCheckbox("Checkbox", row.BoolVariable, ui.NewLabel(ui.Label{
				Text: row.Label,
				Font: row.Font,
			}))
			cb.Supervise(form.Supervisor)
			frame.Pack(cb, ui.Pack{
				Side:  ui.W,
				FillX: true,
			})

			// Tooltip? TODO - make nicer.
			if row.Tooltip.Text != "" || row.Tooltip.TextVariable != nil {
				tt := ui.NewTooltip(cb, row.Tooltip)
				tt.Supervise(form.Supervisor)
			}

			// Handlers
			cb.Handle(ui.Click, func(ed ui.EventData) error {
				if row.OnClick != nil {
					row.OnClick()
				}
				return nil
			})
		}

		// Selectbox? also Radiobox for now.
		if row.Type == Selectbox || row.Type == Radiobox {
			btn := ui.NewSelectBox("Select", ui.Label{
				Font: row.Font,
			})
			frame.Pack(btn, ui.Pack{
				Side:   ui.W,
				FillX:  true,
				Expand: true,
			})

			if row.Options != nil {
				for _, option := range row.Options {
					if option.Separator {
						btn.AddSeparator()
						continue
					}
					btn.AddItem(option.Label, option.Value, func() {})
				}
			}

			if row.SelectValue != nil {
				btn.SetValue(row.SelectValue)
			}

			btn.Handle(ui.Change, func(ed ui.EventData) error {
				if selection, ok := btn.GetValue(); ok {
					if row.OnSelect != nil {
						row.OnSelect(selection.Value)
					}

					// Update bound variables.
					if v, ok := selection.Value.(int); ok && row.IntVariable != nil {
						*row.IntVariable = v
					}
					if v, ok := selection.Value.(string); ok && row.TextVariable != nil {
						*row.TextVariable = v
					}
				}
				return nil
			})

			// Tooltip? TODO - make nicer.
			if row.Tooltip.Text != "" || row.Tooltip.TextVariable != nil {
				tt := ui.NewTooltip(btn, row.Tooltip)
				tt.Supervise(form.Supervisor)
			}

			btn.Supervise(form.Supervisor)
			form.Supervisor.Add(btn)
		}

		// ListBox?
		if row.Type == Listbox {
			btn := ui.NewListBox("List", ui.ListBox{
				Variable: row.SelectValue,
			})
			btn.Configure(ui.Config{
				Height: 120,
			})
			frame.Pack(btn, ui.Pack{
				Side:   ui.W,
				FillX:  true,
				Expand: true,
			})

			if row.Options != nil {
				for _, option := range row.Options {
					if option.Separator {
						// btn.AddSeparator()
						continue
					}
					fmt.Printf("LISTBOX: Insert label '%s' with value %+v\n", option.Label, option.Value)
					btn.AddLabel(option.Label, option.Value, func() {})
				}
			}

			if row.SelectValue != nil {
				fmt.Printf("LISTBOX: Set value to %s\n", row.SelectValue)
				btn.SetValue(row.SelectValue)
			}

			btn.Handle(ui.Change, func(ed ui.EventData) error {
				if selection, ok := btn.GetValue(); ok {
					if row.OnSelect != nil {
						row.OnSelect(selection.Value)
					}

					// Update bound variables.
					if v, ok := selection.Value.(int); ok && row.IntVariable != nil {
						*row.IntVariable = v
					}
					if v, ok := selection.Value.(string); ok && row.TextVariable != nil {
						*row.TextVariable = v
					}
				}
				return nil
			})

			// Tooltip? TODO - make nicer.
			if row.Tooltip.Text != "" || row.Tooltip.TextVariable != nil {
				tt := ui.NewTooltip(btn, row.Tooltip)
				tt.Supervise(form.Supervisor)
			}

			btn.Supervise(form.Supervisor)
			// form.Supervisor.Add(btn) // for btn.Handle(Change) to work??
		}
	}
}

/*
Infer the type if the field was of type Auto.

Returns the first Type inferred from the field by checking in
this order:

- Frame if the field has a *Frame
- Checkbox if there is a *BoolVariable
- Selectbox if there are Options
- Textbox if there is a *TextVariable
- Text if there is a Label

May return Auto if none of the above and be ignored.
*/
func (field Field) Infer() Type {
	if field.Frame != nil {
		return Frame
	}

	if field.BoolVariable != nil {
		return Checkbox
	}

	if field.Options != nil && len(field.Options) > 0 {
		return Selectbox
	}

	if field.TextVariable != nil || field.IntVariable != nil {
		return Textbox
	}

	if field.Label != "" {
		return Text
	}

	if field.Pager != nil {
		return Pager
	}

	return Auto
}
