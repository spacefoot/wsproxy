package serializer

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Courier5000 struct {
	isContinous bool
	// Lock last weight to avoid flooding
	lastWeight float64
	// Last read time. If 2 read <1s, continous mode is enabled
	lastRead time.Time
}

func (c *Courier5000) Read(msg []byte) (any, error) {
	if len(msg) == 0 {
		return nil, nil
	}

	lines := strings.Fields(string(msg))

	// Continous mode
	if len(lines) == 3 && lines[2] == "?" {
		c.isContinous = true
		c.lastWeight = 0 // Reset last weight if same stable weight
		return nil, nil
	}

	if len(lines) != 2 {
		return nil, nil
	}

	// Continous mode
	if !c.isContinous {
		now := time.Now()
		if now.Sub(c.lastRead) < time.Second {
			c.isContinous = true
		}
		c.lastRead = now
	}

	weight, err := strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return nil, err
	}

	// Continous mode, ignore invalid weight, ignore same weight
	if weight <= 0 || c.isContinous && weight == c.lastWeight {
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
