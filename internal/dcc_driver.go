package conductor

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"go.bug.st/serial"
)

type DccExDriver struct {
	portName     string
	port         *serial.Port
	eventChannel eventChannel
}

func NewDccExDriver(portName string, eventChannel eventChannel) *DccExDriver {
	d := &DccExDriver{
		portName:     portName,
		port:         nil,
		eventChannel: eventChannel,
	}

	return d
}

func (d *DccExDriver) Start() (err error) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return err
	}

	d.port = &port

	d.reader()
	return nil
}

func listPorts() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
}

func (d *DccExDriver) reader() {
	go func() {
		scanner := bufio.NewScanner(*d.port)
		for scanner.Scan() {
			packet := scanner.Text()
			fmt.Printf("[RECEIVED] %+v\n", packet) // Println will add back the final '\n'
			event, _ := parsePacket(packet)
			if event != nil {
				d.eventChannel <- event
			}

		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()
}

func parsePacket(packet string) (*Event, error) {
	if !strings.HasPrefix(packet, "<") || !strings.HasSuffix(packet, ">") {
		return nil, fmt.Errorf("invalid packet")
	}

	packet = strings.Trim(packet, "<>")
	params := strings.Fields(packet)
	if len(params) < 1 {
		return nil, fmt.Errorf("not enough fields in the packet")
	}

	opCode, first := params[0][0], params[0][1:]
	params[0] = first

	switch opCode {
	case 'i':
		// Do not attempt to parse further for now
		statusRaw := strings.Join(params[1:], " ")
		data := NewEvent("status", StatusEvent{StatusRaw: statusRaw})
		return data, nil
	case 'p':
		switch params[0] {
		case "0":
			return NewEvent("power", PowerEvent{On: false}), nil
		case "1":
			return NewEvent("power", PowerEvent{On: true}), nil
		default:
			return nil, fmt.Errorf("unknown power event")
		}
	case '*':
		// debug events
	default:
		// unknown events
	}

	return nil, nil
}

func (d *DccExDriver) SendRawCommand(rawCommand string) (err error) {
	n, err := (*d.port).Write([]byte(rawCommand + "\r\n"))
	if err != nil {
		return err
	}

	fmt.Printf("[SENT] %v bytes\n", n)
	return nil
}
