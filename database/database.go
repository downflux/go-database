package database

import (
	"fmt"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/container"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-data/internal/agent"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	pagent "github.com/downflux/go-data/agent"
)

type O struct {
	LeafSize  int
	Tolerance float64
}

type DB struct {
	agents map[id.ID]*agent.A

	agentsBVH container.C

	counter uint64
}

func New(o O) *DB {
	return &DB{
		agents: make(map[id.ID]*agent.A, 1024),
		agentsBVH: bvh.New(bvh.O{
			LeafSize:  o.LeafSize,
			Tolerance: o.Tolerance,
		}),
	}
}

// AgentGetOrDie is a read-only operation and may be called concurrently with
// other read-only operations.
func (db *DB) AgentGetOrDie(x id.ID) pagent.RO {
	if a, ok := db.agents[x]; !ok {
		panic(fmt.Sprintf("cannot find agent %v", x))
	} else {
		return a
	}
}

// AgentInsert mutates the DB and must be called serially.
func (db *DB) AgentInsert(o pagent.O) pagent.RO {
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
func (db *DB) AgentQuery(q hyperrectangle.R, filter func(a pagent.RO) bool) []pagent.RO {
	candidates := db.agentsBVH.BroadPhase(q)

	results := make([]pagent.RO, 0, len(candidates))
	for _, x := range candidates {
		a := db.agents[x]
		if filter(a) {
			results = append(results, a)
		}
	}
	return results
}
