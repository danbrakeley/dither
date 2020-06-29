package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()

	fnUsageAndQuit := func(err error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		fmt.Fprintf(os.Stderr, "usage: %s <width>x<height> <rgb1> <rgb2> <output.png>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "example:\n        %s 320x240 0000FF 000000 win31.png\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	if len(args) != 4 {
		fnUsageAndQuit(fmt.Errorf("expected 4 args, got %d", len(args)))
	}

	width, height, err := ParseWidthXHeight(args[0])
	if err != nil {
		fnUsageAndQuit(err)
	}

	c1, err := ParseRGBA(args[1])
	if err != nil {
		fnUsageAndQuit(err)
	}

	c2, err := ParseRGBA(args[2])
	if err != nil {
		fnUsageAndQuit(err)
	}

	// create target palette
	var pal color.Palette
	maxColors := 16
	for i := 0; i < maxColors; i++ {
		pal = append(pal, LerpRGB(c1, c2, uint8(i*255/(maxColors-1))))
	}

	// create outputimage
	outBounds := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{width, height},
	}

	outImg := image.NewRGBA(outBounds)

	// bayer4x4 := [][]uint8{
	// 	{0, 8, 2, 10},
	// 	{12, 4, 14, 6},
	// 	{3, 11, 1, 9},
	// 	{15, 7, 13, 5},
	// }

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			outImg.Set(x, y, pal.Convert(LerpRGB(c1, c2, uint8(y*255/(height-1)))))
		}
	}

	outFile, err := os.Create(args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating %s: %v\n", args[3], err)
		os.Exit(1)
	}
	err = png.Encode(outFile, outImg)
	outFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encoding %s: %v\n", args[3], err)
		os.Exit(1)
	}
}

func ParseWidthXHeight(arg string) (int, int, error) {
	res := strings.Split(strings.ToLower(arg), "x")
	if len(res) != 2 {
		return 0, 0, fmt.Errorf("unable to parse '%s' as '<width>x<height>'", arg)
	}

	width, err := strconv.Atoi(res[0])
	if err != nil {
		return 0, 0, err
	}

	height, err := strconv.Atoi(res[1])
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func ParseRGBA(arg string) (color.RGBA, error) {
	arg = strings.TrimPrefix(arg, "#")
	if len(arg) != 6 {
		return color.RGBA{}, fmt.Errorf("unable to parse '%s' as a hex color in the form RRGGBB", arg)
	}

	red, err := strconv.ParseUint(arg[:2], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	green, err := strconv.ParseUint(arg[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	blue, err := strconv.ParseUint(arg[4:], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: 0xff,
	}, nil
}

func LerpRGB(c1, c2 color.RGBA, t uint8) color.RGBA {
	it := int(t)
	return color.RGBA{
		R: uint8(((int(c2.R) - int(c1.R)) * it / 255) + int(c1.R)),
		G: uint8(((int(c2.G) - int(c1.G)) * it / 255) + int(c1.G)),
		B: uint8(((int(c2.B) - int(c1.B)) * it / 255) + int(c1.B)),
		A: 0xff,
	}
}
