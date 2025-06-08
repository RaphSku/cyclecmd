# cyclecmd
## Description
cyclecmd helps you to write console applications based on events that are triggered by user input. It is an unopinionated library, thus you can define events in any way that you need.

## How to get it
You can get cyclecmd with the following command:
```bash
go get github.com/RaphSku/cyclecmd@latest
```

## How to use it
First, define a default event that will be triggered whenever user input doesn't match any other registered event. In this example, the default event will simply print the user's input (referred to as the "token").
```Go
type DefaultEvent struct{}

func (de *DefaultEvent) Handle(token string) (error, *cyclecmd.ControlEvent) {
	fmt.Print(token)
	return nil
}
```
**Note**, that all events have to comply with the following interface:
```Go
type Event interface {
    Handle(token string) (error, *ControlEvent)
}
```
Additionally, every `Handle` function returns an error and a `ControlEvent`, which can be used to instruct cyclecmd to trigger a specific event that may alter the application's flow. In the current version, however, the only available control event is for terminating the application (via the constant `CYCLE_TERMINATE`).
The default event is the only required event, as it serves to initialize the event registry and ensure there's always a fallback handler for unrecognized input:
```Go
defaultEventInformation := cyclecmd.EventInformation{
    EventName: "Default",
    Event:     &DefaultEvent{},
}
eventRegistry := cyclecmd.NewEventRegistry(defaultEventInformation)
```
For demonstration purposes, let us register another event, the backspace event that will print a backspace for us.
```Go
type BackspaceEvent struct{}

func (be *BackspaceEvent) Handle(token string) (error, *cyclecmd.ControlEvent) {
	fmt.Print("\b \b")
	return nil
}

backspaceEventInformation := cyclecmd.EventInformation{
    EventName: "Backspace",
    Event:     &BackspaceEvent{},
}
eventRegistry.RegisterEvent("\x7f", backspaceEventInformation)
```
Next, we need to initialize the event history, this is relatively simple:
```Go
eventHistory := cyclecmd.NewEventHistory()
```
And with the event registry and event history we can finally initialise our console app and let it run and handle our custom events.
```Go
consoleApp := cyclecmd.NewConsoleApp(
    "Test",
    "v0.1.0",
    "Example description...",
    eventRegistry,
    eventHistory,
)
consoleApp.SetLineDelimiter("\n\r>>> ", "\x7f")
consoleApp.Start()
```
Note that you can also set a line delimiter—for example, `"\n\r>>> "` in this case. If you want each new line to begin with `>>>`, be sure to include `"\n\r"` in the delimiter. This design is intentional, allowing you to customize the delimiter freely, even omitting new lines if needed. Finally, calling the `Start()` method begins the event loop.

## Example Projects
If you want to see a complete example on how to leverage cyclecmd, please have a look at the following projects that use cyclecmd: 
- [notewolfy](https://github.com/RaphSku/notewolfy)
