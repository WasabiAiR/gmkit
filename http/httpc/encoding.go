package httpc

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
)

type (
	// EncodeFn is an encoder func.
	EncodeFn func(interface{}) (io.Reader, error)

	// DecodeFn is a decoder func.
	DecodeFn func(r io.Reader) error
)

// JSONEncode sets the client's encodeFn to a json encoder.
func JSONEncode() EncodeFn {
	return func(v interface{}) (io.Reader, error) {
		var buf bytes.Buffer
		return &buf, json.NewEncoder(&buf).Encode(v)
	}
}

// JSONDecode sets the client's decodeFn to a json decoder.
func JSONDecode(v interface{}) DecodeFn {
	return func(r io.Reader) error {
		return json.NewDecoder(r).Decode(v)
	}
}

// GobEncode sets the client's encodeFn to a gob encoder.
func GobEncode() EncodeFn {
	return func(v interface{}) (io.Reader, error) {
		var buf bytes.Buffer
		return &buf, gob.NewEncoder(&buf).Encode(v)
	}
}

// GobDecode sets the client's decodeFn to a gob decoder.
func GobDecode(v interface{}) DecodeFn {
	return func(r io.Reader) error {
		return gob.NewDecoder(r).Decode(v)
	}
}
