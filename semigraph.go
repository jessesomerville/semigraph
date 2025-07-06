// The semigraph binary renders images in your terminal using semigraphics.
package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"strings"

	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	semigraph "github.com/jessesomerville/semigraph/src"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fatalf("usage: semigraph <input_path>")
	}
	input, err := readInput(flag.Arg(0))
	if err != nil {
		fatalf("%v", err)
	}

	mid := image.NewRGBA(input.Bounds())
	draw.FloydSteinberg.Draw(mid, input.Bounds(), input, image.Point{})

	s := semigraph.Render(mid)
	fmt.Println(s)
}

func fatalf(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}

func readInput(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	return img, err
}
