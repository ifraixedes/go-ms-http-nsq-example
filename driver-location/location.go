package drvloc

import "time"

// Location represents a geolocation position at certain time.
type Location struct {
	Lat float64   `json:"latitude"`
	Lng float64   `json:"longitude"`
	At  time.Time `json:"updated_at"`
}
