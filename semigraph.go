// The semigraph binary renders images in your terminal using semigraphics.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	semigraph "github.com/jessesomerville/semigraph/src"
)

var (
	cpuprof = flag.String("cpuprof", "", "write a CPU profile to `file`")
	memprof = flag.String("memprof", "", "write a memory profile to `file`")
	noprint = flag.Bool("noprint", false, "render the input but don't output the results")
)

func main() {
	flag.Parse()

	if *cpuprof != "" {
		f, err := os.Create(*cpuprof)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

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
		if !*noprint {
			stop := gg.Play()
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			stop()
		}
	case "png", "jpeg":
		input, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			fatalf("semigraph: %v", err)
		}
		out := semigraph.Render(input)
		if !*noprint {
			fmt.Println(out)
		}
	}

	if *memprof != "" {
		f, err := os.Create(*memprof)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func fatalf(format string, args ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}
