package flags

import (
	"testing"
)

func TestValidate(t *testing.T) {
	type config struct {
		name string
		f    F
		want bool
	}

	configs := []config{
		{
			name: "Valid/TerrainAir",
			f:    FTerrainAccessibleAir | FTerrainAir,
			want: true,
		},
		{
			name: "Invalid/TerrainAir",
			f:    FTerrainAir,
			want: false,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Validate(c.f); got != c.want {
				t.Errorf("Validate() = %v, want = %v", got, c.want)
			}
		})
	}
}
