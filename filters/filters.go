package filters

import (
	"github.com/downflux/go-database/agent"
	"github.com/downflux/go-database/feature"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	dhr "github.com/downflux/go-database/internal/geometry/hyperrectangle"
)

// IsColliding checks if two agents are actually physically overlapping. This
// does not care about the extra logic for e.g. squishing.
func IsColliding(a agent.RO, b agent.RO) bool {
	if a.ID() == b.ID() {
		return false
	}

	m, n := a.Flags(), b.Flags()

	// Agents are allowed to overlap if (only) one of them is in the air.
	if (m^n)&flags.FTerrainAir == flags.FTerrainAir {
		return false
	}

	r := a.Radius() + b.Radius()
	if vector.SquaredMagnitude(vector.Sub(a.Position(), b.Position())) > r*r {
		return false
	}
	return true

}

func IsSquishableColliding(a agent.RO, b agent.RO) bool {
	if IsColliding(a, b) {
		// TODO(minkezhang): Check for team.
		if a.Flags()&flags.SizeCheck > b.Flags()&flags.SizeCheck {
			return false
		}
		return true
	}
	return false
}

func IsCollidingFeature(a agent.RO, f feature.RO) bool {
	m, n := a.Flags(), f.Flags()

	// Feature and agent are allowed to overlap if (only) one of them is in
	// the air.
	if (m^n)&flags.FTerrainAir == flags.FTerrainAir {
		return false
	}

	if hyperrectangle.Disjoint(a.AABB(), f.AABB()) {
		return false
	}

	return dhr.IntersectCircle(f.AABB(), a.Position(), a.Radius())
}
