package drvloc

type code uint8

// The list of specific error codes that Redis Service Location Service can
// return.
const (
	ErrInvalidRules code = iota + 1
	ErrRequiredDrvLocSvc
)

func (c code) String() string {
	switch c {
	case ErrInvalidRules:
		return "InvalidRules"
	case ErrRequiredDrvLocSvc:
		return "RequiredDrvLocSvc"
	}

	return ""
}

func (c code) Message() string {
	switch c {
	case ErrInvalidRules:
		return "Rules values aren't valid"
	case ErrRequiredDrvLocSvc:
		return "Driver Service Location cannot be nil beause it's required"
	}

	return ""
}
