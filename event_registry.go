package cyclecmd

import (
	"fmt"
)

// nonDefaultEvent is used when an event trigger is not registered and the byte length of
// the token is greater than 3.
type nonDefaultEvent struct{}

// Handle does for nonDefaultEvent nothing since no event trigger has been specified.
//
// Parameters:
//   - `token` : Token that should be handled given as a string
//
// Returns:
//   - `error` : Returns no error in this case
//   - `*ControlEvent` : Returns no control event in this case
func (nde *nonDefaultEvent) Handle(token string) (error, *ControlEvent) {
	return nil, nil
}

// EventRegistry contains all related information with custom events that are registered with
// this registry.
type EventRegistry struct {
	nonDefaultEventInformation EventInformation

	// registry is a key-value data structure, the key contains the event trigger and the value contains
	// the EventInformation related to the event that was triggered by an event
	registry map[string]EventInformation

	// DefaultEventInformation contains information related to the default event that is triggered whenever
	// a token does not match with any other event that is registered.
	DefaultEventInformation EventInformation
}

// NewEventRegistry initialises the event registry.
//
// Parameters:
//   - `defaultEventInformation` : Information related to the default event
//
// Returns:
//   - `*EventRegistry` : Returns an instance of the event registry
func NewEventRegistry(defaultEventInformation EventInformation) *EventRegistry {
	eventRegistry := &EventRegistry{
		DefaultEventInformation: defaultEventInformation,
	}
	eventRegistry.registry = make(map[string]EventInformation)

	eventRegistry.nonDefaultEventInformation = EventInformation{
		EventName: "NonDefault",
		Event:     &nonDefaultEvent{},
	}

	return eventRegistry
}

// ResetEventRegistry resets the event registry, so all registered events are deleted from the registry.
func (er *EventRegistry) ResetEventRegistry() {
	er.registry = make(map[string]EventInformation)
}

// RegisterEvent registers an event with an event trigger.
//
// Parameters:
//   - `eventTrigger` : Trigger that will kick off the event
//   - `eventInformation` : Information related to the event that will be triggered by `eventTrigger`
//
// Returns:
//   - `error` : Returns an error when the event is already registered
func (er *EventRegistry) RegisterEvent(eventTrigger string, eventInformation EventInformation) error {
	_, ok := er.registry[eventTrigger]
	if ok {
		return fmt.Errorf("event is already registered under event trigger %v", eventTrigger)
	}
	er.registry[eventTrigger] = eventInformation
	return nil
}

// GetMatchingEventInformation retrieves the information related to the event that gets triggered by `eventTrigger`.
// The default event is returned when the event trigger matches no event registered in the event registry.
//
// Parameters:
//   - `eventTrigger` : Trigger for the event that should be returned
//
// Returns:
//   - `EventInformation` : Information related to the event triggered by `eventTrigger`
//   - `error` : An error is only returned when no default event is defined
func (er *EventRegistry) GetMatchingEventInformation(eventTrigger string) (EventInformation, error) {
	eventInformation, ok := er.registry[eventTrigger]
	if !ok {
		if len([]byte(eventTrigger)) == 1 {
			defaultEvent := er.DefaultEventInformation.Event
			if defaultEvent == nil {
				return EventInformation{}, fmt.Errorf("default event is not set! Please set it via InitEventRegistry")
			}
			return er.DefaultEventInformation, nil
		}
		if len([]byte(eventTrigger)) > 1 {
			return er.nonDefaultEventInformation, nil
		}
	}
	return eventInformation, nil
}
