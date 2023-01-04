package projectile

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/internal/projectile"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type O projectile.O

type RO interface {
	ID() id.ID

	Position() vector.V
	Velocity() vector.V
	TargetVelocity() vector.V
	Heading() polar.V

	Radius() float64

	Flags() flags.F

	AABB() hyperrectangle.R
}
