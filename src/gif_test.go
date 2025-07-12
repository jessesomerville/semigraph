package semigraph

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"testing"
)

// helper to create a mock GIF with n frames of size w x h
func createMockGIF(n, w, h int) *gif.GIF {
	g := &gif.GIF{
		Image:     make([]*image.Paletted, n),
		Delay:     make([]int, n),
		Disposal:  make([]byte, n),
		Config:    image.Config{Width: w, Height: h},
		LoopCount: 0,
	}
	p := color.Palette{color.Black, color.White}

	for i := 0; i < n; i++ {
		img := image.NewPaletted(image.Rect(0, 0, w, h), p)
		draw.Draw(img, img.Rect, &image.Uniform{color.White}, image.Point{}, draw.Src)
		g.Image[i] = img
		g.Delay[i] = 10
		g.Disposal[i] = gif.DisposalNone
	}
	return g
}

func BenchmarkRenderGIF_Small(b *testing.B) {
	g := createMockGIF(5, 32, 32)

	b.ResetTimer()
	for b.Loop() {
		_, err := RenderGIF(g)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderGIF_Medium(b *testing.B) {
	g := createMockGIF(10, 128, 128)

	b.ResetTimer()
	for b.Loop() {
		_, err := RenderGIF(g)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderGIF_Large(b *testing.B) {
	g := createMockGIF(20, 256, 256)

	b.ResetTimer()
	for b.Loop() {
		_, err := RenderGIF(g)
		if err != nil {
			b.Fatal(err)
		}
	}
}
