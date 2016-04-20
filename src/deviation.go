package main

import (
	"image"
)

func meanContrast(img image.Image) float64 {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	var sum float64
	for ix := 0; ix < w; ix++ {
		for iy := 0; iy < h; iy++ {
			sum += contrast(img, ix, iy)
		}
	}
	return sum / float64(w*h)
}

func mean(floats []float64) float64 {
	var sum float64
	sum = 0
	for n := range floats {
		sum += floats[n]
	}
	return sum / float64(len(floats))
}

/*
func contrastZScore(img image.Image) float64 {
	mean := meanContrast(image.Image)

}
*/
