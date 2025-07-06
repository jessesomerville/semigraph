package semigraph

import (
	"cmp"
	"fmt"
	"image"
	"math"
	"slices"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

type pixel struct {
	idx   int
	color colorful.Color
}

// Render renders the img using semigraphic characters and ANSI escapes.
func Render(img image.Image) string {
	srcw := img.Bounds().Dx()
	srch := img.Bounds().Dy()
	w := int(math.Floor(float64(srcw) / 2))
	h := int(math.Floor(float64(srch) / 4))

	var out strings.Builder

	cs := make([]pixel, 8)
	var rmin, rmax, gmin, gmax, bmin, bmax float64
	for ty := range h {
		for tx := range w {
			for i := range 8 {
				srcx := tx*2 + i/4
				srcy := ty*4 + i%4
				c, _ := colorful.MakeColor(img.At(srcx, srcy))
				cl, ca, cb := c.Lab()
				rmin, rmax = min(rmin, cl), max(rmax, cb)
				gmin, gmax = min(gmin, ca), max(gmax, ca)
				bmin, bmax = min(bmin, cb), max(bmax, cl)
				cs[i] = pixel{i, c}
			}
			rRange := rmax - rmin
			gRange := gmax - gmin
			bRange := bmax - bmin
			var fn func(pixel, pixel) int
			switch max(rRange, gRange, bRange) {
			case rRange:
				fn = sortRed
			case gRange:
				fn = sortGreen
			case bRange:
				fn = sortBlue
			}
			slices.SortFunc(cs, fn)
			avgA, avgB := average(cs[:4]), average(cs[4:])

			var mask uint8
			for _, c := range cs[:4] {
				mask |= 1 << c.idx
			}
			out.WriteString(toANSI(avgA, false))
			out.WriteString(toANSI(avgB, true))
			out.WriteRune(blocks[mask])
			out.WriteString("\x1b[0m")
		}
		out.WriteByte('\n')
	}

	return out.String()
}

func toANSI(c colorful.Color, bg bool) string {
	r, g, b := c.Clamped().RGB255()
	csi := 38
	if bg {
		csi = 48
	}
	return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm", csi, r, g, b)
}

func sortRed(a, b pixel) int {
	al, _, _ := a.color.Lab()
	bl, _, _ := b.color.Lab()
	return cmp.Compare(al, bl)
}

func sortGreen(a, b pixel) int {
	_, ag, _ := a.color.Lab()
	_, bg, _ := b.color.Lab()
	return cmp.Compare(ag, bg)
}

func sortBlue(a, b pixel) int {
	_, _, ab := a.color.Lab()
	_, _, bb := b.color.Lab()
	return cmp.Compare(ab, bb)
}

func average(colors []pixel) colorful.Color {
	if len(colors) == 0 {
		return colorful.Color{}
	}

	var nl, na, nb float64
	for _, c := range colors {
		l, a, b := c.color.Lab()
		nl += l
		na += a
		nb += b
	}
	n := float64(len(colors))
	return colorful.Lab(nl/n, na/n, nb/n).Clamped()
}
