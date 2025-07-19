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

var Transparent = Color{alpha: true}

// Color is a color.
type Color struct {
	// Red, Green, and Blue channels.
	R, G, B uint8

	alpha bool

	// This is used to preserve pixel location when quantizing.
	// It lives here to prevent allocating memory to store it
	// elsewhere.
	idx int
}

func RGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

// to8bit returns the color's corresponding code from the 6x6x6 color cube
// defined by the range [16,231] and whether that conversion was successful or
// not.
//
// See https://en.wikipedia.com/wiki/ANSI_escape_code#8-bit.
func to8bit(c Color) (uint8, bool) {
	clamp := func(v uint8) (uint8, bool) {
		if v == 0 {
			return 0, true
		}
		if v < 0x5f {
			return 0, false
		}
		v -= 55
		if v%40 == 0 {
			return v / 40, true
		}
		return 0, false
	}

	r, rok := clamp(c.R)
	g, gok := clamp(c.G)
	b, bok := clamp(c.B)
	if !(rok && gok && bok) {
		return 0, false
	}
	return 16 + (r * 36) + (g * 6) + b, true
}

func writeANSI(buf *strings.Builder, c Color) {
	if c.alpha {
		return
	}

	if v, ok := to8bit(c); ok {
		buf.WriteString("5;")
		buf.WriteString(colorLUT[v])
		return
	}
	buf.WriteString("2;")
	buf.WriteString(colorLUT[c.R])
	buf.WriteByte(';')
	buf.WriteString(colorLUT[c.G])
	buf.WriteByte(';')
	buf.WriteString(colorLUT[c.B])
}

// WriteStyled writes the foreground and background ANSI sequences to w.
//
// If neither color is [Transparent], the written escape sequence is in the
// following format:
//
//	\x1b[48;<bg>;38;<fg>m
//
// Where <bg> and <fg> are either 8-bit or 24-bit color codes. The escape
// sequence for a color will be omitted if that color is [Transparent], and
// calling WriteStyled with both fg and bg == Transparent is a no-op.
func WriteStyled(buf *strings.Builder, fg, bg Color) {
	if fg.alpha && bg.alpha {
		return
	}
	buf.WriteString("\x1b[")
	if !bg.alpha {
		buf.WriteString("48;")
		writeANSI(buf, bg)
	}
	if !fg.alpha {
		if !bg.alpha {
			buf.WriteByte(';')
		}
		buf.WriteString("38;")
		writeANSI(buf, fg)
	}
	buf.WriteByte('m')
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
			return Transparent
		}
		if c.A == 0xff {
			return RGB(c.R, c.G, c.B)
		}
		return RGB(
			uint8(uint32(c.R*c.A)>>8),
			uint8(uint32(c.G*c.A)>>8),
			uint8(uint32(c.B*c.A)>>8),
		)
	}
}

func newColorAtFuncNRGBA(p *image.NRGBA) ColorAtFunc {
	return func(x, y int) Color {
		c := p.NRGBAAt(x, y)
		if c.A == 0x00 {
			return Transparent
		}
		if c.A == 0xff {
			return RGB(c.R, c.G, c.B)
		}
		rr, gg, bb, aa := color.NRGBA{c.R, c.G, c.B, c.A}.RGBA()
		return RGB(
			uint8(rr*aa>>8),
			uint8(gg*aa>>8),
			uint8(bb*aa>>8),
		)
	}
}

func newColorAtFuncYCbCr(p *image.YCbCr) ColorAtFunc {
	return func(x, y int) Color {
		c := p.YCbCrAt(x, y)
		r, g, b := color.YCbCrToRGB(c.Y, c.Cb, c.Cr)
		return RGB(r, g, b)
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
	if c <= 0.0031308 {
		c = 12.92 * c
	} else {
		c = 1.055*math.Pow(c, 1.0/2.4) - 0.055
	}
	v := uint8(c * 255)
	return max(0, min(v, 255))
}
