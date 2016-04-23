package main

import (
	"image"
)

//where n is the number of cols totranspose and offset1 is the number of
//cols from the left edge where they should be lifted from and offset2 is
//the numb of cols from the left edge where they should be transposed to
func transposeX(img image.Image, out *image.RGBA, n int, offset1 int, offset2 int) {
	h := img.Bounds().Max.Y
	for ix := 0; ix < n; ix++ {
		for iy := 0; iy < h; iy++ {
			if isInBounds(img, offset1+ix, iy) && isInBounds(out, offset2+ix, iy) {
				out.Set(ix+offset2, iy, img.At(ix+offset1, iy))
			}
		}
	}
}

func transposeY(img image.Image, out *image.RGBA, n int, offset1 int, offset2 int) {
	w := img.Bounds().Max.X
	for iy := 0; iy < n; iy++ {
		for ix := 0; ix < w; ix++ {
			if isInBounds(img, ix, offset1+iy) && isInBounds(out, ix, offset2+iy) {
				out.Set(ix, iy+offset2, img.At(ix, iy+offset2))
			}
		}
	}
}
