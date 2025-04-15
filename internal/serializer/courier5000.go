package serializer

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const STABLE_DELAY = 2 * time.Second

type Lock int

const (
	LOCK_NONE Lock = iota
	LOCK_WEIGHT
	LOCK_UNSTABLE
)

type Courier5000 struct {
	isContinuous bool
	// Last read time. If 2 read <1s, continuous mode is enabled
	detectContinuousLastMsg time.Time
	// Last time weight changed
	continuousChangedTime time.Time
	// Lock last weight to avoid flooding
	lastWeight float64
	// Lock to send weight only once
	lock Lock
}

func (c *Courier5000) readContinuous(lines []string) (any, error) {
	if !(len(lines) == 2 || (len(lines) == 3 && lines[2] == "?")) {
		return nil, nil
	}

	weight, err := strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return nil, nil
	}

	if weight != c.lastWeight {
		c.lastWeight = weight
		c.continuousChangedTime = time.Now()

		if c.lock != LOCK_UNSTABLE {
			c.lock = LOCK_UNSTABLE
			return &Unstable{}, nil
		}
		return nil, nil
	}

	// Don't send weight if already sent
	if c.lock == LOCK_WEIGHT {
		return nil, nil
	}

	// Wait stable weight
	if time.Since(c.continuousChangedTime) < STABLE_DELAY || len(lines) == 3 /* unstable */ {
		return nil, nil
	}

	c.lock = LOCK_WEIGHT
	return &Weight{
		Weight: weight,
		Unit:   lines[1],
	}, nil
}

func (c *Courier5000) readStable(lines []string) (any, error) {
	// Continuous mode if has unstable weight
	if len(lines) == 3 && lines[2] == "?" {
		c.isContinuous = true
		return c.readContinuous(lines)
	}

	if len(lines) != 2 {
		return nil, nil
	}

	// Auto detect continuous mode
	if time.Since(c.detectContinuousLastMsg) < 500*time.Millisecond {
		c.isContinuous = true
		return c.readContinuous(lines)
	}
	c.detectContinuousLastMsg = time.Now()

	weight, err := strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return nil, nil
	}

	return &Weight{
		Weight: weight,
		Unit:   lines[1],
	}, nil
}

func (c *Courier5000) Read(msg []byte) (any, error) {
	if len(msg) == 0 {
		return nil, nil
	}

	lines := strings.Fields(string(msg))

	if c.isContinuous {
		return c.readContinuous(lines)
	}
	return c.readStable(lines)
}

func (*Courier5000) Write(data any) ([]byte, error) {
	switch data.(type) {
	case *Zero:
		return []byte("Z\r\n"), nil
	default:
		return nil, errors.New("unsupported data type")
	}
}
