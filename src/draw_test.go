package semigraph

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
)

func BenchmarkRender(b *testing.B) {
	data, err := os.ReadFile("testdata/benchRGB.png")
	if err != nil {
		b.Fatal(err)
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		b.Fatal(err)
	}
	input, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(cfg.Width * cfg.Height * 4))
	b.ReportAllocs()
	for b.Loop() {
		Render(input)
	}
}

var (
	red    = color.RGBA{0xff, 0x00, 0x00, 0xff} // 48;5;196
	orange = color.RGBA{0xff, 0xa5, 0x00, 0xff} // 48;2;255;165;0
	yellow = color.RGBA{0xff, 0xff, 0x00, 0xff} // 48;5;226
	green  = color.RGBA{0x00, 0x80, 0x00, 0xff} // 48;2;0;128;0
	blue   = color.RGBA{0x00, 0x00, 0xff, 0xff} // 48;5;21
	indigo = color.RGBA{0x4b, 0x00, 0x82, 0xff} // 48;2;75;0;130
	violet = color.RGBA{0xee, 0x82, 0xee, 0xff} // 48;2;238;130;238

	rainbow = [...]color.RGBA{red, orange, yellow, green, blue, indigo, violet}
)

func TestRender(t *testing.T) {
	testCases := []struct {
		name  string
		input image.Image
		want  string
	}{
		{
			name:  "empty_image",
			input: image.NewRGBA(image.Rect(0, 0, 0, 0)),
			want:  "",
		},
		{
			name:  "img_too_small",
			input: image.NewRGBA(image.Rect(0, 0, 1, 2)),
			want:  "",
		},
		{
			name:  "all_transparent",
			input: image.NewRGBA(image.Rect(0, 0, 2, 4)),
			want:  " ",
		},
		{
			name: "alternating_block",
			input: drawFn(20, 4, func(x, _ int) color.Color {
				if x%4/2 == 0 {
					return color.Black
				}
				return color.White
			}),
			want: strings.Repeat("\x1b[48;5;16m \x1b[48;5;231m ", 5) + "\x1b[m",
		},
		{
			name: "rainbow",
			input: drawFn(14, 4, func(x, _ int) color.Color {
				return rainbow[x/2]
			}),
			want: "\x1b[48;5;196m \x1b[48;2;255;165;0m \x1b[48;5;226m \x1b[48;2;0;128;0m \x1b[48;5;21m \x1b[48;2;75;0;130m \x1b[48;2;238;130;238m \x1b[m",
		},
		{
			name: "split_block",
			input: drawFn(2, 4, func(_, y int) color.Color {
				if y < 2 {
					return color.White
				}
				return color.Black
			}),
			want: "\x1b[48;5;231;38;5;16m▄\x1b[m",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Render(tc.input)
			if got != tc.want {
				t.Errorf("Render(img) returned unexpected result:\ngot:  %q\nwant: %q", got, tc.want)
			}
		})
	}
}

func drawFn(x, y int, fn func(int, int) color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, x, y))
	for yy := range y {
		for xx := range x {
			img.Set(xx, yy, fn(xx, yy))
		}
	}
	return img
}
