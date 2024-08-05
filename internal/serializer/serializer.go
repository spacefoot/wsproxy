package serializer

import (
	"encoding/json"
	"errors"
)

type Serializer interface {
	Read(msg []byte) (any, error)
	Write(data any) ([]byte, error)
}

type TypeMeta struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func MarshalJSON(v any) ([]byte, error) {
	switch v := v.(type) {
	case *Weight:
		return marshalJSON("weight", v)
	case *Status:
		return marshalJSON("status", v)
	case *Zeroed:
		return marshalJSON("zeroed", v)
	default:
		return nil, errors.New("invalid message")
	}
}

func UnmarshalJSON(data []byte) (any, error) {
	var t TypeMeta
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	switch t.Type {
	case "weight":
		return unmarshalJSON(t.Data, &Weight{})
	case "status":
		return &RequestStatus{}, nil
	default:
		return nil, errors.New("invalid message")
	}
}

func marshalJSON(name string, data any) ([]byte, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return json.Marshal(TypeMeta{
		Type: name,
		Data: d,
	})
}

func unmarshalJSON(data []byte, v any) (any, error) {
	err := json.Unmarshal(data, v)
	return v, err
}
