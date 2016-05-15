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

type perPixel func(img image.Image) float64

/*
func pixelFuncStDev(img image.Image, fun perPixel) float64 {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	mean := meanContrast(image.Image)
	var sum_squares float64
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			sum_squares += math.Pow(math.Abs(fun(img, x, y)-mean), 2)
		}
	}
	mean_sum_squares := sum_squares / (w * h)
	return math.Sqrt(mean_sum_squares)
}
*/
