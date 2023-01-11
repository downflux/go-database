package flags

type F uint64

const (
	FNone F = iota

	// FTerrainAccessibleAir defines the agent can fly.
	FTerrainAccessibleAir = 1 << iota
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
)

// Validate ensures the input mask is valid. Additional checks may be added on
// top of this for per-instance validation.
func Validate(f F) bool {
	if f&TerrainAirCheck == FTerrainAir {
		return false
	}
	if f&TerrainLandCheck == FTerrainLand {
		return false
	}
	if f&TerrainSeaCheck == FTerrainSea {
		return false
	}

	return true
}
