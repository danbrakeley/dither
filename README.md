# dither

## Overview

Generate a dithered gradient.

Trying to recreate the classic "Windows 3.1 Setup" look.

## Usage

```text
usage: dither <width>x<height> <rgb1> <rgb2> <num-colors> <smooth> <output.png>

       <width>: width of resulting image
      <height>: height of resulting image
        <rgb1>: top color of the gradient
        <rgb2>: bottom color of the gradient
  <num-colors>: number of colors in the palette to dither to
      <smooth>: smoothing type, one of 'none', 'both', 'out'
  <output.png>: file to write results to (contenst will be png, regardless of file extension)
```

## Example

`dither 320x240 0000FF 000000 16 out example.png`

![example dithered gradient](example.png)
