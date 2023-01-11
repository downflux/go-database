package filters

import (
	"testing"

	"github.com/downflux/go-database/agent"
	"github.com/downflux/go-database/agent/mock"
	"github.com/downflux/go-database/flags"
	"github.com/downflux/go-geometry/2d/vector"
	"github.com/downflux/go-geometry/2d/vector/polar"
)

func TestAgentIsColliding(t *testing.T) {
	type config struct {
		name string
		a    agent.RO
		b    agent.RO
		want bool
	}

	configs := []config{
		func() config {
			a := mock.New(1, agent.O{
				Heading:  polar.V{1, 0},
				Velocity: vector.V{0, 0},
				Position: vector.V{1, 1},
			})
			return config{
				name: "NoCollide/SelfID",
				a:    a,
				b:    a,
				want: false,
			}
		}(),
		func() config {
			a := mock.New(1, agent.O{
				Position: vector.V{1, 1},
				Radius:   1,
				Velocity: vector.V{0, 0},
				Heading:  polar.V{1, 0},
				Flags:    flags.FTerrainAir | flags.FTerrainAccessibleAir,
			})
			b := mock.New(2, agent.O{
				Position: vector.V{1, 1},
				Radius:   1,
				Velocity: vector.V{0, 0},
				Heading:  polar.V{1, 0},
				Flags:    flags.FTerrainLand | flags.FTerrainAccessibleLand,
			})
			return config{
				name: "NoCollide/ExclusiveAir",
				a:    a,
				b:    b,
				want: false,
			}
		}(),
		func() config {
			a := mock.New(1, agent.O{
				Position: vector.V{1, 1},
				Radius:   1,
				Velocity: vector.V{0, 0},
				Heading:  polar.V{1, 0},
				Flags:    flags.FTerrainAir | flags.FTerrainAccessibleAir,
			})
			b := mock.New(2, agent.O{
				Position: vector.V{1, 1},
				Radius:   1,
				Velocity: vector.V{0, 0},
				Heading:  polar.V{1, 0},
				Flags:    flags.FTerrainAir | flags.FTerrainAccessibleAir,
			})
			return config{
				name: "Collide/BothAir",
				a:    a,
				b:    b,
				want: true,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := AgentIsColliding(c.a, c.b); got != c.want {
				t.Errorf("AgentIsColliding() = %v, want = %v", got, c.want)
			}
		})
	}
}
