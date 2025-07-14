package semigraph

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"strings"
	"time"
)

type GIF struct {
	frames []*frame
}

type frame struct {
	delay    time.Duration
	contents string
	lines    int
}

// RenderGIF parses the frames of the input GIF into a [GIF] that can be
// rendered in a terminal using [GIF.Play].
func RenderGIF(g *gif.GIF) (*GIF, error) {
	nFrames := len(g.Image)

	if nFrames == 0 {
		return nil, errors.New("semigraph: GIF has no frames")
	}
	if nFrames != len(g.Delay) || nFrames != len(g.Disposal) {
		return nil, errors.New("semigraph: mismatched GIF frame, disposal, and delay lengths")
	}

	out := &GIF{
		frames: make([]*frame, nFrames),
	}
	gBounds := image.Rect(0, 0, g.Config.Width, g.Config.Height)
	prev := image.NewRGBA(gBounds) // Might defer this in case there's only 1 frame.
	base := image.NewRGBA(gBounds)
	for i, frm := range g.Image {
		if i > 0 {
			draw.Draw(base, gBounds, prev, image.Point{}, draw.Src)
			draw.Draw(base, gBounds, frm, image.Point{}, draw.Over)
		} else {
			draw.Draw(base, frm.Bounds(), frm, image.Point{}, draw.Src)
		}
		contents := Render(base)
		prev.Pix = clonePix(base.Pix)
		clear(base.Pix)
		fr := &frame{
			contents: contents,
			delay:    time.Millisecond * time.Duration(g.Delay[i]) * 10,
			lines:    strings.Count(contents, "\n"),
		}
		out.frames[i] = fr
	}

	return out, nil
}

func clonePix(b []uint8) []byte {
	c := make([]uint8, len(b))
	copy(c, b)
	return c
}

func (g *GIF) Play() {
	n := len(g.frames)
	if n == 0 {
		return
	}
	fmt.Print("\x1b[2J\x1b[H")

	i := 0
	for {
		f := g.frames[i]
		i++
		i %= n
		fmt.Print(f.contents)
		fmt.Printf("\x1b[%dF", f.lines)
		time.Sleep(f.delay)
	}
}

func (g *GIF) RenderFrame(n int) (string, error) {
	if n < 0 || n >= len(g.frames) {
		return "", errors.New("semigraph: frame out of bounds")
	}
	return g.frames[n].contents, nil
}
