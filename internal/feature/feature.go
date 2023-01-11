package feature

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/flags/team"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	v2d "github.com/downflux/go-geometry/2d/vector"
)

type O struct {
	Min   v2d.V
	Max   v2d.V
	Flags flags.F
	Team  team.F
}

type F struct {
	id    id.ID
	aabb  hyperrectangle.M
	flags flags.F
	team  team.F
}

func New(o O) *F {
	if !Validate(o) {
		panic("cannot create feature")
	}

	f := &F{
		aabb:  hyperrectangle.New(vector.V{0, 0}, vector.V{0, 0}).M(),
		flags: o.Flags,
		team:  o.Team,
	}

	f.aabb.Min().Copy(vector.V(o.Min))
	f.aabb.Max().Copy(vector.V(o.Max))

	return f
}

func (f *F) ID() id.ID      { return f.id }
func (f *F) Flags() flags.F { return f.flags }
func (f *F) Team() team.F   { return f.team }

func (f *F) AABB() hyperrectangle.R {
	buf := hyperrectangle.New(vector.V{0, 0}, vector.V{0, 0}).M()
	buf.Copy(f.aabb.R())
	return buf.R()
}

func (f *F) SetID(x id.ID) { f.id = x }

func Validate(o O) bool {
	return flags.Validate(o.Flags)
}
