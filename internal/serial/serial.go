package serial

import (
	"bufio"
	"bytes"
	"errors"
	"log/slog"
	"math"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type ISerial interface {
	Run()
	RequestStatus()
}

type Status struct {
	Open bool
}

type Serial struct {
	port   serial.Port
	opened bool

	read   chan<- []byte
	write  <-chan []byte
	status chan<- Status
}

func NewSerial(read chan<- []byte, write <-chan []byte, status chan<- Status) *Serial {
	return &Serial{
		read:   read,
		write:  write,
		status: status,
	}
}

func (s *Serial) Open() error {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return err
	}

	if len(ports) == 0 {
		return errors.New("no serial ports found")
	}

	if len(ports) > 1 {
		return errors.New("multiple serial ports found. not supported")
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(ports[0].Name, mode)
	if err != nil {
		return err
	}

	slog.Info("serial port opened", "port", ports[0].Name)
	s.port = port
	s.opened = true
	s.RequestStatus()
	return nil
}

const (
	// BACKOFF_MAX_DURATION is the maximum duration to wait between retries
	BACKOFF_MAX_DURATION = 5 * time.Second
	// BACKOFF_SCALE is the factor to increase the backoff duration each retry
	BACKOFF_SCALE = 2
	// BACKOFF_SCALE_COUNT is the number of retries to increase the backoff duration
	BACKOFF_SCALE_COUNT = 5
)

func (s *Serial) OpenWithBackoff() {
	backoff := time.Duration(float64(BACKOFF_MAX_DURATION) / math.Pow(BACKOFF_SCALE, BACKOFF_SCALE_COUNT))

	for {
		err := s.Open()
		if err == nil {
			return
		}

		slog.Warn("serial open error", "err", err)
		slog.Info("retrying in", "duration", backoff)

		time.Sleep(backoff)
		if backoff < BACKOFF_MAX_DURATION {
			backoff *= BACKOFF_SCALE
		}
	}
}

func (s *Serial) reader() error {
	defer s.Close()

	reader := bufio.NewReader(s.port)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		if len(line) == 0 {
			break
		}

		line = bytes.TrimRight(line, "\r\n")
		slog.Debug("read from serial", "data", string(line))
		s.read <- line
	}

	return nil
}

func (s *Serial) readLoop() {
	for {
		s.OpenWithBackoff()

		err := s.reader()
		if err != nil {
			slog.Error("serial reader error", "err", err)
		}
	}
}

func (s *Serial) Run() {
	go s.readLoop()

	for msg := range s.write {
		if !s.opened {
			slog.Warn("serial port not open")
			continue
		}

		slog.Debug("write to serial", "data", string(msg))
		if _, err := s.port.Write(msg); err != nil {
			slog.Error("serial write error", "err", err)
		}
	}
}

func (s *Serial) Close() {
	s.port.Close()
	s.opened = false
	s.RequestStatus()
}

func (s *Serial) RequestStatus() {
	s.status <- Status{
		Open: s.opened,
	}
}
