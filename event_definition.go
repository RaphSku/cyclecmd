package cyclecmd

// Event is an interface that defines the behavior of the custom events
// that the User can define themselves.
//
// Event expects the following method to be implemented by all events:
//
// Behavior:
//   - `Handle(token string) (error, *ControlEvent)` : it expects the token that is associated with the event
type Event interface {
	Handle(token string) (error, *ControlEvent)
}

// EventInformation stores the event itself but also the event name that was given to the custom event.
type EventInformation struct {
	// Name of the event
	EventName string
	// The event instance itself
	Event Event
}

// EventHistoryEntry stores the event and the event name but also the token that triggered the event.
// This is especially useful for the DefaultEvent since that event gets triggered by every token
// that is not already registered with another event.
type EventHistoryEntry struct {
	// Token is the same as event trigger
	Token string
	// The name of the event
	EventName string
	// The event instance itself
	Event Event
}
