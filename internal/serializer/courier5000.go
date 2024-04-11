package serializer

import (
	"errors"
	"strconv"
	"strings"
)

type Courier5000 struct{}

func (Courier5000) Read(msg []byte) (any, error) {
	if len(msg) == 0 {
		return nil, nil
	}

	lines := strings.Fields(string(msg))

	if len(lines) != 2 {
		return nil, nil
	}

	weight, err := strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return nil, err
	}

	return &Weight{
		Weight: weight,
		Unit:   lines[1],
	}, nil
}

func (Courier5000) Write(data any) ([]byte, error) {
	return nil, errors.New("not implemented")
}
