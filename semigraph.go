// The semigraph binary renders images in your terminal using semigraphics.
package main

import (
	"flag"
	"fmt"
	"image"
	"math"
	"os"
	"strings"
	"time"

	"image/gif"
	_ "image/jpeg"
	_ "image/png"

	semigraph "github.com/jessesomerville/semigraph/src"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fatalf("usage: semigraph <input_path>")
	}
	input, format, err := readInput(flag.Arg(0))
	if err != nil {
		fatalf("%v", err)
	}
	if format == "gif" {
		// image.Decode only returns the first frame when reading a GIF
		// so it has to be read again to get the full contents.
		g, err := readGIF(flag.Arg(0))
		if err != nil {
			fatalf("%v", err)
		}

		gg, err := semigraph.RenderGIF(g)
		if err != nil {
			fatalf(err.Error())
		}
		// gg.ShowFrame(0)
		// gg.ShowFrame(1)
		gg.Play()

		// frames := make([]gifFrame, len(g.Image))
		// delays := make([]time.Duration, len(g.Delay))
		// for i, img := range g.Image {
		// 	frames[i] = newGIFFrame(img)
		// 	delays[i] = time.Millisecond * time.Duration(g.Delay[i]) * 10
		// }
		// loopGIF(frames, delays)

	} else {
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

func readInput(path string) (image.Image, string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer reader.Close()

	return image.Decode(reader)
}

func readGIF(path string) (*gif.GIF, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return gif.DecodeAll(reader)
}

type gifFrame struct {
	offsetX, offsetY int
	contents         []string
}

func newGIFFrame(img image.Image) gifFrame {
	b := img.Bounds()
	return gifFrame{
		offsetX:  int(math.Floor(float64(b.Min.X) / 2)),
		offsetY:  int(math.Floor(float64(b.Min.Y) / 4)),
		contents: strings.Split(semigraph.Render(img), "\n"),
	}
}

func (g gifFrame) Draw() {
	if g.offsetY > 0 {
		fmt.Printf("\x1b[%dE", g.offsetY)
	}
	for _, line := range g.contents {
		if g.offsetX > 0 {
			fmt.Printf("\x1b[%dC", g.offsetX)
		}
		fmt.Print(line)
		fmt.Print("\x1b[E")
	}
	fmt.Print("\x1b[H")
}

func loopGIF(frames []gifFrame, delays []time.Duration) {
	fmt.Print("\x1b[2J")
	i := 0
	n := len(frames)
	for {
		frame := frames[i]
		i++
		i %= n
		frame.Draw()
		time.Sleep(delays[i])
	}
}
