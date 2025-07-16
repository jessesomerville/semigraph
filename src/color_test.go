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
			input: []Color{{1, 2, 3, 0}},
			want:  Color{1, 2, 3, 0},
		},
		{
			name: "black_and_white",
			input: []Color{
				{0, 0, 0, 0},
				{255, 255, 255, 0},
			},
			want: Color{188, 188, 188, 0},
		},
		{
			name: "red_and_blue",
			input: []Color{
				{255, 0, 0, 0},
				{0, 0, 255, 0},
			},
			want: Color{188, 0, 188, 0},
		},
		{
			name: "rainbow",
			input: []Color{
				{255, 0, 0, 0},
				{255, 128, 0, 0},
				{128, 128, 0, 0},
				{128, 255, 0, 0},
				{0, 255, 0, 0},
				{0, 255, 128, 0},
				{0, 255, 255, 0},
				{0, 128, 128, 0},
				{0, 0, 128, 0},
			},
			want: Color{142, 190, 119, 0},
		},
	}

	for _, tc := range testCases {
		got := Average(tc.input)
		if got != tc.want {
			t.Errorf("%s: Average(%+v) = %v, want %v", tc.name, tc.input, got, tc.want)
		}
	}
}
