//go:build unit_test

package cyclecmd_test

import (
	"fmt"
	"testing"

	"github.com/RaphSku/cyclecmd"
	"github.com/stretchr/testify/assert"
)

func TestInitDefaultEvent(t *testing.T) {
	t.Parallel()
	defaultEventInformation := setupDefaultEventInformation()
	eventRegistry := cyclecmd.NewEventRegistry(defaultEventInformation)

	actDefaultEventInformation, err := eventRegistry.GetMatchingEventInformation("d")
	assert.NoError(t, err)
	assert.Equal(t, defaultEventInformation.Event, actDefaultEventInformation.Event)
}

func TestEventRegistration(t *testing.T) {
	t.Parallel()

	defaultEventInformation := setupDefaultEventInformation()
	eventRegistry := cyclecmd.NewEventRegistry(defaultEventInformation)

	expEventTrigger := "t"
	expEventInformation := cyclecmd.EventInformation{
		EventName: "Test",
		Event:     &TestEvent{},
	}
	err := eventRegistry.RegisterEvent(expEventTrigger, expEventInformation)
	assert.NoError(t, err)

	actDefaultEventInformation, err := eventRegistry.GetMatchingEventInformation(expEventTrigger)
	assert.NoError(t, err)
	assert.Equal(t, &TestEvent{}, actDefaultEventInformation.Event)
}

func TestResetEventRegistry(t *testing.T) {
	t.Parallel()

	defaultEventInformation := setupDefaultEventInformation()
	eventRegistry := cyclecmd.NewEventRegistry(defaultEventInformation)

	eventTrigger := "t"
	eventInformation := cyclecmd.EventInformation{
		EventName: "Test",
		Event:     &TestEvent{},
	}
	err := eventRegistry.RegisterEvent(eventTrigger, eventInformation)
	assert.NoError(t, err)

	eventRegistry.ResetEventRegistry()

	// If the registry reset was successful, the default event information should be returned
	// instead of the event information belonging to the trigger "t"
	actEventInformation, err := eventRegistry.GetMatchingEventInformation(eventTrigger)
	assert.Equal(t, defaultEventInformation, actEventInformation)
}

func TestDuplicateEventRegistration(t *testing.T) {
	t.Parallel()

	defaultEventInformation := setupDefaultEventInformation()
	eventRegistry := cyclecmd.NewEventRegistry(defaultEventInformation)

	eventTrigger := "t"
	eventInformation := cyclecmd.EventInformation{
		EventName: "Test",
		Event:     &TestEvent{},
	}
	err := eventRegistry.RegisterEvent(eventTrigger, eventInformation)
	assert.NoError(t, err)

	err = eventRegistry.RegisterEvent(eventTrigger, eventInformation)
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("event is already registered under event trigger %v", eventTrigger), err)
	}
}

func TestNoDefaultEventInRegistry(t *testing.T) {
	t.Parallel()

	eventRegistry := &cyclecmd.EventRegistry{}
	actEventInformation, err := eventRegistry.GetMatchingEventInformation("\x1b")
	assert.Equal(t, cyclecmd.EventInformation{}, actEventInformation)
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("default event is not set! Please set it via InitEventRegistry"), err)
	}
}
