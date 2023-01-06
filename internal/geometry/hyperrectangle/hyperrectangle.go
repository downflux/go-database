package hyperrectangle

import (
	"fmt"

	"github.com/downflux/go-geometry/2d/hypersphere"
	"github.com/downflux/go-geometry/2d/line"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

type Side uint64

const (
	SideN = 1 << iota
	SideE
	SideS
	SideW

	CornerNE = SideN | SideE
	CornerSE = SideS | SideE
	CornerSW = SideS | SideW
	CornerNW = SideN | SideW
)

// Normal finds the appropriate normal vector of the hyperrectangle which is
// closest to the input vector v.
//
// TODO(minkezhang): Return distance as well.
func Normal(r hyperrectangle.R, v vector.V) vector.V {
	vx, vy := v.X(), v.Y()

	xmin, xmax := r.Min().X(vnd.AXIS_X), r.Max().X(vnd.AXIS_X)
	ymin, ymax := r.Min().X(vnd.AXIS_Y), r.Max().X(vnd.AXIS_Y)

	var domain Side
	if dnorth := vy - ymax; dnorth >= 0 {
		domain |= SideN
	}
	if dsouth := ymin - vy; dsouth >= 0 {
		domain |= SideS
	}
	if deast := vx - xmax; deast >= 0 {
		domain |= SideE
	}
	if dwest := xmin - vx; dwest >= 0 {
		domain |= SideW
	}

	n := vector.M{0, 0}
	n.Copy(v)

	switch domain {
	case SideN:
		return vector.V{0, 1}
	case SideE:
		return vector.V{1, 0}
	case SideS:
		return vector.V{0, -1}
	case SideW:
		return vector.V{-1, 0}
	case CornerNE:
		n.Sub(vector.V{xmax, ymax})
		if epsilon.Within(vector.Magnitude(n.V()), 0) {
			n.Copy(vector.V{1, 1})
		}
	case CornerSE:
		n.Sub(vector.V{xmax, ymin})
		if epsilon.Within(vector.Magnitude(n.V()), 0) {
			n.Copy(vector.V{1, -1})
		}
	case CornerSW:
		n.Sub(vector.V{xmin, ymin})
		if epsilon.Within(vector.Magnitude(n.V()), 0) {
			n.Copy(vector.V{-1, -1})
		}
	case CornerNW:
		n.Sub(vector.V{xmin, ymax})
		if epsilon.Within(vector.Magnitude(n.V()), 0) {
			n.Copy(vector.V{-1, 1})
		}
	default:
		panic(fmt.Sprintf("invalid domain: %v", domain))
	}

	n.Unit()
	return n.V()
}

// CollideCircle checks if a circle overlaps an AABB. Note that this can be
// decomposed into three checks --
//
//  1. if the circle center lies inside the rectangle,
//  1. if the an edge of the rectangle crosses the circle at some point, and
//  1. if the rectangle lies entirely within the circle, but does not overlap the
//     circle center
//
// See https://stackoverflow.com/a/402019/873865 for more information.
func CollideCircle(r hyperrectangle.R, p vector.V, radius float64) bool {
	if r.In(vnd.V(p)) {
		return true
	}

	c := *hypersphere.New(p, radius)

	xmin, ymin := r.Min().X(vnd.AXIS_X), r.Min().X(vnd.AXIS_Y)
	xmax, ymax := r.Max().X(vnd.AXIS_X), r.Max().X(vnd.AXIS_Y)

	// Check corners.
	if c.In(vector.V{xmin, ymin}) {
		return true
	}
	if c.In(vector.V{xmin, ymax}) {
		return true
	}
	if c.In(vector.V{xmax, ymax}) {
		return true
	}
	if c.In(vector.V{xmax, ymin}) {
		return true
	}

	// Check edges.
	if _, _, ok := line.New(vector.V{xmin, ymin}, vector.V{0, 1}).IntersectCircle(c); ok {
		return true
	}
	if _, _, ok := line.New(vector.V{xmin, ymax}, vector.V{1, 0}).IntersectCircle(c); ok {
		return true
	}
	if _, _, ok := line.New(vector.V{xmax, ymax}, vector.V{0, -1}).IntersectCircle(c); ok {
		return true
	}
	if _, _, ok := line.New(vector.V{xmax, ymin}, vector.V{-1, 0}).IntersectCircle(c); ok {
		return true
	}

	return false
}
