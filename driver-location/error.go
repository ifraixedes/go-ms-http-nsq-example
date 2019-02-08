package drvloc

type code uint8

// The list of error codes that any implementation of the Driver Location
// Service can return.
const (
	ErrNotFoundDriver code = iota + 1

	ErrUnexpectedErrorStore
)

func (c code) String() string {
	switch c {
	case ErrNotFoundDriver:
		return "NotFoundDriver"
	case ErrUnexpectedErrorStore:
		return "UnexpectedErrorStore"
	}

	return ""
}

func (c code) Message() string {
	switch c {
	case ErrNotFoundDriver:
		return "The driver is not found"
	case ErrUnexpectedErrorStore:
		return "The underlying store has returned an unexpected error"
	}

	return ""
}
