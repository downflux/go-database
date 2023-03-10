package feature

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/flags/team"
	"github.com/downflux/go-database/internal/feature"
	"github.com/downflux/go-geometry/2d/hyperrectangle"
)

type O feature.O

type RO interface {
	ID() id.ID

	Flags() flags.F
	Team() team.F

	AABB() hyperrectangle.R
}
