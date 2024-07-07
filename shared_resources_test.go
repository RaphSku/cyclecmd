package cyclecmd_test

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/RaphSku/cyclecmd"
)

type DefaultEvent struct{}

func (de *DefaultEvent) Handle(token string) error {
	fmt.Print(token)
	return nil
}

func setupDefaultEventInformation() cyclecmd.EventInformation {
	return cyclecmd.EventInformation{
		EventName: "Default",
		Event:     &DefaultEvent{},
	}
}

type TestEvent struct{}

func (te *TestEvent) Handle(token string) error {
	fmt.Print("Testing this event")
	return nil
}

type BackspaceEvent struct{}

func (be *BackspaceEvent) Handle(token string) error {
	fmt.Print("\b \b")
	return nil
}

func captureStdOutput(f func()) (string, error) {
	originalStdOut := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	outputC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputC <- buf.String()
	}()

	f()
	w.Close()

	os.Stdout = originalStdOut
	out := <-outputC

	return out, nil
}
