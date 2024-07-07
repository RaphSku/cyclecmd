//go:build e2e_test

package cyclecmd_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RaphSku/cyclecmd"
	"github.com/stretchr/testify/assert"
)

func TestEventLifecycle(t *testing.T) {
	t.Parallel()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	defaultEventInformation := cyclecmd.EventInformation{
		EventName: "Default",
		Event:     &DefaultEvent{},
	}
	eventRegistry := cyclecmd.NewEventRegistry(defaultEventInformation)

	eventHistory := cyclecmd.NewEventHistory()

	consoleApp := cyclecmd.NewConsoleApp(
		ctxTimeout,
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

	tmpFile, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(userInput)
	assert.NoError(t, err)

	_, err = tmpFile.Seek(0, 0)
	assert.NoError(t, err)

	prevStdin := os.Stdin
	defer func() { os.Stdin = prevStdin }()
	os.Stdin = tmpFile

	actOutput, err := captureStdOutput(func() {
		consoleApp.Start()
	})
	assert.NoError(t, err)
	expOutput := "Welcome to TestConsoleApp! Version: v0.1.0\r\nTest Console Application\rHello W\b \borld"
	assert.Equal(t, expOutput, actOutput)
}
