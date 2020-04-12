package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// Flags
var flagFrames int

var h, w, x int

func main() {
	// Parse input flags.
	flag.IntVar(&flagFrames, "f", 0, "instead of concatenating, crops input image into a specified number of separate frames; they are cut from a single row horizontal image.")
	flag.Usage = func() {
		fmt.Println("usage: seqconcat [glob_pattern]")
		fmt.Println("seqconcat concatenates glob matched png images into one image horizontally from left to right and saves it as png.")
		flag.PrintDefaults()
	}
	flag.Parse()

	input := flag.Args()

	if len(input) < 1 {
		flag.Usage()
		os.Exit(0)
	}

	for _, input := range input {
		files, err := filepath.Glob(input)
		if err != nil {
			panic(err)
		}

		switch {
		case flagFrames > 0:
			for _, fileName := range files {
				f, err := os.Open(fileName)
				if err != nil {
					panic(err)
				}

				src, err := png.Decode(f)
				f.Close()

				for i := 1; i <= flagFrames; i++ {
					w = src.Bounds().Max.X / flagFrames
					h = src.Bounds().Max.Y

					dst := image.NewRGBA(image.Rect(0, 0, w, h))

					fmt.Println(fileName, image.Rectangle{image.Point{x, 0}, image.Point{w, h}})
					draw.Draw(dst, image.Rectangle{image.Point{0, 0}, image.Point{w, h}}, src, image.Point{x, 0}, draw.Over)
					x += w

					outputName := fmt.Sprintf("%v_%05d.png", strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName)), i)

					out, err := os.Create(outputName)
					if err != nil {
						panic(err)
					}
					err = png.Encode(out, dst)
					if err != nil {
						panic(err)
					}
					err = out.Close()
					if err != nil {
						panic(err)
					}
					fmt.Printf("saved %q\n", outputName)
				}
			}
		default:
			for _, fileName := range files {
				f, err := os.Open(fileName)
				if err != nil {
					panic(err)
				}

				img, err := png.Decode(f)
				f.Close()

				w += img.Bounds().Max.X
				if h < img.Bounds().Max.Y {
					h = img.Bounds().Max.Y
				}
			}

			dst := image.NewRGBA(image.Rect(0, 0, w, h))

			for _, fileName := range files {
				f, err := os.Open(fileName)
				if err != nil {
					panic(err)
				}

				src, err := png.Decode(f)
				f.Close()

				fmt.Println(fileName, image.Rectangle{image.Point{x, 0}, image.Point{src.Bounds().Max.X + x, src.Bounds().Max.Y}})
				draw.Draw(dst, image.Rectangle{image.Point{x, 0}, image.Point{src.Bounds().Max.X + x, src.Bounds().Max.Y}}, src, image.ZP, draw.Over)
				x += src.Bounds().Max.X
			}

			out, err := os.Create("output.png")
			if err != nil {
				panic(err)
			}
			err = png.Encode(out, dst)
			if err != nil {
				panic(err)
			}
			err = out.Close()
			if err != nil {
				panic(err)
			}
			fmt.Println("saved \"output.png\"")
		}
	}
}
