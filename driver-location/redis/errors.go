package redis

type code uint8

// The list of specific error codes that Redis Service Location Service can
// return.
const (
	ErrInvalidRedisConfig code = iota + 1
)

func (c code) String() string {
	switch c {
	case ErrInvalidRedisConfig:
		return "InvalidRedisConfig"
	}

	return ""
}

func (c code) Message() string {
	switch c {
	case ErrInvalidRedisConfig:
		return "Impossible to connect to Redis with the specified configuration"
	}

	return ""
}
