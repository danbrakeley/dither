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
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "usage: %s <width>x<height> <rgb1> <rgb2> <num-colors> <smooth> <output.png>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "       <width>: width of resulting image\n")
		fmt.Fprintf(os.Stderr, "      <height>: height of resulting image\n")
		fmt.Fprintf(os.Stderr, "        <rgb1>: top color of the gradient\n")
		fmt.Fprintf(os.Stderr, "        <rgb2>: bottom color of the gradient\n")
		fmt.Fprintf(os.Stderr, "  <num-colors>: number of colors in the palette to dither to\n")
		fmt.Fprintf(os.Stderr, "      <smooth>: smoothing type, one of 'none', 'both', 'out'\n")
		fmt.Fprintf(os.Stderr, "  <output.png>: file to write results to (contenst will be png, regardless of file extension)\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "example:\n")
		fmt.Fprintf(os.Stderr, "        %s 320x240 0000FF 000000 16 out winsetup.png\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	if len(args) != 6 {
		fnUsageAndQuit(fmt.Errorf("expected 6 args, got %d", len(args)))
	}

	width, height, err := ParseWidthXHeight(args[0])
	if err != nil {
		fnUsageAndQuit(fmt.Errorf("error parsing <width>x<height>: %v", err))
	}

	c1, err := ParseRGBA(args[1])
	if err != nil {
		fnUsageAndQuit(fmt.Errorf("error parsing <rgb1>: %v", err))
	}

	c2, err := ParseRGBA(args[2])
	if err != nil {
		fnUsageAndQuit(fmt.Errorf("error parsing <rgb2>: %v", err))
	}

	numColors, err := strconv.Atoi(args[3])
	if err != nil {
		fnUsageAndQuit(fmt.Errorf("error parsing <palette-size>: %v", err))
	}

	fnSmooth := NoSmoothing
	switch args[4] {
	case "none":
		break
	case "both":
		fnSmooth = SmoothStepByte
	case "out":
		fnSmooth = SmoothStepByteOut
	default:
		fnUsageAndQuit(fmt.Errorf("invalid value for <smooth>"))
	}

	// create palette
	var pal color.Palette = make([]color.Color, numColors)
	for i := 0; i < numColors; i++ {
		pal[i] = LerpRGB(c1, c2, uint8(i*255/(numColors-1)))
	}

	// create outputimage
	outBounds := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{width, height},
	}

	outImg := image.NewRGBA(outBounds)

	// bayer 4x4
	m := [][]int{
		{0, 8, 2, 10},
		{12, 4, 14, 6},
		{3, 11, 1, 9},
		{15, 7, 13, 5},
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := LerpRGB(c1, c2, fnSmooth(uint8(y*255/(height-1))))
			offset := m[x%4][y%4] - 7
			c.R = ClampUint8(int(c.R) + offset)
			c.G = ClampUint8(int(c.G) + offset)
			c.B = ClampUint8(int(c.B) + offset)
			outImg.Set(x, y, pal.Convert(c))
		}
	}

	outFile, err := os.Create(args[5])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating %s: %v\n", args[5], err)
		os.Exit(1)
	}
	err = png.Encode(outFile, outImg)
	outFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encoding %s: %v\n", args[5], err)
		os.Exit(1)
	}
}

func ParseWidthXHeight(arg string) (int, int, error) {
	res := strings.Split(strings.ToLower(arg), "x")
	if len(res) != 2 {
		return 0, 0, fmt.Errorf("does not contain exactly one 'x'")
	}

	width, err := strconv.Atoi(res[0])
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse width: %v", err)
	}

	height, err := strconv.Atoi(res[1])
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse height: %v", err)
	}

	return width, height, nil
}

func ParseRGBA(arg string) (color.RGBA, error) {
	arg = strings.TrimPrefix(arg, "#")
	if len(arg) != 6 {
		return color.RGBA{}, fmt.Errorf("more or less than 6 hex digits")
	}

	red, err := strconv.ParseUint(arg[:2], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("error parsing red: %v", err)
	}
	green, err := strconv.ParseUint(arg[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("error parsing green: %v", err)
	}
	blue, err := strconv.ParseUint(arg[4:], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("error parsing blue: %v", err)
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

func ClampUint8(n int) uint8 {
	if n > 255 {
		return 255
	} else if n < 0 {
		return 0
	}
	return uint8(n)
}

func NoSmoothing(x uint8) uint8 {
	return x
}

func SmoothStepByte(x uint8) uint8 {
	f := float64(x) / 255.0
	f = f * f * (3 - 2*f)
	return uint8(f * 255.0)
}

func SmoothStepByteOut(x uint8) uint8 {
	f := (float64(x) / 510.0)
	f = f * f * (3 - 2*f)
	return uint8(f * 510.0)
}
