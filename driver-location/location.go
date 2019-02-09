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
