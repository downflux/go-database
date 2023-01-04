package flags

type F uint64

const (
	FNone F = iota

	// FSizeProjectile defines bullets and missiles. Only one size class can
	// be active on an agent at any given time.
	FSizeProjectile = 1 << iota
	FSizeSmall
	FSizeFedium
	FSizeLarge

	// FTerrainAccessibleAir defines the agent can fly.
	FTerrainAccessibleAir
	FTerrainAccessibleLand
	FTerrainAccessibleSea

	// FTerrainAir defines the current map layer the agent is occupying. Ony
	// one terrain layer can be active on an agent at any given time. The
	// agent must always be occupying a map layer which it has access to,
	// e.g. tanks activate FTerrainAccessibleLand and cannot activate
	// FTerrainAir.
	FTerrainAir
	FTerrainLand
	FTerrainSea
)

const (
	TerrainAirCheck  = FTerrainAccessibleAir | FTerrainAir
	TerrainLandCheck = FTerrainAccessibleLand | FTerrainLand
	TerrainSeaCheck  = FTerrainAccessibleSea | FTerrainSea

	SizeCheck = FSizeProjectile | FSizeSmall | FSizeFedium | FSizeLarge
)

// Validate ensures the input mask is valid. Additional checks may be added on
// top of this for per-instance validation.
func Validate(m F) bool {
	n := 0
	if m&FSizeProjectile != 0 {
		n++
	}
	if m&FSizeSmall != 0 {
		n++
	}
	if m&FSizeFedium != 0 {
		n++
	}
	if m&FSizeLarge != 0 {
		n++
	}
	if n > 1 {
		return false
	}

	if m&TerrainAirCheck == FTerrainAir {
		return false
	}
	if m&TerrainLandCheck == FTerrainLand {
		return false
	}
	if m&TerrainSeaCheck == FTerrainSea {
		return false
	}

	return true
}
