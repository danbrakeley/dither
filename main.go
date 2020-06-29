package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <input.png|gif|jpg> <output.png>\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	inFile, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening %s: %v\n", args[0], err)
		os.Exit(1)
	}

	inImg, _, err := image.Decode(inFile)
	inFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decoding %s: %v\n", args[0], err)
		os.Exit(1)
	}

	// choose palette
	// var pal color.Palette = palette.WebSafe
	var pal color.Palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0x80, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0x80, 0x00, 0xff},
		color.RGBA{0x80, 0x80, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0x80, 0xff},
		color.RGBA{0x80, 0x00, 0x80, 0xff},
		color.RGBA{0x00, 0x80, 0x80, 0xff},
		color.RGBA{0xC0, 0xC0, 0xC0, 0xff},
		color.RGBA{0x80, 0x80, 0x80, 0xff},
		color.RGBA{0xFF, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0xFF, 0x00, 0xff},
		color.RGBA{0xFF, 0xFF, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0xFF, 0xff},
		color.RGBA{0xFF, 0x00, 0xFF, 0xff},
		color.RGBA{0x00, 0xFF, 0xFF, 0xff},
		color.RGBA{0xFF, 0xFF, 0xFF, 0xff},
	}

	// create outputimage
	inBounds := inImg.Bounds()
	if inBounds.Min.X != 0 || inBounds.Min.Y != 0 {
		fmt.Fprintf(os.Stderr, "assumption that images start at 0,0 broken\n")
		os.Exit(1)
	}
	outImg := image.NewRGBA(inBounds)

	width := inBounds.Dx()
	height := inBounds.Dy()

	// bayer4x4 := [][]uint8{
	// 	{0, 8, 2, 10},
	// 	{12, 4, 14, 6},
	// 	{3, 11, 1, 9},
	// 	{15, 7, 13, 5},
	// }

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			outImg.Set(x, y, pal.Convert(inImg.At(x, y)))
		}
	}

	outFile, err := os.Create(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating %s: %v\n", args[1], err)
		os.Exit(1)
	}
	err = png.Encode(outFile, outImg)
	outFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encoding %s: %v\n", args[1], err)
		os.Exit(1)
	}
}
