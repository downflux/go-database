package feature

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-data/flags"
	"github.com/downflux/go-data/internal/feature"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type O feature.O

type RO interface {
	ID() id.ID
	AABB() hyperrectangle.R
	Flags() flags.F
}
