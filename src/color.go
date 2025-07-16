package semigraph

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
)

var toLinearLUT [256]float64

func init() {
	for i := range 256 {
		v := float64(i) / 255
		if v < 0.04045 {
			toLinearLUT[i] = v / 12.92
		} else {
			toLinearLUT[i] = math.Pow((v+0.055)/1.055, 2.4)
		}
	}
}

// Color is a color.
type Color struct {
	R, G, B uint8

	// This is used to preserve pixel location when quantizing.
	// It lives here to prevent allocating memory to store it
	// elsewhere.
	idx int
}

func (c Color) WriteANSI(buf *strings.Builder) {
	buf.WriteString(colorLUT[c.R])
	buf.WriteByte(';')
	buf.WriteString(colorLUT[c.G])
	buf.WriteByte(';')
	buf.WriteString(colorLUT[c.B])
}

func (c Color) Show() string {
	return fmt.Sprintf("\x1b[48;2;%s;%s;%sm  \x1b[m", colorLUT[c.R], colorLUT[c.G], colorLUT[c.B])
}

var colorLUT = [256]string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
	"20", "21", "22", "23", "24", "25", "26", "27", "28", "29",
	"30", "31", "32", "33", "34", "35", "36", "37", "38", "39",
	"40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
	"50", "51", "52", "53", "54", "55", "56", "57", "58", "59",
	"60", "61", "62", "63", "64", "65", "66", "67", "68", "69",
	"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
	"100", "101", "102", "103", "104", "105", "106", "107", "108", "109",
	"110", "111", "112", "113", "114", "115", "116", "117", "118", "119",
	"120", "121", "122", "123", "124", "125", "126", "127", "128", "129",
	"130", "131", "132", "133", "134", "135", "136", "137", "138", "139",
	"140", "141", "142", "143", "144", "145", "146", "147", "148", "149",
	"150", "151", "152", "153", "154", "155", "156", "157", "158", "159",
	"160", "161", "162", "163", "164", "165", "166", "167", "168", "169",
	"170", "171", "172", "173", "174", "175", "176", "177", "178", "179",
	"180", "181", "182", "183", "184", "185", "186", "187", "188", "189",
	"190", "191", "192", "193", "194", "195", "196", "197", "198", "199",
	"200", "201", "202", "203", "204", "205", "206", "207", "208", "209",
	"210", "211", "212", "213", "214", "215", "216", "217", "218", "219",
	"220", "221", "222", "223", "224", "225", "226", "227", "228", "229",
	"230", "231", "232", "233", "234", "235", "236", "237", "238", "239",
	"240", "241", "242", "243", "244", "245", "246", "247", "248", "249",
	"250", "251", "252", "253", "254", "255",
}

type ColorAtFunc func(x, y int) Color

// NewColorAtFunc returns a function which returns the color at a given
// set of coordinates in an image.
//
// This is faster than using the `At` func in go's image.Image package.
// See https://github.com/golang/go/issues/15759.
func NewColorAtFunc(img image.Image) ColorAtFunc {
	switch p := img.(type) {
	case *image.RGBA:
		return newColorAtFuncRGBA(p)
	case *image.NRGBA:
		return newColorAtFuncNRGBA(p)
	case *image.YCbCr:
		return newColorAtFuncYCbCr(p)
	default:
		panic(fmt.Sprintf("unsupported image color mode %T", img))
	}
}

func newColorAtFuncRGBA(p *image.RGBA) ColorAtFunc {
	return func(x, y int) Color {
		c := p.RGBAAt(x, y)
		if c.A == 0x00 {
			return Color{}
		}
		if c.A == 0xff {
			return Color{c.R, c.G, c.B, 0}
		}
		return Color{
			uint8(uint32(c.R*c.A) >> 8),
			uint8(uint32(c.G*c.A) >> 8),
			uint8(uint32(c.B*c.A) >> 8),
			0,
		}
	}
}

func newColorAtFuncNRGBA(p *image.NRGBA) ColorAtFunc {
	return func(x, y int) Color {
		c := p.NRGBAAt(x, y)
		if c.A == 0x00 {
			return Color{}
		}
		if c.A == 0xff {
			return Color{c.R, c.G, c.B, 0}
		}
		rr, gg, bb, aa := color.NRGBA{c.R, c.G, c.B, c.A}.RGBA()
		return Color{
			uint8(rr * aa >> 8),
			uint8(gg * aa >> 8),
			uint8(bb * aa >> 8),
			0,
		}
	}
}

func newColorAtFuncYCbCr(p *image.YCbCr) ColorAtFunc {
	return func(x, y int) Color {
		c := p.YCbCrAt(x, y)
		r, g, b := color.YCbCrToRGB(c.Y, c.Cb, c.Cr)
		return Color{r, g, b, 0}
	}
}

// Average returns a color representing the average of colors.
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
	r := fromLinear(rsum / n)
	g := fromLinear(gsum / n)
	b := fromLinear(bsum / n)
	return Color{R: r, G: g, B: b}
}

// toLinear converts an sRGB channel to linear RGB.
// https://en.wikipedia.org/wiki/SRGB#Transfer_function_(%22gamma%22)
func toLinear(c uint8) float64 {
	return toLinearLUT[c]
}

// fromLinear converts a linear RGB channel to sRGB.
func fromLinear(c float64) uint8 {
	closest := 0
	minDelta := math.MaxFloat64
	for i := range 256 {
		delta := c - toLinearLUT[i]
		delta = max(delta, -delta)
		if delta < minDelta {
			minDelta = delta
			closest = i
		}
	}
	return uint8(closest)
}
