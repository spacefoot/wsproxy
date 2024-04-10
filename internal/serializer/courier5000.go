package serializer

import (
	"strconv"
	"strings"
)

type Courier5000 struct{}

func (Courier5000) Read(msg []byte) ([]byte, error) {
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

	return Weight{
		Weight: weight,
		Unit:   lines[1],
	}.ToJSON()
}
