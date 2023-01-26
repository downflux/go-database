package hyperrectangle

import (
	"fmt"
	"math"

	"github.com/downflux/go-geometry/2d/hyperrectangle"
	"github.com/downflux/go-geometry/2d/hypersphere"
	"github.com/downflux/go-geometry/2d/line"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/epsilon"
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
// closest to the input vector v. Also returns the distance to the corresponding
// edge or corner.
func Normal(r hyperrectangle.R, v vector.V) (float64, vector.V) {
	vx, vy := v.X(), v.Y()

	var d float64
	xmin, xmax := r.Min().X(), r.Max().X()
	ymin, ymax := r.Min().Y(), r.Max().Y()

	var domain Side
	if dnorth := vy - ymax; dnorth >= 0 {
		d += dnorth * dnorth
		domain |= SideN
	}
	if dsouth := ymin - vy; dsouth >= 0 {
		d += dsouth * dsouth
		domain |= SideS
	}
	if deast := vx - xmax; deast >= 0 {
		d += deast * deast
		domain |= SideE
	}
	if dwest := xmin - vx; dwest >= 0 {
		d += dwest * dwest
		domain |= SideW
	}

	d = math.Sqrt(d)

	n := vector.M{0, 0}
	n.Copy(v)

	switch domain {
	case SideN:
		return d, vector.V{0, 1}
	case SideE:
		return d, vector.V{1, 0}
	case SideS:
		return d, vector.V{0, -1}
	case SideW:
		return d, vector.V{-1, 0}
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
	return d, n.V()
}

// IntersectCircle checks if a circle overlaps an AABB. Note that this can be
// decomposed into three checks --
//
//  1. if the circle center lies inside the rectangle,
//  1. if the an edge of the rectangle crosses the circle at some point, and
//  1. if the rectangle lies entirely within the circle, but does not overlap the
//     circle center
//
// See https://stackoverflow.com/a/402019/873865 for more information.
func IntersectCircle(r hyperrectangle.R, p vector.V, radius float64) bool {
	if r.In(vector.V(p)) {
		return true
	}

	c := *hypersphere.New(p, radius)

	xmin, ymin := r.Min().X(), r.Min().Y()
	xmax, ymax := r.Max().X(), r.Max().Y()

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
