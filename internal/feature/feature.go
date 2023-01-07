package feature

import (
	"fmt"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/team"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	v2d "github.com/downflux/go-geometry/2d/vector"
)

type O struct {
	Min   v2d.V
	Max   v2d.V
	Flags flags.F
	Team  team.T
}

type F struct {
	id    id.ID
	aabb  hyperrectangle.M
	flags flags.F
	team  team.T
}

func New(o O) *F {
	if !Validate(o.Flags) {
		panic(fmt.Sprintf("cannot create feature: invalid mask %v", o.Flags))
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
func (f *F) Team() team.T   { return f.team }

func (f *F) AABB() hyperrectangle.R {
	buf := hyperrectangle.New(vector.V{0, 0}, vector.V{0, 0}).M()
	buf.Copy(f.aabb.R())
	return buf.R()
}

func (f *F) SetID(x id.ID)      { f.id = x }
func (f *F) SetFlags(g flags.F) { f.flags = g }

func Validate(f flags.F) bool {
	if f&flags.SizeCheck != 0 {
		return false
	}
	return flags.Validate(f)
}
