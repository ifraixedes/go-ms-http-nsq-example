package drvloc

import (
	"context"
)

//go:generate mockgen -package mock -destination mock/service.go -mock_names=Service=Service github.com/ifraixedes/go-ms-http-nsq-example/driver-location Service

// Service specifies the operations that the Driver Location Service exposes.
//
// All its methods can return the ErrAbortedCtx error code when the operation is
// aborted because context cancellation or deadline exceeded. Furthermore
// ErrUnexpectedErrorSystem can also be returned by any of the methods when
// the implementation encounter an error which cannot classified to a more
// precise error code.
//
// When the service implementation is a client, any method can return a
// ErrUnexpectedErrorTransport, ErrUnexpectedErrorSerializer and
// ErrUnexpectedErrorDeserializer
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
