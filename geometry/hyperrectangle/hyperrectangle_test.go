package hyperrectangle

import (
	"testing"

	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	vnd "github.com/downflux/go-geometry/nd/vector"
)

func TestNormal(t *testing.T) {
	r := *hyperrectangle.New(vnd.V{0, 0}, vnd.V{10, 10})
	type config struct {
		name string
		v    vector.V
		want vector.V
	}

	configs := []config{
		{
			name: "North",
			v:    vector.V{5, 20},
			want: vector.V{0, 1},
		},
		{
			name: "South",
			v:    vector.V{5, -10},
			want: vector.V{0, -1},
		},
		{
			name: "East",
			v:    vector.V{20, 5},
			want: vector.V{1, 0},
		},
		{
			name: "West",
			v:    vector.V{-10, 5},
			want: vector.V{-1, 0},
		},

		{
			name: "Corner/NE",
			v:    vector.V{10, 10},
			want: vector.Unit(vector.V{1, 1}),
		},
		{
			name: "Corner/SE",
			v:    vector.V{10, 0},
			want: vector.Unit(vector.V{1, -1}),
		},
		{
			name: "Corner/SW",
			v:    vector.V{0, 0},
			want: vector.Unit(vector.V{-1, -1}),
		},
		{
			name: "Corner/NW",
			v:    vector.V{0, 10},
			want: vector.Unit(vector.V{-1, 1}),
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Normal(r, c.v); !vector.Within(got, c.want) {
				t.Errorf("Normal() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestIntersectCircle(t *testing.T) {
	type config struct {
		name   string
		r      hyperrectangle.R
		p      vector.V
		radius float64
		want   bool
	}

	configs := []config{
		{
			name:   "Center",
			r:      *hyperrectangle.New(vnd.V{0, 0}, vnd.V{10, 10}),
			p:      vector.V{5, 5},
			radius: 1,
			want:   true,
		},
		{
			name:   "Corner",
			r:      *hyperrectangle.New(vnd.V{0, 0}, vnd.V{10, 10}),
			p:      vector.V{-1, -1},
			radius: 2,
			want:   true,
		},
		{
			name:   "Edge",
			r:      *hyperrectangle.New(vnd.V{0, 0}, vnd.V{10, 10}),
			p:      vector.V{-1, 5},
			radius: 2,
			want:   true,
		},
		{
			name:   "Outside",
			r:      *hyperrectangle.New(vnd.V{0, 0}, vnd.V{10, 10}),
			p:      vector.V{12, 12},
			radius: 1,
			want:   false,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := IntersectCircle(c.r, c.p, c.radius); got != c.want {
				t.Errorf("IntersectCircle() = %v, want = %v", got, c.want)
			}
		})
	}
}
