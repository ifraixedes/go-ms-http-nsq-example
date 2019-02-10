package drvloc

import (
	"encoding/json"
	"time"
)

// Location represents a geolocation position at certain time.
type Location struct {
	Lat float64   `json:"latitude"`
	Lng float64   `json:"longitude"`
	At  time.Time `json:"updated_at"`
}

// MarshalBinary satisfies the encoding.BinaryMarshaler interface marsahllin
// the location into a binary form.
func (l *Location) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

// UnmarshalBinary satisfies the encoding.BinaryUnmarshaler interface
// unmarshalling the location from a binary form.
func (l *Location) UnmarshalBinary(d []byte) error {
	return json.Unmarshal(d, &l)
}
