package filters

import (
	"github.com/downflux/go-database/agent"
	"github.com/downflux/go-database/feature"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-geometry/2d/hyperrectangle"
	"github.com/downflux/go-geometry/2d/vector"

	dhr "github.com/downflux/go-database/geometry/hyperrectangle"
)

func AgentOnDifferentLayers(a agent.RO, b agent.RO) bool {
	m, n := a.Flags(), b.Flags()

	// Agents are allowed to overlap if (only) one of them is in the air.
	return (m^n)&flags.FTerrainAir == flags.FTerrainAir
}

// AgentIsSquishable checks if the agent a may be run over by b.
func AgentIsSquishable(a agent.RO, b agent.RO) bool {
	if AgentIsTeammate(a, b) {
		return false
	}
	if AgentOnDifferentLayers(a, b) {
		return false
	}
	return a.Size() < b.Size()
}

func AgentIsTeammate(a agent.RO, b agent.RO) bool { return a.Team() == b.Team() }

// AgentIsColliding checks if two agents are actually physically overlapping.
func AgentIsColliding(a agent.RO, b agent.RO) bool {
	if a.ID() == b.ID() {
		return false
	}

	if AgentOnDifferentLayers(a, b) {
		return false
	}

	r := a.Radius() + b.Radius()
	if vector.SquaredMagnitude(vector.Sub(a.Position(), b.Position())) > r*r {
		return false
	}
	return true

}

func AgentIsCollidingNotSquishable(a agent.RO, b agent.RO) bool {
	return !AgentIsSquishable(a, b) && AgentIsColliding(a, b)
}

func FeatureOnDifferentLayers(a agent.RO, f feature.RO) bool {
	m, n := a.Flags(), f.Flags()

	// Feature and agent are allowed to overlap if (only) one of them is in
	// the air.
	if (m^n)&flags.FTerrainAir == flags.FTerrainAir {
		return true
	}
	return false
}

func AgentIsCollidingWithFeature(a agent.RO, f feature.RO) bool {
	if FeatureOnDifferentLayers(a, f) {
		return false
	}

	if hyperrectangle.Disjoint(a.AABB(), f.AABB()) {
		return false
	}

	return dhr.IntersectCircle(f.AABB(), a.Position(), a.Radius())
}
