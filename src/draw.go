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
	at := NewColorAtFunc(img)
	minx, miny := img.Bounds().Min.X, img.Bounds().Min.Y

	var out strings.Builder
	for ty := range h {
		ok := false
		for tx := range w {
			fg, bg, r := quantize(tx, ty, minx, miny, at)
			if !fg.alpha || !bg.alpha {
				ok = true
			}
			WriteStyled(&out, fg, bg)
			out.WriteRune(r)
		}
		if ok {
			// Only write the reset sequence if we wrote color in the first place.
			out.WriteString("\x1b[m")
		}
		if ty+1 < h {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func quantize(x, y, minx, miny int, at ColorAtFunc) (fg, bg Color, contents rune) {
	cs := make([]Color, 8)
	var rmin, gmin, bmin uint8 = 255, 255, 255
	var rmax, gmax, bmax uint8
	for i := range 8 {
		srcx := x*2 + i%2 + minx
		srcy := y*4 + i/2 + miny
		c := at(srcx, srcy)
		c.idx = i
		cs[i] = c
		rmin, rmax = min(rmin, c.R), max(rmax, c.R)
		gmin, gmax = min(gmin, c.G), max(gmax, c.G)
		bmin, bmax = min(bmin, c.B), max(bmax, c.B)
	}
	rRange := rmax - rmin
	gRange := gmax - gmin
	bRange := bmax - bmin
	// All 8 pixels are the same color.
	if rRange+gRange+bRange == 0 {
		return Transparent, cs[0], ' '
	}
	switch max(rRange, gRange, bRange) {
	case rRange:
		slices.SortFunc(cs, sortR)
	case gRange:
		slices.SortFunc(cs, sortG)
	case bRange:
		slices.SortFunc(cs, sortB)
	}
	avgA, avgB := Average(cs[:4]), Average(cs[4:])

	var mask uint8
	for _, c := range cs[:4] {
		mask |= 1 << c.idx
	}
	return avgA, avgB, blocks[mask]
}

func sortR(a, b Color) int {
	return cmp.Compare(a.R, b.R)
}

func sortG(a, b Color) int {
	return cmp.Compare(a.G, b.G)
}

func sortB(a, b Color) int {
	return cmp.Compare(a.B, b.B)
}
