# cyclecmd
## Description
cyclecmd helps you to write console applications based on events that are triggered by the user input. It is an unopinionated library, thus you can define events in any way that you need.

## How to get it
You can get cyclecmd with the following command:
```bash
go get github.com/RaphSku/cyclecmd@latest
```

## How to use it
First of all, you have to define a default event that will be triggered by any user input that does not trigger any other registered event. For this example, let's print the user input (that we call token) on our default event.
```Go
type DefaultEvent struct{}

func (de *DefaultEvent) Handle(token string) error {
	fmt.Print(token)
	return nil
}
```
Note, that all events have to comply with the following interface:
```Go
type Event interface {
    Handle(token string) error
}
```
The default event is the only event that we are required to create, because we can initialize the event registry with this event:
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

func (be *BackspaceEvent) Handle(token string) error {
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
    context.Background(),
    "Test",
    "v0.1.0",
    "Example description...",
    eventRegistry,
    eventHistory,
)
consoleApp.SetLineDelimiter("\n\r>>> ", "\x7f")
consoleApp.Start()
```
Note, that we can also set a line delimiter that will print in our case "\n\r>>> ". Be ware, that you have to add "\n\r" to the delimiter if you want to print a new line where each line begins with >>>. This is intentional, such that you can define any delimiter that you want, even ones with no new lines if that is what you require. The `Start()` method will kick off the event loop. 

## Example Projects
If you want to see a complete example on how to leverage cyclecmd, please have a look at the following projects that use cyclecmd: 

