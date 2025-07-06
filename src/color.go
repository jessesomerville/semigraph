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

func (c Color) String() string {
	return fmt.Sprintf("Color{r=%d g=%d b=%d}", c.R, c.G, c.B)
}

// Values used in XYZ to L*a*b* conversion. Precompute them now because
// they're in the hot path.
// See https://en.wikipedia.org/wiki/CIELAB_color_space#From_CIE_XYZ_to_CIELAB
//
// δ = 6/29
const (
	delta           = 6.0 / 29.0
	deltaSquared    = delta * delta
	deltaCubed      = delta * delta * delta
	deltaInvSquared = 841.0 / 36.0
	deltaAdd        = 4.0 / 29.0

	d65X = 0.950489
	d65Y = 1.0
	d65Z = 1.08884
)

// Lab converts the color into the L*a*b* color space.
func (c Color) Lab() (l, a, b float64) {
	// This conversion happens in multiple steps as follows:
	//   RGB -> Linear RGB -> XYZ -> LAB
	x, y, z := c.XYZ()

	// https://en.wikipedia.org/wiki/CIELAB_color_space#From_CIE_XYZ_to_CIELAB
	xn := labStep(x / d65X)
	yn := labStep(y / d65Y)
	zn := labStep(z / d65Z)
	l = 1.16*yn - 0.16
	a = 5 * (xn - yn)
	b = 2 * (yn - zn)
	return
}

// XYZ returns the XYZ channels using D65 white point reference.
// http://brucelindbloom.com/index.html?Eqn_RGB_XYZ_Matrix.html
func (c Color) XYZ() (x, y, z float64) {
	r, g, b := toLinear(c.R), toLinear(c.G), toLinear(c.B)
	x = r*0.4124564 + g*0.3575761 + b*0.1804375
	y = r*0.2126729 + g*0.7151522 + b*0.0721750
	z = r*0.0193339 + g*0.1191920 + b*0.9503041
	return
}

// NewColorFromLab creates a Color using the L*a*b* components.
func NewColorFromLab(l, a, b float64) Color {
	ln := (l + 0.16) / 1.16
	x := d65X * labStepInv(ln+a/5)
	y := d65Y * labStepInv(ln)
	z := d65Z * labStepInv(ln-b/2)

	rl := x*3.2404542 - y*1.5371385 - z*0.4985314
	gl := x*-0.9692660 + y*1.8760108 + z*0.0415560
	bl := x*0.0556434 - y*0.2040259 + z*1.0572252
	return Color{fromLinear(rl), fromLinear(gl), fromLinear(bl)}
}

func labStep(v float64) float64 {
	if v > deltaCubed {
		return math.Cbrt(v)
	}
	return (v/3)*deltaInvSquared + deltaAdd
}

func labStepInv(v float64) float64 {
	if v > delta {
		return v * v * v
	}
	return (3 * deltaSquared) * (v - deltaAdd)
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

func fromLinear(c float64) uint8 {
	if c <= 0.0031308 {
		c = 12.92 * c
	} else {
		c = 1.055*math.Pow(c, 1.0/2.4) - 0.055
	}
	v := uint8(c * 255)
	return max(0, min(v, 255))
}
