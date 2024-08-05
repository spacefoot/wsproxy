package serializer

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const STABLE_DELAY = 2 * time.Second
const MIN_WEIGHT_TO_SEND = 50

type Courier5000 struct {
	isContinuous bool
	// Lock last weight to avoid flooding
	lastWeight float64
	// Ensure scale is zeroed between measures
	lastStableWeight float64
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
		c.reset()
		return nil, nil
	}

	if len(lines) != 2 {
		return nil, nil
	}

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
		return nil, err
	}

	// Continuous mode, ignore invalid weight
	if weight < 0 {
		c.reset()
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

		if time.Since(c.lastChange) < STABLE_DELAY {
			return nil, nil
		}
	}

	c.sendLock = true

	if c.isContinuous {
		// Ensure scale is zeroed between measures
		if weight != 0 && c.lastStableWeight != 0 {
			return nil, nil
		}

		c.lastStableWeight = weight

		if weight == 0 {
			return &Zeroed{}, nil
		}
	}

	// skip below MIN_WEIGHT_TO_SEND threshold
	if lines[1] == "g" && weight < MIN_WEIGHT_TO_SEND || lines[1] == "kg" && weight*1000 < MIN_WEIGHT_TO_SEND {
		return nil, nil
	}

	return &Weight{
		Weight: weight,
		Unit:   lines[1],
	}, nil
}

func (*Courier5000) Write(data any) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (c *Courier5000) setLastWeight(weight float64) {
	c.lastWeight = weight
	c.lastChange = time.Now()
	c.sendLock = false
}

func (c *Courier5000) reset() {
	c.setLastWeight(0)
}
