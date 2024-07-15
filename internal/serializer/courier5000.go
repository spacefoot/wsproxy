package serializer

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Courier5000 struct {
	isContinuous bool
	// Lock last weight to avoid flooding
	lastWeight float64
	// Last read time. If 2 read <1s, continuous mode is enabled
	lastRead time.Time
}

func (c *Courier5000) Read(msg []byte) (any, error) {
	if len(msg) == 0 {
		return nil, nil
	}

	lines := strings.Fields(string(msg))

	// Continuous mode
	if len(lines) == 3 && lines[2] == "?" {
		c.isContinuous = true
		c.lastWeight = 0 // Reset last weight if same stable weight
		return nil, nil
	}

	if len(lines) != 2 {
		return nil, nil
	}

	// Continuous mode
	if !c.isContinuous {
		now := time.Now()
		if now.Sub(c.lastRead) < time.Second {
			c.isContinuous = true
		}
		c.lastRead = now
	}

	weight, err := strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return nil, err
	}

	// Continuous mode, ignore invalid weight
	if weight <= 0 {
		c.lastWeight = 0
		return nil, nil
	}

	// Continuous mode, ignore same weight
	if c.isContinuous && weight == c.lastWeight {
		return nil, nil
	}

	c.lastWeight = weight
	return &Weight{
		Weight: weight,
		Unit:   lines[1],
	}, nil
}

func (*Courier5000) Write(data any) ([]byte, error) {
	return nil, errors.New("not implemented")
}
