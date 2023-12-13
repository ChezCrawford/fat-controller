package conductor

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
)

type eventChannel chan *Event

type Event struct {
	Name string
	Data interface{}
}

type StatusEvent struct {
	StatusRaw string
}

type PowerEvent struct {
	On bool
}

type FunctionState struct {
	Id string
	On bool
}

type Conductor struct {
	driver       DccDriver
	eventChannel eventChannel
	hornState    *FunctionState
	log          *slog.Logger
}

func NewConductor(log *slog.Logger, portName string, useSimDriver bool) *Conductor {
	eventChannel := make(eventChannel, 1024)

	var driver DccDriver
	if useSimDriver {
		driver = NewSimDriver()
	} else {
		driver = NewDccExDriver(log, portName, eventChannel)
	}

	hornState := &FunctionState{
		Id: "f3",
		On: false,
	}

	return &Conductor{
		driver:       driver,
		eventChannel: eventChannel,
		hornState:    hornState,
		log:          log,
	}
}

func (c *Conductor) Conduct(ctx context.Context) {
	c.log.InfoContext(ctx, "Hello.")

	err := c.driver.Start()
	if err != nil {
		c.log.WarnContext(ctx, "Unable to start driver", "error", err)
		listPorts()
		c.log.WarnContext(ctx, "Quitting.")
		os.Exit(1)
	}

	inputChannel := inputter()

	for {
		select {
		case data := <-c.eventChannel:
			c.log.InfoContext(ctx, "Received event", "event", *data)
		case input := <-inputChannel:
			c.log.InfoContext(ctx, "Received command", "input", input)
			if err = c.driver.SendRawCommand(input); err != nil {
				c.log.ErrorContext(ctx, "error", err)
				os.Exit(1)
			}
		case <-ctx.Done():
			c.log.InfoContext(ctx, "Hanging up the hat.")
			return
		}
	}
}

func (c *Conductor) LightsOn() {
	cmd := "<F 65 5 1>"
	c.driver.SendRawCommand(cmd)
}

func (c *Conductor) LightsOff() {
	cmd := "<F 65 5 0>"
	c.driver.SendRawCommand(cmd)
}

func (c *Conductor) Horn() {
	cmd := "<F 65 3 %d>"
	if c.hornState.On {
		c.driver.SendRawCommand(fmt.Sprintf(cmd, 0))
		c.hornState.On = false
	} else {
		c.driver.SendRawCommand(fmt.Sprintf(cmd, 1))
		c.hornState.On = true
	}
}

func inputter() chan string {
	c := make(chan string)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter command: ")
			text, _ := reader.ReadString('\n')
			c <- text
		}
	}()

	return c
}

func NewEvent(name string, data interface{}) *Event {
	return &Event{Name: name, Data: data}
}
