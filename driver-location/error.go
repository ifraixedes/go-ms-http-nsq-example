package drvloc

type code uint8

// The list of error codes that any implementation of the Driver Location
// Service can return.
const (
	ErrAbortedCtx code = iota + 1

	ErrInvalidDataFormatStore

	ErrNotFoundDriver

	ErrUnexpectedErrorDeserializer
	ErrUnexpectedErrorSerializer
	ErrUnexpectedErrorSystem
	ErrUnexpectedErrorStore
	ErrUnexpectedErrorTransport
)

func (c code) String() string {
	switch c {
	case ErrAbortedCtx:
		return "AbortedCtx"
	case ErrInvalidDataFormatStore:
		return "ErrInvalidDataFormatStore"
	case ErrNotFoundDriver:
		return "NotFoundDriver"
	case ErrUnexpectedErrorDeserializer:
		return "UnexpectedErrorDeserializer"
	case ErrUnexpectedErrorSerializer:
		return "UnexpectedErrorSerializer"
	case ErrUnexpectedErrorStore:
		return "UnexpectedErrorStore"
	case ErrUnexpectedErrorSystem:
		return "UnexpectedErrorSystem"
	case ErrUnexpectedErrorTransport:
		return "UnexpectedErrorTransport"
	}

	return ""
}

func (c code) Message() string {
	switch c {
	case ErrAbortedCtx:
		return "The operation has been aborted due to a context cancellation or deadline exceeded"
	case ErrInvalidDataFormatStore:
		return "The store contains data of unexpected format"
	case ErrNotFoundDriver:
		return "The driver is not found"
	case ErrUnexpectedErrorDeserializer:
		return "Aunexpected error when de-serializing the data happened"
	case ErrUnexpectedErrorSerializer:
		return "An unexpected error when serializing the data happened"
	case ErrUnexpectedErrorSystem:
		return "A unexpected system error happened"
	case ErrUnexpectedErrorStore:
		return "The underlying store has returned an unexpected error"
	case ErrUnexpectedErrorTransport:
		return "The underlying communication transport returned an unexpected error"
	}

	return ""
}
