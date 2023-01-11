package move

// F is the move mode of the object.
type F uint64

const (
	FNone F = iota

	FAvoidance = 1 << iota
	FSeek
	FArrival
	FAlignment
	FCoherence
	FSeparation
)

const (
	FFlocking = FAlignment | FCoherence | FSeparation
)

func Validate(f F) bool {
	// An object can either seek or arrive, but not both.
	if f&FSeek == FSeek && f&FArrival == FArrival {
		return false
	}
	return true
}
