package semigraph

import (
	"cmp"
	"image"
	"math"
	"slices"
	"strings"
)

type pixel struct {
	idx   int
	color Color
}

// Render renders the img using semigraphic characters and ANSI escapes.
func Render(img image.Image) string {
	srcw := img.Bounds().Dx()
	srch := img.Bounds().Dy()
	w := int(math.Floor(float64(srcw) / 2))
	h := int(math.Floor(float64(srch) / 4))

	var out strings.Builder

	for ty := range h {
		for tx := range w {
			fg, bg, r := quantize(tx, ty, img)
			out.WriteString(fg.Foreground())
			out.WriteString(bg.Background())
			out.WriteRune(r)
			out.WriteString("\x1b[0m")
		}
		out.WriteByte('\n')
	}

	return out.String()
}

func quantize(x, y int, img image.Image) (fg, bg Color, contents rune) {
	cs := make([]pixel, 8)
	var rmin, rmax, gmin, gmax, bmin, bmax uint8
	for i := range 8 {
		srcx := x*2 + i%2 + img.Bounds().Min.X
		srcy := y*4 + i/2 + img.Bounds().Min.Y
		c := NewColor(img.At(srcx, srcy))
		rmin, rmax = min(rmin, c.R), max(rmax, c.R)
		gmin, gmax = min(gmin, c.G), max(gmax, c.G)
		bmin, bmax = min(bmin, c.B), max(bmax, c.B)
		cs[i] = pixel{i, c}
	}
	rRange := rmax - rmin
	gRange := gmax - gmin
	bRange := bmax - bmin
	var fn func(pixel, pixel) int
	switch max(rRange, gRange, bRange) {
	case rRange:
		fn = sortR
	case gRange:
		fn = sortG
	case bRange:
		fn = sortB
	}
	slices.SortFunc(cs, fn)
	avgA, avgB := average(cs[:4]), average(cs[4:])

	var mask uint8
	for _, c := range cs[:4] {
		mask |= 1 << c.idx
	}
	return avgA, avgB, blocks[mask]
}

func sortR(a, b pixel) int {
	return cmp.Compare(a.color.R, b.color.R)
}

func sortG(a, b pixel) int {
	return cmp.Compare(a.color.G, b.color.G)
}

func sortB(a, b pixel) int {
	return cmp.Compare(a.color.B, b.color.B)
}

func average(colors []pixel) Color {
	cs := make([]Color, len(colors))
	for i, c := range colors {
		cs[i] = c.color
	}
	return Average(cs)
}
