package serial

import (
	"errors"
	"log"
	"time"

	"github.com/tarm/serial"
)

var s *serial.Port

func init() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600}
	tmp, err := serial.OpenPort(c)
	chkErr(err)
	s = tmp
	log.Println("s is ", s)
	go rcv()
}

// Send sends data to serial
func Send(data []byte) (err error) {
	log.Printf("sending data % x\n", data)
	if s == nil {
		log.Println("no connection to serial")
		return
	}
	n, err := s.Write(data)

	if err != nil {
		return
	}
	if n != len(data) {
		return errors.New("incomplete  send")
	}
	return
}

func rcv() {
	buf := make([]byte, 1024)
	for {
		if s == nil {
			break
		}
		n, err := s.Read(buf)
		if err != nil {
			log.Println("error in received data", err)
			time.Sleep(time.Second)
			continue
		}
		log.Printf("received data % x\n", buf[:n])

		processStatusData(buf)
	}
}

// Status is status of light
type Status struct {
	InputPower          byte    `json:"inputPower"`
	Dim                 byte    `json:"dim"`
	OutputVoltage       byte    `json:"outputVoltage"`
	OutputCurrent       float32 `json:"outputCurrent"`
	InternalTemperature byte    `json:"internalTemperatue"`
	Timestamp           string  `json:"timestamp"`
}

// ChanStatus is channel of status of lights struct
var ChanStatus = make(chan *Status, 10)

func processStatusData(data []byte) {
	if len(data) < 10 || data[0] != 0X50 {
		log.Println("invalid status packet")
	}

	status := new(Status)
	status.InputPower = data[2]
	ChanStatus <- status
}

func chkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
