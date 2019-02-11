package zmbdrv

// Rules are the set of rules which are used to figure out if a driver is a
// zombie or not.
type Rules struct {
	// MinDistance is the minimum distance that a driver have to driver for not
	// being considered to be a zombie. The distance is expressed in meters.
	MinDistance uint64
	// LastMinutes is the number of minutes, to calculate the distance that a
	// driver has driven. The time is calculated based on the current time when
	// the question if it's a zombie or not is made.
	LastMinutes uint16
}

// IsValid checks if r is valid, if it's returns true, otherwise false.
func (r Rules) IsValid() bool {
	if r.MinDistance == 0 || r.LastMinutes == 0 {
		return false
	}

	return true
}
