package feature

import (
	"fmt"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type O struct {
	AABB  hyperrectangle.R
	Flags flags.F
}

type F struct {
	id    id.ID
	aabb  hyperrectangle.M
	flags flags.F
}

func New(o O) *F {
	if !flags.Validate(o.Flags) {
		panic(fmt.Sprintf("cannot create feature: invalid mask %v", o.Flags))
	}

	f := &F{
		aabb:  hyperrectangle.New(vector.V{0, 0}, vector.V{0, 0}).M(),
		flags: o.Flags,
	}

	f.aabb.Copy(o.AABB)

	return f
}

func (f *F) ID() id.ID      { return f.id }
func (f *F) Flags() flags.F { return f.flags }

func (f *F) AABB() hyperrectangle.R {
	buf := hyperrectangle.New(vector.V{0, 0}, vector.V{0, 0}).M()
	buf.Copy(f.aabb.R())
	return buf.R()
}

func (f *F) SetID(x id.ID)      { f.id = x }
func (f *F) SetFlags(g flags.F) { f.flags = g }
