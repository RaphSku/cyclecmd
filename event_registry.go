package cyclecmd

import "fmt"

// EventRegistry contains all related information with custom events that are registered with
// this registry.
type EventRegistry struct {
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
//   - `*EventRegistry` - Returns an instance of the event registry
func NewEventRegistry(defaultEventInformation EventInformation) *EventRegistry {
	eventRegistry := &EventRegistry{
		DefaultEventInformation: defaultEventInformation,
	}
	eventRegistry.registry = make(map[string]EventInformation)

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
	err := er.validateEventTrigger(eventTrigger)
	if err != nil {
		return err
	}
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
	err := er.validateEventTrigger(eventTrigger)
	if err != nil {
		return EventInformation{}, err
	}
	eventInformation, ok := er.registry[eventTrigger]
	if !ok {
		defaultEvent := er.DefaultEventInformation.Event
		if defaultEvent == nil {
			return EventInformation{}, fmt.Errorf("default event is not set! Please set it via InitEventRegistry")
		}
		return er.DefaultEventInformation, nil
	}
	return eventInformation, nil
}

// validateEventTrigger validates whether the eventTrigger represents just one token
//
// Parameters:
//   - `eventTrigger` : String that triggers an event
//
// Returns:
//   - `error` : Returns an error if the length of []byte(eventTrigger) is not equal to 1 (so no token)
func (er *EventRegistry) validateEventTrigger(eventTrigger string) error {
	if len([]byte(eventTrigger)) != 1 {
		return fmt.Errorf("ensure that the eventTrigger represents just one token ([]byte length of 1)")
	}
	return nil
}
