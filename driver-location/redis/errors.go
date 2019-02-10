package redis

type code uint8

// The list of specific error codes that Redis Service Location Service can
// return.
const (
	ErrInvalidRedisConfig code = iota + 1

	ErrClosedService
)

func (c code) String() string {
	switch c {
	case ErrInvalidRedisConfig:
		return "InvalidRedisConfig"
	case ErrClosedService:
		return "ClosedService"
	}

	return ""
}

func (c code) Message() string {
	switch c {
	case ErrInvalidRedisConfig:
		return "Impossible to connect to Redis with the specified configuration"
	case ErrClosedService:
		return "Service has been closed so it cannot be used anymore, you must create a new instance"
	}

	return ""
}
