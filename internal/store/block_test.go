package store

import "testing"

func TestOverlaps(t *testing.T) {
	cases := []struct {
		name string
		a, b Block
		want bool
	}{
		{
			name: "overlapping",
			a:    Block{Date: "2026-09-29", Start: 540, End: 600},
			b:    Block{Date: "2026-09-29", Start: 570, End: 630},
			want: true,
		},
		{
			name: "touching, not overlapping",
			a:    Block{Date: "2026-09-29", Start: 540, End: 600},
			b:    Block{Date: "2026-09-29", Start: 600, End: 660},
			want: false,
		},
		{
			name: "entirely inside the other, same date",
			a:    Block{Date: "2026-09-29", Start: 540, End: 660},
			b:    Block{Date: "2026-09-29", Start: 570, End: 630},
			want: true,
		},
		{
			name: "identical times, different date",
			a:    Block{Date: "2026-09-28", Start: 540, End: 660},
			b:    Block{Date: "2026-09-29", Start: 540, End: 660},
			want: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.a.Overlaps(c.b)
			if got != c.want {
				t.Errorf("Overlaps() = %v, want %v", got, c.want)
			}
		})
	}
}
