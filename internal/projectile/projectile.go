package projectile

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/team"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

type O struct {
	Position       vector.V
	TargetPosition vector.V
	Velocity       vector.V
	TargetVelocity vector.V
	Heading        polar.V
	Radius         float64
	Flags          flags.F
	Team           team.T
}

type P struct {
	id             id.ID
	position       vector.M
	targetPosition vector.M
	velocity       vector.M
	targetVelocity vector.M
	heading        polar.M
	radius         float64
	flags          flags.F
	team           team.T
}

func New(o O) *P {
	if !Validate(o) {
		panic("cannot create agent")
	}

	p := &P{
		position:       vector.M{0, 0},
		targetPosition: vector.M{0, 0},
		velocity:       vector.M{0, 0},
		targetVelocity: vector.M{0, 0},
		heading:        polar.M{0, 0},
		radius:         o.Radius,
		flags:          o.Flags,
		team:           o.Team,
	}

	p.position.Copy(o.Position)
	p.targetPosition.Copy(o.TargetPosition)
	p.velocity.Copy(o.Velocity)
	p.targetVelocity.Copy(o.TargetVelocity)
	p.heading.Copy(o.Heading)

	return p
}

func (p *P) ID() id.ID       { return p.id }
func (p *P) Flags() flags.F  { return p.flags }
func (p *P) Team() team.T    { return p.team }
func (p *P) Radius() float64 { return p.radius }

func (p *P) Position() vector.V {
	buf := vector.M{0, 0}
	buf.Copy(p.position.V())
	return buf.V()
}

func (p *P) TargetPosition() vector.V {
	buf := vector.M{0, 0}
	buf.Copy(p.targetPosition.V())
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
func (p *P) SetTargetPosition(v vector.V) { p.targetPosition.Copy(v) }
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

func Validate(o O) bool {
	if o.Radius == 0 {
		return false
	}
	if o.Flags&flags.FSizeProjectile == 0 {
		return false
	}
	return flags.Validate(o.Flags)
}
