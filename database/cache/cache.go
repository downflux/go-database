package cache

import (
	"fmt"

	"github.com/downflux/go-bvh/container"
	"github.com/downflux/go-bvh/container/bruteforce"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/agent"
	"github.com/downflux/go-database/feature"
	"github.com/downflux/go-database/projectile"
	"github.com/downflux/go-geometry/2d/hyperrectangle"

	roagent "github.com/downflux/go-database/agent"
	rofeature "github.com/downflux/go-database/feature"
	roprojectile "github.com/downflux/go-database/projectile"
	hnd "github.com/downflux/go-geometry/nd/hyperrectangle"
)

type O struct {
	Agents      []agent.RO
	Features    []feature.RO
	Projectiles []projectile.RO
}

type DB struct {
	agents      map[id.ID]agent.RO
	features    map[id.ID]feature.RO
	projectiles map[id.ID]projectile.RO

	agentsBVH   container.C
	featuresBVH container.C
}

func New(o O) *DB {
	db := &DB{
		agents:      make(map[id.ID]agent.RO, len(o.Agents)),
		features:    make(map[id.ID]feature.RO, len(o.Features)),
		projectiles: make(map[id.ID]projectile.RO, len(o.Projectiles)),

		agentsBVH:   bruteforce.New(),
		featuresBVH: bruteforce.New(),
	}

	for _, a := range o.Agents {
		db.agents[a.ID()] = a
		db.agentsBVH.Insert(a.ID(), hnd.R(a.AABB()))
	}

	for _, f := range o.Features {
		db.features[f.ID()] = f
		db.featuresBVH.Insert(f.ID(), hnd.R(f.AABB()))
	}

	for _, p := range o.Projectiles {
		db.projectiles[p.ID()] = p
	}

	return db
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

// GetFeatureOrDie is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) GetFeatureOrDie(x id.ID) rofeature.RO {
	if a, ok := db.features[x]; !ok {
		panic(fmt.Sprintf("cannot find feature %v", x))
	} else {
		return a
	}
}

// GetProjectileOrDie is a read-only operation and may be called concurrently
// with other read-only operations.
func (db *DB) GetProjectileOrDie(x id.ID) roprojectile.RO {
	if a, ok := db.projectiles[x]; !ok {
		panic(fmt.Sprintf("cannot find projectile %v", x))
	} else {
		return a
	}
}

// ListAgents returns all agents in the DB. There are serveral use-cases for
// this method which changes the invocation pattern.
//
//  1. Read-only operations on agents may consume this output concurrently.
//  1. Agent-specific mutations should first iterate over the returned values
//     and create a proposal batch of changes. If these changes do not modify
//     the BVH, they may be run in parallel. Changes to the BVH (e.g.
//     SetAgentPosition) must be done serially.
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

// ListFeatures returns all features in the DB. There are serveral use-cases for
// this method which changes the invocation pattern.
//
// See ListAgents for more information.
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

// ListProjectiles returns all projectiles in the DB. There are serveral
// use-cases for this method which changes the invocation pattern.
//
// See ListAgents for more information.
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

// QueryAgents is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) QueryAgents(q hyperrectangle.R, filter func(a roagent.RO) bool) []roagent.RO {
	candidates := db.agentsBVH.BroadPhase(hnd.R(q))

	results := make([]roagent.RO, 0, len(candidates))
	for _, x := range candidates {
		a := db.agents[x]
		if filter(a) {
			results = append(results, a)
		}
	}
	return results
}

// QueryFeatures is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) QueryFeatures(q hyperrectangle.R, filter func(a rofeature.RO) bool) []rofeature.RO {
	candidates := db.featuresBVH.BroadPhase(hnd.R(q))

	results := make([]rofeature.RO, 0, len(candidates))
	for _, x := range candidates {
		a := db.features[x]
		if filter(a) {
			results = append(results, a)
		}
	}
	return results
}
