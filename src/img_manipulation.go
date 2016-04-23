package main

import (
	"fmt"
	"image"
	"image/color"
)

//where ltr means the image drips left to right and degree represents
//the likelihood that a pixel will drip into the adjacent
func drip(in image.Image, out *image.RGBA, ltr int, ttb int, degree float64) {
	rand.Seed(time.Now().UTC().UnixNano())
	w, h := in.Bounds().Max.X, in.Bounds().Max.Y
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if rand.Float64() < degree {
				//the pixel to drip from the in image
				if isInBounds(out, x+ltr, y+ttb) {
					out.Set(x, y, in.At(x+ltr, y+ttb))
				}
			} else {
				if isInBounds(out, x, y) {
					out.Set(x, y, in.At(x, y))
				}
			}
		}
	}
}

//block is the width of sections to drip
func blockDrip(in image.Image, out *image.RGBA, ltr int, ttb int,
	blockw int, blockh int, degree float64) {
	rand.Seed(time.Now().UTC().UnixNano())
	w, h := in.Bounds().Max.X, in.Bounds().Max.Y
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if rand.Float64() < degree {
				for bx := 0; bx < blockw; bx++ {
					for by := 0; by < blockh; by++ {
						if isInBounds(out, x+ltr+bx, y+ttb+by) {
							out.Set(x+ltr+bx, y+ttb+by, out.At(x, y))
						}
					}
				}
			}
		}
	}
}
