package serializer

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const STABLE_DELAY = 2 * time.Second

type Courier5000 struct {
	isContinuous bool
	// Lock last weight to avoid flooding
	lastWeight float64
	// Last time weight changed
	lastChange time.Time
	// Lock to send weight only once
	sendLock bool
	// Last read time. If 2 read <1s, continuous mode is enabled
	lastRead time.Time
}

func (c *Courier5000) Read(msg []byte) (any, error) {
	if len(msg) == 0 {
		return nil, nil
	}

	lines := strings.Fields(string(msg))

	// Auto detect continuous mode
	if len(lines) == 3 && lines[2] == "?" {
		c.isContinuous = true
	} else if len(lines) != 2 {
		return nil, nil
	}
	unstable := len(lines) == 3

	// Auto detect continuous mode
	if !c.isContinuous {
		now := time.Now()
		if now.Sub(c.lastRead) < time.Second {
			c.isContinuous = true
		}
		c.lastRead = now
	}

	weight, err := strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return nil, nil
	}

	if c.isContinuous {
		if weight != c.lastWeight {
			c.setLastWeight(weight)
			return nil, nil
		}

		if c.sendLock {
			return nil, nil
		}

		if time.Since(c.lastChange) < STABLE_DELAY || unstable {
			return nil, nil
		}
	}

	c.sendLock = true

	return &Weight{
		Weight: weight,
		Unit:   lines[1],
	}, nil
}

func (*Courier5000) Write(data any) ([]byte, error) {
	switch data.(type) {
	case *Zero:
		return []byte("Z\r\n"), nil
	default:
		return nil, errors.New("unsupported data type")
	}
}

func (c *Courier5000) setLastWeight(weight float64) {
	c.lastWeight = weight
	c.lastChange = time.Now()
	c.sendLock = false
}
