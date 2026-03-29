package dto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// JSONValue stores a dynamic JSON value while preserving number fidelity via json.Number.
type JSONValue struct {
	value any
	set   bool
}

func (j *JSONValue) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	var v any
	if err := dec.Decode(&v); err != nil {
		return fmt.Errorf("invalid JSON value: %w", err)
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("invalid JSON value: trailing data")
	}

	j.value = v
	j.set = true
	return nil
}

func (j JSONValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.value)
}

func (j JSONValue) Value() any {
	return j.value
}

func (j JSONValue) IsSet() bool {
	return j.set
}

func NewJSONValue(v any) JSONValue {
	return JSONValue{value: v, set: true}
}

func NewJSONValueFromRaw(data []byte) JSONValue {
	if len(data) == 0 {
		return JSONValue{value: nil, set: true}
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	var v any
	if err := dec.Decode(&v); err != nil {
		return JSONValue{value: nil, set: true}
	}

	return JSONValue{value: v, set: true}
}
