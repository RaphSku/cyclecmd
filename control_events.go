package cyclecmd

const (
	// CYCLE_TERMINATE is a flag indicating the termination event.
	// It has the integer value 1.
	CYCLE_TERMINATE int = 1
)

// ControlEvent represents control signals with boolean flags.
//
// Currently, it contains a single flag:
//   - Terminate: true if the termination flag is set.
type ControlEvent struct {
	Terminate bool
}

// NewControlEvent creates a new ControlEvent from the given flags integer.
//
// It checks if the TERMINATE bit is set in the flags and sets the
// Terminate field accordingly.
//
// Parameters:
//   - `flags` : A number of control events whose bit should be set
//
// Returns:
//   - `*ControlEvent` : A ControlEvent structure that captures all the control events that are activated or deactivated
func NewControlEvent(flags int) *ControlEvent {
	controlEvent := &ControlEvent{}

	controlEvent.Terminate = false
	if flags&CYCLE_TERMINATE != 0 {
		controlEvent.Terminate = true
	}

	return controlEvent
}
