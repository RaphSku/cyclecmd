// Package cyclecmd is an unopinionated library for building your own console applications.
package cyclecmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
)

// ConsoleApp handles all the events in an event loop and serves as an entry point for your console.
type ConsoleApp struct {
	ctx    context.Context
	logger *zap.Logger
	// The Event Registry is the source of truth for all custom events that
	// were registered.
	eventRegistry *EventRegistry
	// The Event History is decoupled from the Event Registry and records
	// all events that were processed.
	eventHistory *EventHistory

	// Name of the console application
	Name string
	// Version of the console application
	Version string
	// Description of the console application. Should be relatively short.
	Description string
	// Delimiter can be used to visibly separate lines in your console application.
	// But the usage is flexible.
	Delimiter string
	// DelimiterEventTrigger defines when a Delimiter will be printed.
	DelimiterEventTrigger string
	// Should be only used when you want to pause your console application for some reason.
	ControlC chan bool
}

// initLogger provides a Zap logger for structured logging.
// Primary usage of the logger is for debugging, in production there will be no logs.
//
// Parameters:
//   - `debug` : Whether the logger should be configured to DebugLevel or FatalLevel
//
// Returns:
//   - `*zap.Logger` : Configured and built Zap logger instance
func initLogger(debug bool) *zap.Logger {
	var config zap.Config

	if debug {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	}

	logger, err := config.Build()
	if err != nil {
		fmt.Print("logger creation failed! Configuration could not be built! error:", err)
		os.Exit(1)
	}
	return logger
}

// NewConsoleApp initializes a new Console App instance that will be the
// main entry point to your console application or to be more specific,
// to the event loop that will handle the custom events.
//
// Parameters:
//   - `ctx` : Context that will be propagated to the event loop
//   - `name` : Name of the Console App
//   - `version` : Version of the Console App
//   - `description` : Description of the Console App. Should be relatively short
//   - `eventRegistry` : The single source of truth for custom events that were registered
//   - `eventHistory` : Records events that were processed
//
// Returns:
//   - `*ConsoleApp` : An instance of a Console Application
func NewConsoleApp(
	ctx context.Context,
	name string,
	version string,
	description string,
	eventRegistry *EventRegistry,
	eventHistory *EventHistory,
) *ConsoleApp {
	logger := initLogger(false)

	consoleApp := &ConsoleApp{
		ctx:           ctx,
		logger:        logger,
		eventRegistry: eventRegistry,
		eventHistory:  eventHistory,
		Name:          name,
		Version:       version,
		Description:   description,
	}
	consoleApp.ControlC = make(chan bool)

	return consoleApp
}

// ChangeToDebugMode allows you to switch to debug mode for logging.
// Don't forget to switch off debug mode once you want to ship your console app
// to production since all logs will be otherwise shown.
func (ca *ConsoleApp) ChangeToDebugMode() {
	ca.logger = initLogger(true)
	fmt.Printf("Attention! You have enabled debug mode (Level: %v)! Turn off if running in production!\r\n", ca.logger.Level())
	ca.logger.Debug("Logger is now set to debug level", zap.String("func", "ChangeToDebugMode"))
}

// SetLineDelimiter allows the User to define a custom delimiter that will be printed
// after each event that is defined by eventTrigger.
//
// Parameters:
//   - `delimiter` : Delimiter should be fairly short.
//   - `eventTrigger` : Delimiter will be printed after each event that is triggered by eventTrigger
func (ca *ConsoleApp) SetLineDelimiter(delimiter string, eventTrigger string) {
	ca.Delimiter = delimiter
	_, ok := ca.eventRegistry.registry[eventTrigger]
	if !ok {
		fmt.Print("line delimiter event needs to be available in the event registry!")
		os.Exit(1)
	}
	ca.DelimiterEventTrigger = eventTrigger
}

