//go:build e2e_test

package cyclecmd_test

import (
	"os"
	"testing"
	"time"

	"github.com/RaphSku/cyclecmd"
	"github.com/stretchr/testify/assert"
)

func TestEventLifecycle(t *testing.T) {
	t.Parallel()

	defaultEventInformation := cyclecmd.EventInformation{
		EventName: "Default",
		Event:     &DefaultEvent{},
	}
	eventRegistry := cyclecmd.NewEventRegistry(defaultEventInformation)

	eventHistory := cyclecmd.NewEventHistory()

	consoleApp := cyclecmd.NewConsoleApp(
		"TestConsoleApp",
		"v0.1.0",
		"Test Console Application",
		eventRegistry,
		eventHistory,
	)
	consoleApp.ChangeToDebugMode()

	backspaceEventInformation := cyclecmd.EventInformation{
		EventName: "Backspace",
		Event:     &BackspaceEvent{},
	}
	err := eventRegistry.RegisterEvent("\b", backspaceEventInformation)
	assert.NoError(t, err)

	userInput := []byte("Hello W\borld")

	r, w, err := os.Pipe()
	assert.NoError(t, err)

	prevStdin := os.Stdin
	defer func() { os.Stdin = prevStdin }()
	os.Stdin = r

	go func() {
		for _, b := range userInput {
			w.Write([]byte{b})
			time.Sleep(10 * time.Millisecond)
		}
		w.Close()
	}()

	actOutput, err := captureStdOutput(func() {
		consoleApp.Start()
	})
	assert.NoError(t, err)
	expOutput := "Welcome to TestConsoleApp! Version: v0.1.0\r\nTest Console Application\rHello W\b \borld"
	assert.Equal(t, expOutput, actOutput)
}
