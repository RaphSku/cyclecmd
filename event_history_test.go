//go:build unit_test

package cyclecmd_test

import (
	"fmt"
	"testing"

	"github.com/RaphSku/cyclecmd"
	"github.com/stretchr/testify/assert"
)

func setupPopulatedEventHistory() (*cyclecmd.EventHistory, error) {
	eventHistory := cyclecmd.NewEventHistory()
	eventHistoryEntryA := cyclecmd.EventHistoryEntry{
		Token:     "a",
		EventName: "A",
		Event:     &TestEvent{},
	}
	eventHistory.AddEvent(eventHistoryEntryA)
	eventHistoryEntryB := cyclecmd.EventHistoryEntry{
		Token:     "b",
		EventName: "B",
		Event:     &TestEvent{},
	}
	eventHistory.AddEvent(eventHistoryEntryB)
	eventHistoryEntryC := cyclecmd.EventHistoryEntry{
		Token:     "c",
		EventName: "C",
		Event:     &TestEvent{},
	}
	eventHistory.AddEvent(eventHistoryEntryC)

	return eventHistory, nil
}

func TestEventHistoryLength(t *testing.T) {
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

	eventHistory := cyclecmd.NewEventHistory()
	assert.Equal(t, 0, eventHistory.Len())

	eventEntry := cyclecmd.EventHistoryEntry{
		Token:     "t",
		EventName: "Test",
		Event:     &TestEvent{},
	}
	eventHistory.AddEvent(eventEntry)
	assert.Equal(t, 1, eventHistory.Len())
}

func TestRemovingNthElementFromHistory(t *testing.T) {
	t.Parallel()

	eventHistory, err := setupPopulatedEventHistory()
	assert.NoError(t, err)

	eventHistory.RemoveNthEventFromHistory(1)

	expEventHistoryEntry0 := cyclecmd.EventHistoryEntry{
		Token:     "a",
		EventName: "A",
		Event:     &TestEvent{},
	}
	actEventHistory0, err := eventHistory.RetrieveEventEntryByIndex(0)
	assert.NoError(t, err)
	assert.Equal(t, expEventHistoryEntry0, actEventHistory0)

	expEventHistoryEntry1 := cyclecmd.EventHistoryEntry{
		Token:     "c",
		EventName: "C",
		Event:     &TestEvent{},
	}
	actEventHistory1, err := eventHistory.RetrieveEventEntryByIndex(1)
	assert.NoError(t, err)
	assert.Equal(t, expEventHistoryEntry1, actEventHistory1)

	eventHistory.RemoveNthEventFromHistory(3)
	eventHistory.RemoveNthEventFromHistory(-1)
	actEventHistory0, err = eventHistory.RetrieveEventEntryByIndex(0)
	assert.NoError(t, err)
	actEventHistory1, err = eventHistory.RetrieveEventEntryByIndex(1)
	assert.NoError(t, err)
	assert.Equal(t, expEventHistoryEntry0, actEventHistory0)
	assert.Equal(t, expEventHistoryEntry1, actEventHistory1)
}

func TestRetrieveEventEntryByIndex(t *testing.T) {
	t.Parallel()

	eventHistory, err := setupPopulatedEventHistory()
	assert.NoError(t, err)

	expEventHistoryEntry0 := cyclecmd.EventHistoryEntry{
		Token:     "a",
		EventName: "A",
		Event:     &TestEvent{},
	}
	actEventHistoryEntry0, err := eventHistory.RetrieveEventEntryByIndex(0)
	assert.NoError(t, err)
	assert.Equal(t, expEventHistoryEntry0, actEventHistoryEntry0)

	actEventHistoryEntry3, err := eventHistory.RetrieveEventEntryByIndex(3)
	assert.Equal(t, cyclecmd.EventHistoryEntry{}, actEventHistoryEntry3)
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("index %v error, index is either smaller than 0 or larger than the length of the event history", 3), err)
	}
}

