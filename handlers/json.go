package handlers

import (
	"encoding/json"
	"io"
)

// ToJSON serializes the given interface into a
// string based JSON and write to writer
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// FromJSON deserializes an content of the reader
// into a given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
