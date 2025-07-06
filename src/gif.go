package semigraph

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"time"
)

type GIF struct {
	frames []*frame
	lines  int
}

type frame struct {
	delay    time.Duration
	contents string
}

func RenderGIF(g *gif.GIF) (*GIF, error) {
	if len(g.Image) == 0 {
		return nil, errors.New("semigraph: GIF has no frames")
	}
	if len(g.Image) != len(g.Delay) {
		return nil, errors.New("semigraph: mismatched GIF image and delay lengths")
	}
	if len(g.Image) != len(g.Disposal) {
		return nil, errors.New("semigraph: mismatched GIF image and disposal lengths")
	}

	srcWidth := g.Config.Width
	srcHeight := g.Config.Height

	out := &GIF{
		frames: make([]*frame, len(g.Image)),
		lines:  srcHeight,
	}
	gBounds := image.Rect(0, 0, srcWidth, srcHeight)
	norms := make([]image.Image, len(g.Image))
	for i, frm := range g.Image {
		base := image.NewRGBA(gBounds)
		if i > 0 {
			prev := norms[i-1]
			draw.Draw(base, gBounds, prev, image.Point{}, draw.Src)

			draw.Draw(base, gBounds, frm, image.Point{}, draw.Over)
		} else {
			draw.Draw(base, frm.Bounds(), frm, image.Point{}, draw.Src)
		}
		norms[i] = base
		contents := Render(base)
		fr := &frame{
			contents: contents,
			delay:    time.Millisecond * time.Duration(g.Delay[i]) * 10,
		}
		out.frames[i] = fr
	}

	return out, nil
}

func (g *GIF) Play() {
	n := len(g.frames)
	if n == 0 {
		return
	}
	fmt.Print("\x1b[2J\x1b[H")         // clear the screen and save the cursor position
	time.Sleep(time.Millisecond * 300) // wait before drawing

	i := 0
	for {
		f := g.frames[i]
		i++
		i %= n
		fmt.Print(f.contents)
		fmt.Printf("\x1b[%dF", g.lines) // return to the saved cursor position
		time.Sleep(f.delay)
	}
}

func (g *GIF) ShowFrame(n int) {
	fmt.Println(g.frames[n].contents)
}
