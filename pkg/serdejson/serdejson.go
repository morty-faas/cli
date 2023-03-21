package serdejson

import (
	"encoding/json"
	"io"
)

// Serialize the value v in bytes
func Serialize(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Deserialize data into an instance of T
func Deserialize[T interface{}](data io.Reader) (*T, error) {
	t := new(T)
	if err := json.NewDecoder(data).Decode(t); err != nil {
		return nil, err
	}
	return t, nil
}
