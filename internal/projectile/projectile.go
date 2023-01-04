package projectile

import (
	"fmt"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

type O struct {
	Position       vector.V
	Velocity       vector.V
	TargetVelocity vector.V
	Heading        polar.V
	Radius         float64
	MaxVelocity    float64
	Flags          flags.F
}

type P struct {
	id             id.ID
	position       vector.M
	velocity       vector.M
	targetVelocity vector.M
	heading        polar.M
	radius         float64
	maxVelocity    float64
	flags          flags.F
}

func New(o O) *P {
	if !flags.Validate(o.Flags) {
		panic(fmt.Sprintf("cannot create agent: invalid mask %v", o.Flags))
	}

	p := &P{
		position:       vector.M{0, 0},
		velocity:       vector.M{0, 0},
		targetVelocity: vector.M{0, 0},
		heading:        polar.M{0, 0},
		radius:         o.Radius,
		maxVelocity:    o.MaxVelocity,
		flags:          o.Flags,
	}

	p.position.Copy(o.Position)
	p.velocity.Copy(o.Velocity)
	p.targetVelocity.Copy(o.TargetVelocity)
	p.heading.Copy(o.Heading)

	return p
}

func (p *P) ID() id.ID            { return p.id }
func (p *P) Flags() flags.F       { return p.flags }
func (p *P) Radius() float64      { return p.radius }
func (p *P) MaxVelocity() float64 { return p.maxVelocity }

func (p *P) Position() vector.V {
	buf := vector.M{0, 0}
	buf.Copy(p.position.V())
	return buf.V()
}
func (p *P) Velocity() vector.V {
	buf := vector.M{0, 0}
	buf.Copy(p.velocity.V())
	return buf.V()
}

func (p *P) TargetVelocity() vector.V {
	buf := vector.M{0, 0}
	buf.Copy(p.targetVelocity.V())
	return buf.V()
}

func (p *P) Heading() polar.V {
	buf := polar.M{0, 0}
	buf.Copy(p.heading.V())
	return buf.V()
}

func (p *P) SetID(x id.ID)                { p.id = x }
func (p *P) SetPosition(v vector.V)       { p.position.Copy(v) }
func (p *P) SetVelocity(v vector.V)       { p.velocity.Copy(v) }
func (p *P) SetTargetVelocity(v vector.V) { p.targetVelocity.Copy(v) }
func (p *P) SetHeading(v polar.V)         { p.heading.Copy(v) }
func (p *P) SetFlags(f flags.F)           { p.flags = f }

func (p *P) AABB() hyperrectangle.R {
	x, y := p.position.X(), p.position.Y()
	r := p.radius

	return *hyperrectangle.New(
		vnd.V{
			x - r,
			y - r,
		},
		vnd.V{
			x + r,
			y + r,
		},
	)

}