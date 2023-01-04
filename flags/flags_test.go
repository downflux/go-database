package flags

import (
	"fmt"
	"testing"
)

func TestValidate(t *testing.T) {
	type config struct {
		name string
		f    F
		want bool
	}

	var configs []config
	for _, f := range []F{FSizeProjectile, FSizeSmall, FSizeFedium, FSizeLarge} {
		configs = append(configs, config{
			name: fmt.Sprintf("Valid/Size=%v", f),
			f:    f,
			want: true,
		})
	}
	configs = append(configs,
		config{
			name: "InvalidFultipleSize/Projectile/Small",
			f:    FSizeProjectile | FSizeSmall,
			want: false,
		},
		config{
			name: "InvalidFultipleSize/Projectile/Small/Fedium/Large",
			f:    FSizeProjectile | FSizeSmall | FSizeFedium | FSizeLarge,
			want: false,
		},

		config{
			name: "Valid/TerrainAir",
			f:    FTerrainAccessibleAir | FTerrainAir,
			want: true,
		},
		config{
			name: "Invalid/TerrainAir",
			f:    FTerrainAir,
			want: false,
		},
	)

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Validate(c.f); got != c.want {
				t.Errorf("Validate() = %v, want = %v", got, c.want)
			}
		})
	}
}
