package semigraph

import (
	"cmp"
	"image"
	"image/draw"
	"slices"

	"github.com/lucasb-eyer/go-colorful"
)

type qColor struct {
	idx image.Point
	c   colorful.Color
}

// Quantize2 quantizes src into two colors and draws it to dst.
func Quantize2(dst draw.Image, src image.Image) {
	bounds := src.Bounds()
	cs := make([]qColor, 0, bounds.Dx()*bounds.Dy())

	// Convert to linear RGB and calculate the channel ranges.
	var rmin, rmax, gmin, gmax, bmin, bmax float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			v := src.At(x, y)
			c, _ := colorful.MakeColor(v)
			rmin, rmax = min(rmin, c.R), max(rmax, c.R)
			gmin, gmax = min(gmin, c.G), max(gmax, c.G)
			bmin, bmax = min(bmin, c.B), max(bmax, c.B)
			cs = append(cs, qColor{image.Pt(x, y), c})
		}
	}
	rRange := rmax - rmin
	gRange := gmax - gmin
	bRange := bmax - bmin

	switch max(rRange, gRange, bRange) {
	case rRange: // Sort using the red channel
		slices.SortFunc(cs, func(a, b qColor) int {
			return cmp.Compare(a.c.R, b.c.R)
		})
	case gRange: // Sort using the green channel
		slices.SortFunc(cs, func(a, b qColor) int {
			return cmp.Compare(a.c.G, b.c.G)
		})
	case bRange: // Sort using the blue channel
		slices.SortFunc(cs, func(a, b qColor) int {
			return cmp.Compare(a.c.B, b.c.B)
		})
	}

	n := len(cs) >> 1
	avgA, avgB := averageQ(cs[:n]), averageQ(cs[n:])
	for _, c := range cs[:n] {
		dst.Set(c.idx.X, c.idx.Y, avgA)
	}
	for _, c := range cs[n:] {
		dst.Set(c.idx.X, c.idx.Y, avgB)
	}
}

func averageQ(colors []qColor) colorful.Color {
	if len(colors) == 0 {
		return colorful.Color{}
	}

	var nl, na, nb float64
	for _, c := range colors {
		l, a, b := c.c.Lab()
		nl += l
		na += a
		nb += b
	}
	n := float64(len(colors))
	return colorful.Lab(nl/n, na/n, nb/n)
}
