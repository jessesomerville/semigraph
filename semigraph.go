// The semigraph binary renders images in your terminal using semigraphics.
package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"strings"

	"image/gif"
	_ "image/jpeg"
	_ "image/png"

	semigraph "github.com/jessesomerville/semigraph/src"
)

func main() {
	flag.Parse()

	inPath := flag.Arg(0)
	if inPath == "" {
		fatalf("usage: semigraph <input_path>")
	}

	format, err := readImgType(inPath)
	if err != nil {
		fatalf("semigraph: %v", err)
	}

	f, err := os.Open(inPath)
	if err != nil {
		fatalf("semigraph: %v", err)
	}
	defer f.Close()
	switch format {
	case "gif":
		err = playGIF(f)
	case "png", "jpeg":
		err = showImage(f)
	}
	if err != nil {
		fatalf("semigraph: %v", err)
	}
}

func fatalf(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}

func readImgType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, format, err := image.DecodeConfig(f)
	return format, err
}

func showImage(f *os.File) error {
	input, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	fmt.Println(semigraph.Render(input))
	return nil
}

func playGIF(f *os.File) error {
	g, err := gif.DecodeAll(f)
	if err != nil {
		return err
	}
	gg, err := semigraph.RenderGIF(g)
	if err != nil {
		return err
	}
	gg.Play()
	return nil
}