func TestPrintLastEventHistoryEntries(t *testing.T) {
	t.Parallel()

	eventHistory, err := setupPopulatedEventHistory()
	assert.NoError(t, err)

	actOutput, err := captureStdOutput(func() {
		eventHistory.PrintLastEventHistoryEntries(2)
	})
	assert.NoError(t, err)
	expEventHistoryEntryC := cyclecmd.EventHistoryEntry{
		Token:     "c",
		EventName: "C",
		Event:     &TestEvent{},
	}
	expEventHistoryEntryB := cyclecmd.EventHistoryEntry{
		Token:     "b",
		EventName: "B",
		Event:     &TestEvent{},
	}
	expOutput1 := fmt.Sprintf("Event Name: %s, Token: %s\r\n", expEventHistoryEntryC.EventName, expEventHistoryEntryC.Token)
	expOutput2 := fmt.Sprintf("Event Name: %s, Token: %s\r\n", expEventHistoryEntryB.EventName, expEventHistoryEntryB.Token)
	expOutput := expOutput1 + expOutput2
	assert.Equal(t, expOutput, actOutput)
}

func TestGetLastEventsFromHistoryToEventReference(t *testing.T) {
	t.Parallel()

	eventHistory, err := setupPopulatedEventHistory()
	assert.NoError(t, err)

	expEvents := []string{"B", "C"}
	actEvents := eventHistory.GetLastEventsFromHistoryToEventReference("A")
	assert.Equal(t, expEvents, actEvents)
}

func TestMostRecentSpliceEventsOfHistory(t *testing.T) {
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

	backspaceEventTrigger := "\x7f"
	backspaceEventInformation := cyclecmd.EventInformation{
		EventName: "Backspace",
		Event:     &BackspaceEvent{},
	}
	err = eventRegistry.RegisterEvent(backspaceEventTrigger, backspaceEventInformation)
	assert.NoError(t, err)

	eventHistory := cyclecmd.NewEventHistory()

	eventEntries := []cyclecmd.EventHistoryEntry{
		{Token: "a", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "t", EventName: "Test", Event: &TestEvent{}},
		{Token: "h", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "e", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "y", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "t", EventName: "Test", Event: &TestEvent{}},
		{Token: "b", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "c", EventName: "Default", Event: &DefaultEvent{}},
	}
	for _, eventEntry := range eventEntries {
		eventHistory.AddEvent(eventEntry)
	}

	expSplicedEvents := []cyclecmd.EventHistoryEntry{
		{Token: "h", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "e", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "y", EventName: "Default", Event: &DefaultEvent{}},
	}
	actSplicedEvents := eventHistory.MostRecentSpliceEventsOfHistory("Test")

	for i := range actSplicedEvents {
		assert.Equal(t, expSplicedEvents[i].Token, actSplicedEvents[i].Token)
		assert.Equal(t, expSplicedEvents[i].EventName, actSplicedEvents[i].EventName)
		assert.Equal(t, expSplicedEvents[i].Event, actSplicedEvents[i].Event)
	}
}

func TestMostRecentSpliceEventsOfHistoryWithNoEnd(t *testing.T) {
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

	backspaceEventTrigger := "\x7f"
	backspaceEventInformation := cyclecmd.EventInformation{
		EventName: "Backspace",
		Event:     &BackspaceEvent{},
	}
	err = eventRegistry.RegisterEvent(backspaceEventTrigger, backspaceEventInformation)
	assert.NoError(t, err)

	eventHistory := cyclecmd.NewEventHistory()

	eventEntries := []cyclecmd.EventHistoryEntry{
		{Token: "a", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "h", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "e", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "y", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "t", EventName: "Test", Event: &TestEvent{}},
	}
	for _, eventEntry := range eventEntries {
		eventHistory.AddEvent(eventEntry)
	}

	expSplicedEvents := []cyclecmd.EventHistoryEntry{
		{Token: "a", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "h", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "e", EventName: "Default", Event: &DefaultEvent{}},
		{Token: "y", EventName: "Default", Event: &DefaultEvent{}},
	}
	actSplicedEvents := eventHistory.MostRecentSpliceEventsOfHistory("Test")

	for i := range actSplicedEvents {
		assert.Equal(t, expSplicedEvents[i].Token, actSplicedEvents[i].Token)
		assert.Equal(t, expSplicedEvents[i].EventName, actSplicedEvents[i].EventName)
		assert.Equal(t, expSplicedEvents[i].Event, actSplicedEvents[i].Event)
	}
}
