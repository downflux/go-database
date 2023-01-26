package mock

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/flags/team"
	"github.com/downflux/go-database/internal/feature"
	"github.com/downflux/go-geometry/2d/hyperrectangle"
	"github.com/downflux/go-geometry/2d/vector"

	rofeature "github.com/downflux/go-database/feature"
)

var (
	_ rofeature.RO = &F{}
)

type F feature.F

func New(x id.ID, o rofeature.O) *F {
	if o.AABB.Min() == nil || o.AABB.Max() == nil {
		o.AABB = *hyperrectangle.New(vector.V{0, 0}, vector.V{0, 0})
	}

	f := feature.New(feature.O(o))
	f.SetID(x)
	return (*F)(f)
}

func (f *F) ID() id.ID              { return (*feature.F)(f).ID() }
func (f *F) Flags() flags.F         { return (*feature.F)(f).Flags() }
func (f *F) Team() team.F           { return (*feature.F)(f).Team() }
func (f *F) AABB() hyperrectangle.R { return (*feature.F)(f).AABB() }
