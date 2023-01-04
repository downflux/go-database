package mock

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/internal/feature"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	rofeature "github.com/downflux/go-database/feature"
)

var (
	_ rofeature.RO = &F{}
)

type F feature.F

func New(x id.ID, o rofeature.O) *F {
	f := feature.New(feature.O(o))
	f.SetID(x)
	return (*F)(f)
}

func (f *F) ID() id.ID              { return (*feature.F)(f).ID() }
func (f *F) Flags() flags.F         { return (*feature.F)(f).Flags() }
func (f *F) AABB() hyperrectangle.R { return (*feature.F)(f).AABB() }
