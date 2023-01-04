package feature

import (
	"fmt"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type O struct {
	AABB  hyperrectangle.R
	Flags flags.F
}

type F struct {
	id    id.ID
	aabb  hyperrectangle.R
	flags flags.F
}

func New(o O) *F {
	if !flags.Validate(o.Flags) {
		panic(fmt.Sprintf("cannot create feature: invalid mask %v", o.Flags))
	}

	return &F{
		aabb:  o.AABB,
		flags: o.Flags,
	}
}

func (f *F) ID() id.ID              { return f.id }
func (f *F) AABB() hyperrectangle.R { return f.aabb }
func (f *F) Flags() flags.F         { return f.flags }

func (f *F) SetID(x id.ID)      { f.id = x }
func (f *F) SetFlags(g flags.F) { f.flags = g }
