package semigraph

import (
	"bytes"
	"image/png"
	"os"
	"testing"
)

func BenchmarkRender(b *testing.B) {
	data, err := os.ReadFile("testdata/benchRGB.png")
	if err != nil {
		b.Fatal(err)
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		b.Fatal(err)
	}
	input, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(cfg.Width * cfg.Height * 4))
	b.ReportAllocs()
	for b.Loop() {
		Render(input)
	}
}
