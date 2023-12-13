package conductor

import (
	"bufio"
	"fmt"
	"log/slog"
	"strings"

	"go.bug.st/serial"
)

type DccDriver interface {
	Start() (err error)
	SendRawCommand(rawCommand string) (err error)
}

type DccExDriver struct {
	eventChannel eventChannel
	log          *slog.Logger
	portName     string
	port         *serial.Port
}

func NewDccExDriver(log *slog.Logger, portName string, eventChannel eventChannel) DccDriver {
	d := &DccExDriver{
		eventChannel: eventChannel,
		log:          log,
		portName:     portName,
		port:         nil,
	}

	return d
}

func (d *DccExDriver) Start() (err error) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err := serial.Open(d.portName, mode)
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
		slog.Error("Unable to get ports list", "error", err)
		return
	}

	if len(ports) == 0 {
		slog.Error("No serial ports found!")
		return
	}

	for _, port := range ports {
		slog.Info("", "port", port)
	}
}

func (d *DccExDriver) reader() {
	go func() {
		scanner := bufio.NewScanner(*d.port)
		for scanner.Scan() {
			packet := scanner.Text()
			slog.Debug("received", "packet", packet)
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

	d.log.Info(fmt.Sprintf("Sent %d bytes", n))
	return nil
}
