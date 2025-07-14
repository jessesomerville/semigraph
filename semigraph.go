// The semigraph binary renders images in your terminal using semigraphics.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"os"
	"strings"

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

	data, err := os.ReadFile(inPath)
	if err != nil {
		fatalf("semigraph: %v", err)
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		fatalf("semigraph: %v", err)
	}

	switch format {
	case "gif":
		g, err := gif.DecodeAll(bytes.NewReader(data))
		if err != nil {
			fatalf("semigraph: %v", err)
		}
		gg, err := semigraph.RenderGIF(g)
		if err != nil {
			fatalf("semigraph: %v", err)
		}
		// TODO: Gracefully handle cleanup.
		gg.Play()
	case "png", "jpeg":
		input, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			fatalf("semigraph: %v", err)
		}
		fmt.Println(semigraph.Render(input))
	}
}

func fatalf(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}
