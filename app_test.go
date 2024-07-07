//go:build unit_test

package cyclecmd_test

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/RaphSku/cyclecmd"
	"github.com/stretchr/testify/assert"
)

func TestConsoleAppCreation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	eventRegistry := cyclecmd.NewEventRegistry(setupDefaultEventInformation())
	eventHistory := cyclecmd.NewEventHistory()

	actConsoleApp := cyclecmd.NewConsoleApp(
		ctx,
		"test",
		"0.1.0",
		"This is a test console application",
		eventRegistry,
		eventHistory,
	)
	expConsoleApp := &cyclecmd.ConsoleApp{
		Name:        "test",
		Version:     "0.1.0",
		Description: "This is a test console application",
	}

	assert.Equal(t, expConsoleApp.Name, actConsoleApp.Name)
	assert.Equal(t, expConsoleApp.Version, actConsoleApp.Version)
	assert.Equal(t, expConsoleApp.Description, actConsoleApp.Description)
}

func TestChangeToDebugMode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	eventRegistry := cyclecmd.NewEventRegistry(setupDefaultEventInformation())
	eventHistory := cyclecmd.NewEventHistory()

	consoleApp := cyclecmd.NewConsoleApp(
		ctx,
		"test",
		"0.1.0",
		"This is a test console application",
		eventRegistry,
		eventHistory,
	)

	expOutput := "Attention! You have enabled debug mode (Level: debug)! Turn off if running in production!\r\n"
	actOutput, err := captureStdOutput(func() {
		consoleApp.ChangeToDebugMode()
	})
	assert.NoError(t, err)
	assert.Equal(t, expOutput, actOutput)
}

func TestSetLineDelimiter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	eventRegistry := cyclecmd.NewEventRegistry(setupDefaultEventInformation())

	expEventTrigger := "z"
	delimiterEventInformation := cyclecmd.EventInformation{
		EventName: "Z",
		Event:     &TestEvent{},
	}
	err := eventRegistry.RegisterEvent(expEventTrigger, delimiterEventInformation)
	assert.NoError(t, err)
	eventHistory := cyclecmd.NewEventHistory()

	consoleApp := cyclecmd.NewConsoleApp(
		ctx,
		"test",
		"0.1.0",
		"This is a test console application",
		eventRegistry,
		eventHistory,
	)

	expDelimiter := "<<< "
	consoleApp.SetLineDelimiter(expDelimiter, expEventTrigger)

	assert.Equal(t, expDelimiter, consoleApp.Delimiter)
	assert.Equal(t, expEventTrigger, consoleApp.DelimiterEventTrigger)
}

func TestSetLineDelimiterNotInRegistry(t *testing.T) {
	t.Parallel()

	if os.Getenv("UT_SetLineDelimiterNotInRegistry") == "1" {
		ctx := context.Background()
		eventRegistry := cyclecmd.NewEventRegistry(setupDefaultEventInformation())
		eventHistory := cyclecmd.NewEventHistory()

		consoleApp := cyclecmd.NewConsoleApp(
			ctx,
			"test",
			"0.1.0",
			"This is a test console application",
			eventRegistry,
			eventHistory,
		)

		delimiter := "<<< "
		eventTrigger := "z"
		consoleApp.SetLineDelimiter(delimiter, eventTrigger)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSetLineDelimiterNotInRegistry")
	cmd.Env = append(os.Environ(), "UT_SetLineDelimiterNotInRegistry=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process returned err: %v, want exit code 1", err)
}
