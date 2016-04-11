package main

import (
	"container/list"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

func decode() {
	infile, err := os.Open("jihadi-john.png")
	if err != nil {
		fmt.Print(err)
	}
	defer infile.Close()
	src, _, err := image.Decode(infile)
	if err != nil {
		fmt.Print("error: ")
		fmt.Print(err)
	}
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	fmt.Print("width: ", w, "\n")
	fmt.Print("height: ", h, "\n")
	//white := src.At(0, 0)
	out := image.NewRGBA(bounds)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if contrast(src, x, y) < 0.1 && x < w-10 {
				out.Set(x+3, y, color.RGBA{255, 0, 0, 255})
			} else {
				out.Set(x, y, src.At(x, y))
			}
			//shift some pixels randomly
			/*
				rand.Seed(time.Now().UTC().UnixNano())
				random := rand.Intn(30)
				if random > 8 && x+random < w-1 {
					for i := 0; i < random; i++ {
						out.Set(x+i, y, src.At(x, y))
					}
				}
			*/

			//for debugging pixel values
			if x%30 == 0 && y%30 == 0 {
				/*
					fmt.Print("contrast=")
					fmt.Println(contrast(src, x, y))
				*/
			}
		}
	}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if contrast(src, x, y) > 0.8 {
				rand.Seed(time.Now().UTC().UnixNano())
				random := rand.Intn(10)
				if random < 7 {
					out.Set(x, y, color.Black)
				}
			}
			if brightness(src, x, y) < 225 {
				rand.Seed(time.Now().UTC().UnixNano())
				random := rand.Intn(10)
				if random < 3 && x > randX(src) && x < 155 {
					out.Set(x, y, colorCloserTo(src, x, y, color.RGBA{241, 53, 125, 255}))
				}
			}
			if brightness(out, x, y) == 0 && contrast(out, x, y) < 0.3 {
				rand.Seed(time.Now().UTC().UnixNano())
				r := rand.Intn(255)
				g := rand.Intn(255)
				b := rand.Intn(255)
				out.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
			}

		}
	}
	getColorList(src)
	//encode the outfile
	outfile, err := os.Create("out.png")
	if err != nil {
		fmt.Print(err)
	}
	defer outfile.Close()
	png.Encode(outfile, out)
}

func cleanEnv() {
	filename := "out.png"
	if _, err := os.Stat(filename); err == nil {
		err := os.Remove(filename)
		fmt.Print(err)
	}
}

func randX(img image.Image) int {
	bounds := img.Bounds()
	w := bounds.Max.X
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(w)
}

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
				around.PushBack(int(brightness(img, x+xplus, y+yplus)))
			}
		}
	}
	sum := 0
	for c := around.Front(); c != nil; c = c.Next() {
		sum += c.Value.(int)
	}
	avg := (float64(sum) / float64(around.Len()))
	return math.Abs(avg-float64(brightness(img, x, y))) / 255
}

//returns a list of all RGBA colors in the image
func getColorsList(img image.Image) List {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	colors := list.New()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			cellColor := img.At(x, y)
			found := false
			for color = colors.First; colors != nil; colors.Next() {
				if color == cellColor {
					found = true
					break
				}
			}
			if !found {
				colors.PushBack(cellColor)
			}
		}
	}
	fmt.Print(colors.Len(), " colors found")
	return colors
}

func colorCloserTo(img image.Image, x int, y int, newColor color.RGBA) color.RGBA {
	r, g, b, _ := img.At(x, y).RGBA()
	r2, g2, b2, _ := newColor.RGBA()
	nr := uint8(r + r2/2)
	ng := uint8(g + g2/2)
	nb := uint8(b + b2/2)
	return color.RGBA{nr, ng, nb, 255}
}

func main() {
	cleanEnv()
	decode()
}
