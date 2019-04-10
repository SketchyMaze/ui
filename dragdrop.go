package ui

// DragDrop is a state machine to manage draggable UI components.
type DragDrop struct {
	isDragging bool
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

// Start the drag state.
func (dd *DragDrop) Start() {
	dd.isDragging = true
}

// Stop dragging.
func (dd *DragDrop) Stop() {
	dd.isDragging = false
}
