package agent

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/flags/size"
	"github.com/downflux/go-database/flags/team"
	"github.com/downflux/go-database/internal/agent"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type O agent.O

type RO interface {
	ID() id.ID

	Position() vector.V
	TargetPosition() vector.V
	Velocity() vector.V
	TargetVelocity() vector.V
	Heading() polar.V

	Radius() float64
	Mass() float64

	MaxVelocity() float64
	MaxAngularVelocity() float64
	MaxAcceleration() float64

	Flags() flags.F
	Size() size.F
	Team() team.F

	AABB() hyperrectangle.R
}
