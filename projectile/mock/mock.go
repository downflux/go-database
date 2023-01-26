package mock

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/flags/team"
	"github.com/downflux/go-database/internal/projectile"
	"github.com/downflux/go-geometry/2d/hyperrectangle"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"

	roprojectile "github.com/downflux/go-database/projectile"
)

var (
	_ roprojectile.RO = &P{}
)

type P projectile.P

func New(x id.ID, o roprojectile.O) *P {
	if o.Position == nil {
		o.Position = vector.V{0, 0}
	}
	if o.TargetPosition == nil {
		o.TargetPosition = vector.V{0, 0}
	}
	if o.Velocity == nil {
		o.Velocity = vector.V{0, 0}
	}
	if o.TargetVelocity == nil {
		o.TargetVelocity = vector.V{0, 0}
	}
	if o.Heading == nil {
		o.Heading = polar.V{0, 0}
	}

	p := projectile.New(projectile.O(o))
	p.SetID(x)
	return (*P)(p)
}

func (p *P) ID() id.ID                { return (*projectile.P)(p).ID() }
func (p *P) Position() vector.V       { return (*projectile.P)(p).Position() }
func (p *P) TargetPosition() vector.V { return (*projectile.P)(p).TargetPosition() }
func (p *P) Velocity() vector.V       { return (*projectile.P)(p).Velocity() }
func (p *P) TargetVelocity() vector.V { return (*projectile.P)(p).TargetVelocity() }
func (p *P) Heading() polar.V         { return (*projectile.P)(p).Heading() }
func (p *P) Radius() float64          { return (*projectile.P)(p).Radius() }
func (p *P) Flags() flags.F           { return (*projectile.P)(p).Flags() }
func (p *P) Team() team.F             { return (*projectile.P)(p).Team() }
func (p *P) AABB() hyperrectangle.R   { return (*projectile.P)(p).AABB() }
