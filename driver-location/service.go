package drvloc

import (
	"context"
)

// Service specifies the operations that the Driver Location Service expose.
// All its methods can return the ErrAbortedCtx error code when the operation is
// aborted because context cancellation or deadline exceeded.
type Service interface {
	// SetLocation store the location of the driver with the associated id.
	// If the location was already set, it's ignore.
	//
	// Any implementation can returns one of the following error codes
	//
	//  * ErrUnexpectedErrorStore
	SetLocation(ctx context.Context, id uint64, l Location) error

	// LocationsForLastMinutes return the list of locations, sorted from the past
	// to the present, of driver with the associated id in the last minutes.
	//
	// Any implementation can returns one of the following error codes
	//
	//  * ErrInvalidDataFormatStore
	//
	//	* ErrNotFoundDriver
	//
	//  * ErrUnexpectedErrorStore
	LocationsForLastMinutes(ctx context.Context, id uint64, minutes uint16) ([]Location, error)
}
