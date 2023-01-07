package feature

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-database/internal/feature"
	"github.com/downflux/go-database/team"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type O feature.O

type RO interface {
	ID() id.ID

	Flags() flags.F
	Team() team.T

	AABB() hyperrectangle.R
}
