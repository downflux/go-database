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
		agents:      make(map[id.ID]*agent.A, 1024),
		features:    make(map[id.ID]*feature.F, 1024),
		projectiles: make(map[id.ID]*projectile.P, 1024),
		agentsBVH: bvh.New(bvh.O{
			K:         2,
			LeafSize:  o.LeafSize,
			Tolerance: o.Tolerance,
		}),
		featuresBVH: bvh.New(bvh.O{
			K:         2,
			LeafSize:  o.LeafSize,
			Tolerance: o.Tolerance,
		}),
	}
}

// GetAgentOrDie is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) GetAgentOrDie(x id.ID) roagent.RO {
	if a, ok := db.agents[x]; !ok {
		panic(fmt.Sprintf("cannot find agent %v", x))
	} else {
		return a
	}
}

// InsertAgent mutates the DB and must be called serially.
func (db *DB) InsertAgent(o roagent.O) roagent.RO {
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

// ListAgents returns all agents in the DB and may be called concurrently with
// other read-only operations. All values must be consumed (or the channel is
// garbage collected) before a DB mutation occurs.
func (db *DB) ListAgents() <-chan roagent.RO {
	ch := make(chan roagent.RO, 256)
	go func(ch chan<- roagent.RO) {
		defer close(ch)
		for _, a := range db.agents {
			ch <- a
		}
	}(ch)
	return ch
}

// DeleteAgent mutates the DB and must be called serially.
func (db *DB) DeleteAgent(x id.ID) {
	if _, ok := db.agents[x]; !ok {
		panic(fmt.Sprintf("cannot find agent %v", x))
	}

	delete(db.agents, x)
	if err := db.agentsBVH.Remove(x); err != nil {
		panic(fmt.Sprintf("cannot delete agent: %v", err))
	}
}

// QueryAgents is a read-only operation and may be called concurrently with other
// read-only operations.
func (db *DB) QueryAgents(q hyperrectangle.R, filter func(a roagent.RO) bool) []roagent.RO {
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

// SetAgentPosition mutates the DB and must be called serially.
func (db *DB) SetAgentPosition(x id.ID, v vector.V) {
	a := db.GetAgentOrDie(x)

	a.(*agent.A).SetPosition(v)
	db.agentsBVH.Update(x, a.AABB())
}

// SetAgentVelocity mutates the DB, but may be called concurently with other
// invocations on different agents.
func (db *DB) SetAgentVelocity(x id.ID, v vector.V) {
	db.GetAgentOrDie(x).(*agent.A).SetVelocity(v)
}

// SetAgentTargetVelocity mutates the DB, but may be called concurrently with
// other invocations on different agents.
func (db *DB) SetAgentTargetVelocity(x id.ID, v vector.V) {
	db.GetAgentOrDie(x).(*agent.A).SetTargetVelocity(v)
}

// SetAgentHeading mutates the DB, but may be called concurrently with other
// invocations on different agents.
func (db *DB) SetAgentHeading(x id.ID, v polar.V) {
	db.GetAgentOrDie(x).(*agent.A).SetHeading(v)
}

// GetFeatureOrDie is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) GetFeatureOrDie(x id.ID) rofeature.RO {
	if a, ok := db.features[x]; !ok {
		panic(fmt.Sprintf("cannot find feature %v", x))
	} else {
		return a
	}
}

// InsertFeature mutates the DB and must be called serially.
func (db *DB) InsertFeature(o rofeature.O) rofeature.RO {
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

// ListFeatures returns all features in the DB and may be called concurrently
// with other read-only operations. All values must be consumed (or the channel
// is garbage collected) before a DB mutation occurs.
func (db *DB) ListFeatures() <-chan rofeature.RO {
	ch := make(chan rofeature.RO, 256)
	go func(ch chan<- rofeature.RO) {
		defer close(ch)
		for _, a := range db.features {
			ch <- a
		}
	}(ch)
	return ch
}

// DeleteFeature mutates the DB and must be called serially.
func (db *DB) DeleteFeature(x id.ID) {
	if _, ok := db.features[x]; !ok {
		panic(fmt.Sprintf("cannot find feature %v", x))
	}

	delete(db.features, x)
	if err := db.featuresBVH.Remove(x); err != nil {
		panic(fmt.Sprintf("cannot delete feature: %v", err))
	}
}

// QueryFeatures is a read-only operation and may be called concurrently with other
// read-only operations.
func (db *DB) QueryFeatures(q hyperrectangle.R, filter func(a rofeature.RO) bool) []rofeature.RO {
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

// GetProjectileOrDie is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) GetProjectileOrDie(x id.ID) roprojectile.RO {
	if a, ok := db.projectiles[x]; !ok {
		panic(fmt.Sprintf("cannot find projectile %v", x))
	} else {
		return a
	}
}

// InsertProjectile mutates the DB and must be called serially.
func (db *DB) InsertProjectile(o roprojectile.O) roprojectile.RO {
	x := id.ID(db.counter)
	db.counter += 1

	a := projectile.New(projectile.O(o))
	a.SetID(x)

	db.projectiles[x] = a
	return a
}

// ListProjectiles returns all projectiles in the DB and may be called
// concurrently with other read-only operations. All values must be consumed (or
// the channel is garbage collected) before a DB mutation occurs.
func (db *DB) ListProjectiles() <-chan roprojectile.RO {
	ch := make(chan roprojectile.RO, 256)
	go func(ch chan<- roprojectile.RO) {
		defer close(ch)
		for _, a := range db.projectiles {
			ch <- a
		}
	}(ch)
	return ch
}

// DeleteProjectile mutates the DB and must be called serially.
func (db *DB) DeleteProjectile(x id.ID) {
	if _, ok := db.projectiles[x]; !ok {
		panic(fmt.Sprintf("cannot find projectile %v", x))
	}

	delete(db.projectiles, x)
}

// SetProjectilePosition mutates the DB, but may be called concurrently with
// other invocations on different projectiles.
func (db *DB) SetProjectilePosition(x id.ID, v vector.V) {
	a := db.GetProjectileOrDie(x)

	a.(*projectile.P).SetPosition(v)
}

// SetProjectileVelocity mutates the DB, but may be called concurently with
// other invocations on different projectiles.
func (db *DB) SetProjectileVelocity(x id.ID, v vector.V) {
	db.GetProjectileOrDie(x).(*projectile.P).SetVelocity(v)
}

// SetProjectileTargetVelocity mutates the DB, but may be called concurrently
// with other invocations on different projectiles.
func (db *DB) SetProjectileTargetVelocity(x id.ID, v vector.V) {
	db.GetProjectileOrDie(x).(*projectile.P).SetTargetVelocity(v)
}

// SetProjectileHeading mutates the DB, but may be called concurrently with
// other invocations on different projectiles.
func (db *DB) SetProjectileHeading(x id.ID, v polar.V) {
	db.GetProjectileOrDie(x).(*projectile.P).SetHeading(v)
}
