package main

import (
	"image"
	"image/color"
	"math/rand"
	"time"
)

func randX(img image.Image) int {
	w := img.Bounds().Max.X
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(w)
}

func randY(img image.Image) int {
	h := img.Bounds().Max.Y
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(h)
}

func randColor() color.RGBA {
	rand.Seed(time.Now().UTC().UnixNano())
	r := uint8(rand.Intn(255))
	g := uint8(rand.Intn(255))
	b := uint8(rand.Intn(255))
	return color.RGBA{r, g, b, 255}
}