// Start will save the terminal state, handle terminating signals and kick off the event loop. Note, events are recorded
// in the event history before the event handling happens. They are recorded as they occur.
func (ca *ConsoleApp) Start() {
	ca.logger.Debug("Saving current terminal (if Stdin is a terminal) state before entering the event loop", zap.String("func", "Start"))
	prevState := ca.saveTerminalState()
	ca.logger.Debug("Current terminal state has been saved successfully", zap.String("func", "Start"))

	ca.logger.Debug("Starting the goroutine to handle terminating signals", zap.String("func", "Start"))
	ca.handleTerminatingSignals(prevState)
	ca.logger.Debug("Terminating signal handler started successfully", zap.String("func", "Start"))

	ca.logger.Debug("Will enter event loop now", zap.String("func", "Start"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ca.eventLoop(ctx, prevState)
}

// saveTerminalState will save the state of the terminal, if there is no terminal available, no state will be saved.
//
// Returns:
//   - `*term.State` : Returns the reference of the current terminal state
func (ca *ConsoleApp) saveTerminalState() *term.State {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		ca.logger.Debug("Detected that standard input is not a terminal", zap.String("func", "saveTerminalState"))
		return nil
	}

	terminalState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Print("terminal state could not be saved! error:", err)
		os.Exit(1)
	}

	return terminalState
}

// handleTerminatingSignals handels SIGINT and SIGTERM signals and will restore the original terminal state and will
// exit successfully from the application.
func (ca *ConsoleApp) handleTerminatingSignals(terminalState *term.State) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		term.Restore(int(os.Stdin.Fd()), terminalState)
		os.Exit(0)
	}()
}

// eventLoop is a long running process that will capture Stdin and handle incoming events,
// as well as record the event history.
//
// Parameters:
//   - `ctx` : Context that can be used to cancel the application
//   - `prevState`: The previous terminal state that will be restored after the event loop concludes
func (ca *ConsoleApp) eventLoop(ctx context.Context, prevState *term.State) {
	if prevState != nil {
		defer term.Restore(int(os.Stdin.Fd()), prevState)
	}
	defer ca.logger.Sync()
	fmt.Printf("Welcome to %s! Version: %s\r\n%s\r", ca.Name, ca.Version, ca.Description)
	fmt.Printf("%s", ca.Delimiter)
	for {
		select {
		case <-ctx.Done():
			ca.logger.Debug("Context is done", zap.String("func", "eventLoop"))
			return
		case active := <-ca.ControlC:
			if !active {
				ca.logger.Debug("Context has been paused", zap.String("func", "eventLoop"))
				select {
				case <-ca.ControlC:
					ca.logger.Debug("Context has been resumed", zap.String("func", "eventLoop"))
				case <-ctx.Done():
					ca.logger.Debug("Context is done", zap.String("func", "eventLoop"))
					return
				}
			}
		default:
			b := make([]byte, 1)
			n, err := os.Stdin.Read(b)
			if n == 0 && prevState == nil {
				ca.logger.Debug("EOF found", zap.String("func", "eventLoop"))
				return
			}
			if err != nil {
				ca.logger.Debug("Could not read from Stdin", zap.Error(err), zap.String("func", "eventLoop"))
				return
			}
			token := string(b[0])
			ca.logger.Debug("Token captured", zap.String("Token", token), zap.String("func", "eventLoop"))
			eventInformation, err := ca.eventRegistry.GetMatchingEventInformation(token)
			if err != nil {
				ca.logger.Debug("Did not find a matching event", zap.Error(err), zap.String("func", "eventLoop"))
				os.Exit(1)
			}
			eventHistoryEntry := EventHistoryEntry{
				Token:     token,
				EventName: eventInformation.EventName,
				Event:     eventInformation.Event,
			}
			ca.eventHistory.AddEvent(eventHistoryEntry)
			lengthOfHistoryString := strconv.Itoa(ca.eventHistory.Len())
			ca.logger.Debug("Event History Length", zap.String("Length", lengthOfHistoryString), zap.String("func", "eventLoop"))
			err = eventInformation.Event.Handle(token)
			if err != nil {
				ca.logger.Debug("Event handling failed", zap.Error(err), zap.String("func", "eventLoop"))
				os.Exit(1)
			}
			if token == ca.DelimiterEventTrigger {
				fmt.Print(ca.Delimiter)
			}
		}
	}
}
