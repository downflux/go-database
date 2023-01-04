package agent

import (
	"fmt"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-data/flags"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

type O struct {
	ID                 id.ID
	Position           vector.V
	Velocity           vector.V
	TargetVelocity     vector.V
	Heading            polar.V
	Radius             float64
	MaxVelocity        float64
	MaxAngularVelocity float64
	MaxAcceleration    float64
	Flags              flags.F
}

type A struct {
	id id.ID

	position vector.M

	// velocity is the actual tick-to-tick velocity. This is used for smoothing
	// over acceleration values.
	velocity       vector.M
	targetVelocity vector.M

	// heading is a unit polar vector whose angular component is oriented to
	// the positive X-axis. The angle is calculated according to normal 2D
	// rotational rules, i.e. a vector lying on the positive Y-axis has an
	// angular componet of π / 2.
	heading polar.M

	radius float64

	maxVelocity        float64
	maxAngularVelocity float64
	maxAcceleration    float64

	flags flags.F
}

func New(o O) *A {
	if !flags.Validate(o.Flags) {
		panic(fmt.Sprintf("cannot create agent: invalid mask %v", o.Flags))
	}

	a := &A{
		id:                 o.ID,
		position:           vector.M{0, 0},
		velocity:           vector.M{0, 0},
		targetVelocity:     vector.M{0, 0},
		heading:            polar.M{0, 0},
		radius:             o.Radius,
		maxVelocity:        o.MaxVelocity,
		maxAngularVelocity: o.MaxAngularVelocity,
		maxAcceleration:    o.MaxAcceleration,
		flags:              o.Flags,
	}

	a.position.Copy(o.Position)
	a.velocity.Copy(o.Velocity)
	a.targetVelocity.Copy(o.TargetVelocity)
	a.heading.Copy(o.Heading)

	return a
}

func (a *A) ID() id.ID                   { return a.id }
func (a *A) Flags() flags.F              { return a.flags }
func (a *A) Radius() float64             { return a.radius }
func (a *A) MaxVelocity() float64        { return a.maxVelocity }
func (a *A) MaxAngularVelocity() float64 { return a.maxAngularVelocity }
func (a *A) MaxAcceleration() float64    { return a.maxAcceleration }

func (a *A) Position() vector.V {
	buf := vector.M{0, 0}
	buf.Copy(a.position.V())
	return buf.V()
}

func (a *A) Velocity() vector.V {
	buf := vector.M{0, 0}
	buf.Copy(a.velocity.V())
	return buf.V()
}

func (a *A) TargetVelocity() vector.V {
	buf := vector.M{0, 0}
	buf.Copy(a.targetVelocity.V())
	return buf.V()
}

func (a *A) Heading() polar.V {
	buf := polar.M{0, 0}
	buf.Copy(a.heading.V())
	return buf.V()
}

func (a *A) SetID(x id.ID)                { a.id = x }
func (a *A) SetPosition(v vector.V)       { a.position.Copy(v) }
func (a *A) SetVelocity(v vector.V)       { a.velocity.Copy(v) }
func (a *A) SetTargetVelocity(v vector.V) { a.targetVelocity.Copy(v) }
func (a *A) SetHeading(v polar.V)         { a.heading.Copy(v) }
func (a *A) SetFlags(f flags.F)           { a.flags = f }

func (a *A) AABB() hyperrectangle.R {
	x, y := a.position.X(), a.position.Y()
	r := a.radius

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
