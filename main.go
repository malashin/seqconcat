package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

var h, w, x int

func main() {
	if len(os.Args[1:]) < 1 {
		fmt.Println("usage: seqconcat [glob_pattern]")
		fmt.Println("seqconcat concatenates glob matched png images into one image horizontally left to right and saves it as png.")
		os.Exit(0)
	}

	input := os.Args[1]
	files, err := filepath.Glob(input)
	if err != nil {
		panic(err)
	}

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
