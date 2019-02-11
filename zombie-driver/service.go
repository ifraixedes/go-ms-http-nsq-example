package zmbdrv

import "context"

// Service specifies the operations that the Zombie Driver Service exposes.
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
	// IsZombie indicates if the driver with id is a zombie or not.
	//
	// Any implementation can returns one of the following error codes
	//
	//	* ErrNotFoundDriver
	IsZombie(ctx context.Context, id uint64) (bool, error)
}
