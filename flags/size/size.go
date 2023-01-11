package size

type F uint64

const (
	FNone F = iota

	FSmall
	FMedium
	FLarge
)

func Validate(f F) bool { return f > FNone }
