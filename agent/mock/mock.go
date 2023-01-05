package mock

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/internal/agent"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	roagent "github.com/downflux/go-database/agent"
)

var (
	_ roagent.RO = &A{}
)

type A agent.A

func New(x id.ID, o roagent.O) *A {
	if o.Position == nil {
		(&o).Position = vector.V{0, 0}
	}
	if o.TargetPosition == nil {
		(&o).TargetPosition = vector.V{0, 0}
	}
	if o.Velocity == nil {
		(&o).Velocity = vector.V{0, 0}
	}
	if o.TargetVelocity == nil {
		(&o).TargetVelocity = vector.V{0, 0}
	}
	if o.Heading == nil {
		(&o).Heading = polar.V{0, 0}
	}

	a := agent.New(agent.O(o))
	a.SetID(x)
	return (*A)(a)
}

func (a *A) ID() id.ID                   { return (*agent.A)(a).ID() }
func (a *A) Position() vector.V          { return (*agent.A)(a).Position() }
func (a *A) TargetPosition() vector.V    { return (*agent.A)(a).TargetPosition() }
func (a *A) Velocity() vector.V          { return (*agent.A)(a).Velocity() }
func (a *A) TargetVelocity() vector.V    { return (*agent.A)(a).TargetVelocity() }
func (a *A) Heading() polar.V            { return (*agent.A)(a).Heading() }
func (a *A) Radius() float64             { return (*agent.A)(a).Radius() }
func (a *A) Mass() float64               { return (*agent.A)(a).Mass() }
func (a *A) MaxVelocity() float64        { return (*agent.A)(a).MaxVelocity() }
func (a *A) MaxAngularVelocity() float64 { return (*agent.A)(a).MaxAngularVelocity() }
func (a *A) MaxAcceleration() float64    { return (*agent.A)(a).MaxAcceleration() }
func (a *A) Flags() flags.F              { return (*agent.A)(a).Flags() }
func (a *A) AABB() hyperrectangle.R      { return (*agent.A)(a).AABB() }
