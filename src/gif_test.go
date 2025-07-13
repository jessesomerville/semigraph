package semigraph

import (
	"bytes"
	"image/gif"
	"os"
	"testing"
)

func BenchmarkRenderGIF(b *testing.B) {
	data, err := os.ReadFile("testdata/video-001.gif")
	if err != nil {
		b.Fatal(err)
	}
	cfg, err := gif.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		b.Fatal(err)
	}
	input, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(cfg.Width * cfg.Height))
	b.ReportAllocs()
	for b.Loop() {
		RenderGIF(input)
	}
}
