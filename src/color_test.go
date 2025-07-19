package semigraph

import (
	"testing"
)

func TestColorAverage(t *testing.T) {
	testCases := []struct {
		name  string
		input []Color
		want  Color
	}{
		{
			name:  "no_colors",
			input: []Color{},
			want:  Color{},
		},
		{
			name:  "one_color",
			input: []Color{RGB(1, 2, 3)},
			want:  RGB(1, 2, 3),
		},
		{
			name: "black_and_white",
			input: []Color{
				RGB(0, 0, 0),
				RGB(255, 255, 255),
			},
			want: RGB(187, 187, 187),
		},
		{
			name: "red_and_blue",
			input: []Color{
				RGB(255, 0, 0),
				RGB(0, 0, 255),
			},
			want: RGB(187, 0, 187),
		},
		{
			name: "rainbow",
			input: []Color{
				RGB(255, 0, 0),
				RGB(255, 128, 0),
				RGB(128, 128, 0),
				RGB(128, 255, 0),
				RGB(0, 255, 0),
				RGB(0, 255, 128),
				RGB(0, 255, 255),
				RGB(0, 128, 128),
				RGB(0, 0, 128),
			},
			want: RGB(141, 190, 118),
		},
	}

	for _, tc := range testCases {
		got := Average(tc.input)
		if got != tc.want {
			t.Errorf("%s: Average(%+v) = %v, want %v", tc.name, tc.input, got, tc.want)
		}
	}
}

func TestTo8Bit(t *testing.T) {
	reds := map[uint8]uint8{
		0x00: 16,
		0x5f: 52,
		0x87: 88,
		0xaf: 124,
		0xd7: 160,
		0xff: 196,
	}
	r := uint8(0)
	for r = range 0xff {
		want, ok := reds[r]
		got, gotOK := to8bit(RGB(r, 0, 0))
		if gotOK != ok || got != want {
			t.Errorf("to8bit({%#02x, 0, 0}) = (%d, %v), want (%d, %v)", r, got, gotOK, want, ok)
		}
	}
}
