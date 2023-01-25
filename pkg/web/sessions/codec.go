package sessions

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"time"
)

// GobCodec is an implementation of Codec using gob encoding
type GobCodec struct{}

// Encode converts a session deadline and data into a byte slice
func (c GobCodec) Encode(deadline time.Time, data map[string]any) ([]byte, error) {
	temp := struct {
		Deadline time.Time
		Data     map[string]any
	}{
		Deadline: deadline,
		Data:     data,
	}
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(&temp)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode converts a byte slice into a session deadline, and data
func (c GobCodec) Decode(b []byte) (time.Time, map[string]any, error) {
	temp := struct {
		Deadline time.Time
		Data     map[string]any
	}{}
	r := bytes.NewReader(b)
	err := gob.NewDecoder(r).Decode(&temp)
	if err != nil {
		return time.Time{}, nil, err
	}
	return temp.Deadline, temp.Data, nil
}

// JSONCodec is an implementation of Codec using json encoding
type JSONCodec struct{}

// Encode converts a session deadline and values into a byte slice
func (c JSONCodec) Encode(deadline time.Time, data map[string]any) ([]byte, error) {
	temp := struct {
		Deadline time.Time
		Data     map[string]any
	}{
		Deadline: deadline,
		Data:     data,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&temp)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode converts a byte slice into a session deadline, and data
func (c JSONCodec) Decode(b []byte) (time.Time, map[string]any, error) {
	temp := struct {
		Deadline time.Time
		Data     map[string]any
	}{}
	r := bytes.NewReader(b)
	err := json.NewDecoder(r).Decode(&temp)
	if err != nil {
		return time.Time{}, nil, err
	}
	return temp.Deadline, temp.Data, nil
}
