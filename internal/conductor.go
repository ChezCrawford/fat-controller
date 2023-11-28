package conductor

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
)

var portName = "/dev/tty.usbmodem144301"

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
	driver       *DccExDriver
	eventChannel eventChannel
	hornState    *FunctionState
}

func NewConductor() *Conductor {
	eventChannel := make(eventChannel, 1024)
	driver := NewDccExDriver(portName, eventChannel)

	hornState := &FunctionState{
		Id: "f3",
		On: false,
	}

	return &Conductor{
		driver:       driver,
		eventChannel: eventChannel,
		hornState:    hornState,
	}
}

func (c *Conductor) Conduct(context context.Context) {
	fmt.Println("Hello.")

	inputChannel := inputter()
	err := c.driver.Start()
	if err != nil {
		fmt.Printf("Unable to start driver. Error: %+v\n", err)
		listPorts()
		log.Fatalf("\nQutting...\n")
	}

	for {
		select {
		case data := <-c.eventChannel:
			fmt.Printf("Event: %+v\n", *data)
		case input := <-inputChannel:
			if err = c.driver.SendRawCommand(input); err != nil {
				log.Fatal(err)
			}
		case <-context.Done():
			fmt.Printf("\nHanging up the hat...\n")
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
			fmt.Printf("Received command: %s\n", text)
			c <- text
		}
	}()

	return c
}

func NewEvent(name string, data interface{}) *Event {
	return &Event{Name: name, Data: data}
}
