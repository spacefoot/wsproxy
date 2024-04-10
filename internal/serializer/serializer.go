package serializer

import (
	"encoding/json"
)

type Serializer interface {
	Read(msg []byte) ([]byte, error)
}

type Weight struct {
	Weight float64 `json:"weight"`
	Unit   string  `json:"unit"`
}

func (w Weight) ToJSON() ([]byte, error) {
	return marshalJSON("weight", w)
}

func marshalJSON(name string, data any) ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Data any    `json:"data"`
	}{
		Type: name,
		Data: data,
	})
}

func UnmarshalJSON(data []byte) (any, error) {
	var t struct {
		Type string `json:"$type"`
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	switch t.Type {
	case "weight":
		return unmarshalJSON(data, &Weight{})
	default:
		return nil, nil
	}
}

func unmarshalJSON(data []byte, v any) (any, error) {
	err := json.Unmarshal(data, v)
	return v, err
}
