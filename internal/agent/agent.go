package agent

import (
	"fmt"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/flags/move"
	"github.com/downflux/go-database/flags/size"
	"github.com/downflux/go-database/flags/team"
	"github.com/downflux/go-geometry/2d/hyperrectangle"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
)

type O struct {
	Position           vector.V
	TargetPosition     vector.V
	Velocity           vector.V
	TargetVelocity     vector.V
	Heading            polar.V
	Radius             float64
	Mass               float64
	MaxVelocity        float64
	MaxAngularVelocity float64
	MaxAcceleration    float64
	Flags              flags.F
	Size               size.F
	Team               team.F
	Move               move.F
}

type A struct {
	id id.ID

	position       vector.M
	targetPosition vector.M

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
	mass   float64

	maxVelocity        float64
	maxAngularVelocity float64
	maxAcceleration    float64

	flags flags.F
	size  size.F
	team  team.F
	move  move.F
}

func New(o O) *A {
	if !Validate(o) {
		panic("cannot create agent")
	}

	a := &A{
		position:           vector.M{0, 0},
		targetPosition:     vector.M{0, 0},
		velocity:           vector.M{0, 0},
		targetVelocity:     vector.M{0, 0},
		heading:            polar.M{0, 0},
		radius:             o.Radius,
		mass:               o.Mass,
		maxVelocity:        o.MaxVelocity,
		maxAngularVelocity: o.MaxAngularVelocity,
		maxAcceleration:    o.MaxAcceleration,
		flags:              o.Flags,
		size:               o.Size,
		team:               o.Team,
		move:               o.Move,
	}

	a.position.Copy(o.Position)
	a.targetPosition.Copy(o.TargetPosition)
	a.velocity.Copy(o.Velocity)
	a.targetVelocity.Copy(o.TargetVelocity)
	a.heading.Copy(o.Heading)

	return a
}

func (a *A) ID() id.ID                   { return a.id }
func (a *A) Flags() flags.F              { return a.flags }
func (a *A) Team() team.F                { return a.team }
func (a *A) Radius() float64             { return a.radius }
func (a *A) Mass() float64               { return a.mass }
func (a *A) MaxVelocity() float64        { return a.maxVelocity }
func (a *A) MaxAngularVelocity() float64 { return a.maxAngularVelocity }
func (a *A) MaxAcceleration() float64    { return a.maxAcceleration }
func (a *A) Size() size.F                { return a.size }
func (a *A) MoveMode() move.F            { return a.move }

// Position returns the current position of the agent.
//
// External callers must not set the position via the position reference object
// returned here.
func (a *A) Position() vector.V       { return a.position.V() }
func (a *A) TargetPosition() vector.V { return a.targetPosition.V() }
func (a *A) Velocity() vector.V       { return a.velocity.V() }
func (a *A) TargetVelocity() vector.V { return a.targetVelocity.V() }
func (a *A) Heading() polar.V         { return a.heading.V() }

func (a *A) SetID(x id.ID)                { a.id = x }
func (a *A) SetPosition(v vector.V)       { a.position.Copy(v) }
func (a *A) SetTargetPosition(v vector.V) { a.targetPosition.Copy(v) }
func (a *A) SetVelocity(v vector.V)       { a.velocity.Copy(v) }
func (a *A) SetTargetVelocity(v vector.V) { a.targetVelocity.Copy(v) }
func (a *A) SetHeading(v polar.V)         { a.heading.Copy(v) }

func (a *A) SetMoveMode(f move.F) {
	if !move.Validate(f) {
		panic(fmt.Sprintf("invalid move mode: %v", f))
	}
	a.move = f
}

func (a *A) AABB() hyperrectangle.R {
	x, y := a.position.X(), a.position.Y()
	r := a.radius

	return *hyperrectangle.New(
		vector.V{
			x - r,
			y - r,
		},
		vector.V{
			x + r,
			y + r,
		},
	)
}

func Validate(o O) bool {
	if o.Radius == 0 {
		return false
	}
	if o.Mass == 0 {
		return false
	}
	if !size.Validate(o.Size) {
		return false
	}

	return flags.Validate(o.Flags)
}
