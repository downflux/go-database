package feature

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/flags/team"
	"github.com/downflux/go-geometry/2d/hyperrectangle"

	hnd "github.com/downflux/go-geometry/nd/hyperrectangle"
	vnd "github.com/downflux/go-geometry/nd/vector"
)

type O struct {
	AABB  hyperrectangle.R
	Flags flags.F
	Team  team.F
}

type F struct {
	id    id.ID
	aabb  hnd.M
	flags flags.F
	team  team.F
}

func New(o O) *F {
	if !Validate(o) {
		panic("cannot create feature")
	}

	f := &F{
		aabb:  hnd.New(vnd.V{0, 0}, vnd.V{0, 0}).M(),
		flags: o.Flags,
		team:  o.Team,
	}
	f.aabb.Copy(hnd.R(o.AABB))

	return f
}

func (f *F) ID() id.ID              { return f.id }
func (f *F) Flags() flags.F         { return f.flags }
func (f *F) Team() team.F           { return f.team }
func (f *F) AABB() hyperrectangle.R { return hyperrectangle.R(f.aabb.R()) }

func (f *F) SetID(x id.ID) { f.id = x }

func Validate(o O) bool { return flags.Validate(o.Flags) }
