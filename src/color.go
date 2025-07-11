package semigraph

import (
	"fmt"
	"image/color"
	"math"
)

// Average returns a color representing the average of colors.
// TODO: Remove all of the L*a*b* nonsense and just use RGB instead.
func Average(colors []Color) Color {
	switch len(colors) {
	case 0:
		return Color{}
	case 1:
		return colors[0]
	}

	var rsum, gsum, bsum float64
	for _, c := range colors {
		rsum += toLinear(c.R)
		gsum += toLinear(c.G)
		bsum += toLinear(c.B)
	}
	n := float64(len(colors))
	c := Color{
		R: fromLinear(rsum / n),
		G: fromLinear(gsum / n),
		B: fromLinear(bsum / n),
	}
	return c
}

// NewColor creates a Color from a [color.Color].
func NewColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	if a == 0xffff {
		return Color{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
	}
	if a == 0 {
		return Color{}
	}
	// Get the non-alpha-premultiplied channels.
	r = (r * 0xffff) / a
	g = (g * 0xffff) / a
	b = (b * 0xffff) / a
	return Color{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

// Color is a color.
type Color struct {
	R, G, B uint8
}

func (c Color) Foreground() string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

func (c Color) Background() string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
}

// toLinear converts an sRGB channel to linear RGB.
// https://en.wikipedia.org/wiki/SRGB#Transfer_function_(%22gamma%22)
func toLinear(c uint8) float64 {
	v := float64(c) / 255
	if v < 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// fromLinear converts a linear RGB channel to sRGB.
func fromLinear(c float64) uint8 {
	if c <= 0.0031308 {
		c = 12.92 * c
	} else {
		c = 1.055*math.Pow(c, 1.0/2.4) - 0.055
	}
	v := uint8(c * 255)
	return max(0, min(v, 255))
}
