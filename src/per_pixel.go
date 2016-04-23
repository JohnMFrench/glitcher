package main

/*
This file contains functions that are applied to specific indexes
of an img, e.g. individual pixels or colors
*/
import (
	"container/list"
	"image"
	"image/color"
	"math"
)

func brightness(img image.Image, x int, y int) uint8 {
	c := img.At(x, y)
	r, g, b, _ := c.RGBA()
	return uint8((r + g + b) / 3)
}

func isInBounds(img image.Image, x, y int) bool {
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	return (x > 0 && x < w && y > 0 && y < h)
}

func contrast(img image.Image, x, y int) float64 {
	around := list.New()
	for xplus := -1; xplus < 2; xplus++ {
		for yplus := -1; yplus < 2; yplus++ {
			if isInBounds(img, x+xplus, y+yplus) {
				around.PushBack(colorDiff(img.At(x, y), img.At(x+xplus, y+yplus)))
			}
		}
	}
	var sum float64
	sum = 0
	for c := around.Front(); c != nil; c = c.Next() {
		sum += c.Value.(float64)
	}
	var avg float64
	avg = sum / float64(around.Len())
	return avg
}

func colorDiff(color1 color.Color, color2 color.Color) float64 {
	r1, g1, b1, _ := color1.RGBA()
	r2, g2, b2, _ := color2.RGBA()
	var nr1, nr2, ng1, ng2, nb1, nb2, diff float64
	nr1, ng1, nb1 = float64(r1), float64(g1), float64(b1)
	nr2, ng2, nb2 = float64(r2), float64(g2), float64(b2)
	diff = (math.Abs(nr1-nr2) / 255) + (math.Abs(ng1-ng2) / 255) + (math.Abs(nb1-nb2) / 255)
	return diff / 3 / 255
}

func invertColor(c color.RGBA) color.RGBA {
	r, g, b, _ := c.RGBA()
	inverted := color.RGBA{uint8(math.Abs(float64(r - 255))),
		uint8(math.Abs(float64(g - 255))),
		uint8(math.Abs(float64(b - 255))),
		255}
	return inverted
}

func iterate(out *image.RGBA, foo func(c color.RGBA)) {
	w, h := out.Bounds().Max.X, out.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			foo(out.At(x, y).(color.RGBA))
		}
	}
}

func colorCloserTo(img image.Image, x int, y int, newColor color.RGBA) color.RGBA {
	r, g, b, _ := img.At(x, y).RGBA()
	r2, g2, b2, _ := newColor.RGBA()
	nr := uint8(r + r2/2)
	ng := uint8(g + g2/2)
	nb := uint8(b + b2/2)
	return color.RGBA{nr, ng, nb, 255}
}
