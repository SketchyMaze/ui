package ui

// DragDrop is a state machine to manage draggable UI components.
type DragDrop struct {
	isDragging bool

	// If the subject of the drag is a widget, it can store itself here.
	widget Widget
}

// NewDragDrop initializes the DragDrop struct. Normally your Supervisor
// will manage the drag/drop object, but you can use your own if you don't
// use a Supervisor.
func NewDragDrop() *DragDrop {
	return &DragDrop{}
}

// IsDragging returns whether the drag state is active.
func (dd *DragDrop) IsDragging() bool {
	return dd.isDragging
}

// SetWidget attaches the widget to the drag state, but does not start the
// drag; you call Start() after this if the subject is a widget.
func (dd *DragDrop) SetWidget(w Widget) {
	dd.widget = w
}

// Widget returns the attached widget or nil.
func (dd *DragDrop) Widget() Widget {
	return dd.widget
}

// Start the drag state.
func (dd *DragDrop) Start() {
	dd.isDragging = true
}

// Stop dragging. This will also clear the stored widget, if any.
func (dd *DragDrop) Stop() {
	dd.isDragging = false
	dd.widget = nil
}
