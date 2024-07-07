package cyclecmd

import "fmt"

// EventHistory records events and offers behavior to manipulate the history and to
// read events from history.
type EventHistory struct {
	// entries is a sequence of events and tokens that triggered those events
	entries []EventHistoryEntry
}

// NewEventHistory initialises an event history instance that can be used to record past events.
//
// Returns:
//   - `*EventHistory` : Returns an instance of EventHistory
func NewEventHistory() *EventHistory {
	return &EventHistory{}
}

// Len returns the number of events that has been recorded.
//
// Returns:
//   - `int` : Number of events
func (eh *EventHistory) Len() int {
	return len(eh.entries)
}

// AddEvent will add an event to the history.
//
// Parameters:
//   - `eventEntry` : Entry that will be recorded and contains all information related to an event.
func (eh *EventHistory) AddEvent(eventEntry EventHistoryEntry) {
	eh.entries = append(eh.entries, eventEntry)
}

// RetrieveEventEntryByIndex will return the event entry at index position
//
// Parameters:
//   - `index` : Position in the event history that you want to access
//
// Returns:
//   - `EventHistoryEntry` : The event history entry at index position
//   - `error` : Returns an error when there is no event history entry at position index
func (eh *EventHistory) RetrieveEventEntryByIndex(index int) (EventHistoryEntry, error) {
	if index < 0 || index >= eh.Len() {
		return EventHistoryEntry{}, fmt.Errorf("index %v error, index is either smaller than 0 or larger than the length of the event history", index)
	}
	return eh.entries[index], nil
}

// RemoveNthEventFromHistory removes the nth event from the history. If the nth element does not exist,
// nothing will happen.
//
// Parameters:
//   - `n` : nth event that should be removed from history
func (eh *EventHistory) RemoveNthEventFromHistory(n int) {
	eventHistoryLength := eh.Len()
	if n > eventHistoryLength || n < 0 {
		return
	}
	eh.entries = append(eh.entries[:n], eh.entries[n+1:]...)
}

// PrintLastEventHistoryEntries will print information related to the last n events that
// were recorded.
//
// Parameters:
//   - `n` : Number of events
func (eh *EventHistory) PrintLastEventHistoryEntries(n int) {
	count := 0
	for i := eh.Len() - 1; i >= 0; i-- {
		if count == n {
			break
		}
		fmt.Printf("Event Name: %s, Token: %s\r\n", eh.entries[i].EventName, eh.entries[i].Token)
		count += 1
	}
}

// GetLastEventsFromHistoryToEventReference will return all event names that followed after a specific event happened.
//
// Parameters:
//   - `eventName` : The name of the event that serves as a reference
//
// Returns:
//   - `[]string` : Returns a series of event names that happened after the reference event
func (eh *EventHistory) GetLastEventsFromHistoryToEventReference(eventName string) []string {
	var eventNames []string
	for i := eh.Len() - 1; i >= 0; i-- {
		eventNameFromHistory := eh.entries[i].EventName
		if eventNameFromHistory == eventName {
			break
		}
		eventNames = append(eventNames, eventNameFromHistory)
	}
	// We need to reverse the array since we appended the elements from the back to the front
	return reverseArray(eventNames)
}

// MostRecentSpliceEventsOfHistory will return a range of events that occured between some event that
// is specified by eventName.
//
// Parameters:
//   - `eventName` : Name of the reference event
//
// Returns:
//   - `[]EventHistoryEntry` : Sequence of event history entries
func (eh *EventHistory) MostRecentSpliceEventsOfHistory(eventName string) []EventHistoryEntry {
	var splicedEvents []EventHistoryEntry
	lastIndex := eh.Len() - 1
	foundStart := false
	foundEnd := false
	for i := lastIndex; i >= 0; i-- {
		if foundStart {
			splicedEvents = append(splicedEvents, eh.entries[i])
		}
		if eventName == eh.entries[i].EventName && !foundStart {
			foundStart = true
			continue
		}
		if eventName == eh.entries[i].EventName && foundStart {
			foundEnd = true
			splicedEvents = splicedEvents[:len(splicedEvents)-1]
		}
		if foundEnd {
			break
		}
	}
	// We need to reverse the array since we appended the elements from the back to the front
	return reverseArray(splicedEvents)
}

// Generic implementation to reverse a slice.
//
// Parameters:
//   - `arr` : Slice that should be reversed
//
// Returns:
//   - `[]T` : Reversed slice
func reverseArray[T any](arr []T) []T {
	n := len(arr)
	for i := 0; i < n/2; i++ {
		arr[i], arr[n-i-1] = arr[n-i-1], arr[i]
	}
	return arr
}
