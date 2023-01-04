package database

import (
	"fmt"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/container"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/internal/agent"
	"github.com/downflux/go-database/internal/feature"
	"github.com/downflux/go-database/internal/projectile"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	roagent "github.com/downflux/go-database/agent"
	rofeature "github.com/downflux/go-database/feature"
	roprojectile "github.com/downflux/go-database/projectile"
)

var (
	// DefaultO provides a default set of options for setting up the
	// database. The values here are tailored to an N = 1000 simulation, and
	// is dependent on a variety of factors, e.g. surface area coverage.
	DefaultO = O{
		LeafSize:  8,
		Tolerance: 1.15,
	}
)

type O struct {
	LeafSize  int
	Tolerance float64
}

type DB struct {
	agents      map[id.ID]*agent.A
	features    map[id.ID]*feature.F
	projectiles map[id.ID]*projectile.P

	agentsBVH   container.C
	featuresBVH container.C

	counter uint64
}

func New(o O) *DB {
	return &DB{
		agents: make(map[id.ID]*agent.A, 1024),
		agentsBVH: bvh.New(bvh.O{
			LeafSize:  o.LeafSize,
			Tolerance: o.Tolerance,
		}),
		featuresBVH: bvh.New(bvh.O{
			LeafSize:  o.LeafSize,
			Tolerance: o.Tolerance,
		}),
	}
}

// AgentGetOrDie is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) AgentGetOrDie(x id.ID) roagent.RO {
	if a, ok := db.agents[x]; !ok {
		panic(fmt.Sprintf("cannot find agent %v", x))
	} else {
		return a
	}
}

// AgentInsert mutates the DB and must be called serially.
func (db *DB) AgentInsert(o roagent.O) roagent.RO {
	x := id.ID(db.counter)
	db.counter += 1

	a := agent.New(agent.O(o))
	a.SetID(x)

	db.agents[x] = a
	if err := db.agentsBVH.Insert(x, a.AABB()); err != nil {
		panic(fmt.Sprintf("cannot insert agent: %v", err))
	}

	return a
}

// AgentDelete mutates the DB and must be called serially.
func (db *DB) AgentDelete(x id.ID) {
	if _, ok := db.agents[x]; !ok {
		panic(fmt.Sprintf("cannot find agent %v", x))
	}

	delete(db.agents, x)
	if err := db.agentsBVH.Remove(x); err != nil {
		panic(fmt.Sprintf("cannot delete agent: %v", err))
	}
}

// AgentQuery is a read-only operation and may be called concurrently with other
// read-only operations.
func (db *DB) AgentQuery(q hyperrectangle.R, filter func(a roagent.RO) bool) []roagent.RO {
	candidates := db.agentsBVH.BroadPhase(q)

	results := make([]roagent.RO, 0, len(candidates))
	for _, x := range candidates {
		a := db.agents[x]
		if filter(a) {
			results = append(results, a)
		}
	}
	return results
}

// AgentSetPosition mutates the DB and must be called serially.
func (db *DB) AgentSetPosition(x id.ID, v vector.V) {
	a := db.AgentGetOrDie(x)

	a.(*agent.A).SetPosition(v)
	db.agentsBVH.Update(x, a.AABB())
}

// AgentSetVelocity mutates the DB, but may be called concurently with other
// invocations on different agents.
func (db *DB) AgentSetVelocity(x id.ID, v vector.V) {
	db.AgentGetOrDie(x).(*agent.A).SetVelocity(v)
}

// AgentSetTargetVelocity mutates the DB, but may be called concurrently with
// other invocations on different agents.
func (db *DB) AgentSetTargetVelocity(x id.ID, v vector.V) {
	db.AgentGetOrDie(x).(*agent.A).SetTargetVelocity(v)
}

// AgentSetHeading mutates the DB, but may be called concurrently with other
// invocations on different agents.
func (db *DB) AgentSetHeading(x id.ID, v polar.V) {
	db.AgentGetOrDie(x).(*agent.A).SetHeading(v)
}

// FeatureGetOrDie is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) FeatureGetOrDie(x id.ID) rofeature.RO {
	if a, ok := db.features[x]; !ok {
		panic(fmt.Sprintf("cannot find feature %v", x))
	} else {
		return a
	}
}

// FeatureInsert mutates the DB and must be called serially.
func (db *DB) FeatureInsert(o rofeature.O) rofeature.RO {
	x := id.ID(db.counter)
	db.counter += 1

	a := feature.New(feature.O(o))
	a.SetID(x)

	db.features[x] = a
	if err := db.featuresBVH.Insert(x, a.AABB()); err != nil {
		panic(fmt.Sprintf("cannot insert feature: %v", err))
	}

	return a
}

// FeatureDelete mutates the DB and must be called serially.
func (db *DB) FeatureDelete(x id.ID) {
	if _, ok := db.features[x]; !ok {
		panic(fmt.Sprintf("cannot find feature %v", x))
	}

	delete(db.features, x)
	if err := db.featuresBVH.Remove(x); err != nil {
		panic(fmt.Sprintf("cannot delete feature: %v", err))
	}
}

// FeatureQuery is a read-only operation and may be called concurrently with other
// read-only operations.
func (db *DB) FeatureQuery(q hyperrectangle.R, filter func(a rofeature.RO) bool) []rofeature.RO {
	candidates := db.featuresBVH.BroadPhase(q)

	results := make([]rofeature.RO, 0, len(candidates))
	for _, x := range candidates {
		a := db.features[x]
		if filter(a) {
			results = append(results, a)
		}
	}
	return results
}

// ProjectileGetOrDie is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) ProjectileGetOrDie(x id.ID) roprojectile.RO {
	if a, ok := db.projectiles[x]; !ok {
		panic(fmt.Sprintf("cannot find projectile %v", x))
	} else {
		return a
	}
}

// ProjectileInsert mutates the DB and must be called serially.
func (db *DB) ProjectileInsert(o roprojectile.O) roprojectile.RO {
	x := id.ID(db.counter)
	db.counter += 1

	a := projectile.New(projectile.O(o))
	a.SetID(x)

	db.projectiles[x] = a
	return a
}

// ProjectileDelete mutates the DB and must be called serially.
func (db *DB) ProjectileDelete(x id.ID) {
	if _, ok := db.projectiles[x]; !ok {
		panic(fmt.Sprintf("cannot find projectile %v", x))
	}

	delete(db.projectiles, x)
}

// ProjectileSetPosition mutates the DB, but may be called concurrently with
// other invocations on different projectiles.
func (db *DB) ProjectileSetPosition(x id.ID, v vector.V) {
	a := db.ProjectileGetOrDie(x)

	a.(*projectile.P).SetPosition(v)
}

// ProjectileSetVelocity mutates the DB, but may be called concurently with
// other invocations on different projectiles.
func (db *DB) ProjectileSetVelocity(x id.ID, v vector.V) {
	db.ProjectileGetOrDie(x).(*projectile.P).SetVelocity(v)
}

// ProjectileSetTargetVelocity mutates the DB, but may be called concurrently
// with other invocations on different projectiles.
func (db *DB) ProjectileSetTargetVelocity(x id.ID, v vector.V) {
	db.ProjectileGetOrDie(x).(*projectile.P).SetTargetVelocity(v)
}

// ProjectileSetHeading mutates the DB, but may be called concurrently with
// other invocations on different projectiles.
func (db *DB) ProjectileSetHeading(x id.ID, v polar.V) {
	db.ProjectileGetOrDie(x).(*projectile.P).SetHeading(v)
}
