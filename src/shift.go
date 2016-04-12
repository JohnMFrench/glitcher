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

func shift() {
	src := decode("jihadi-john.png")
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	fmt.Print("width: ", w, "\n")
	fmt.Print("height: ", h, "\n")
	//white := src.At(0, 0)
	out := image.NewRGBA(bounds)
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
	w := img.Bounds().Max.X
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(w)
}

func randY(img image.Image) int {
	h := img.Bounds().Max.Y
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(h)
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

func colorDiff(color1 color.Color, color2 color.Color) float64 {
	r1, g1, b1, _ := color1.RGBA()
	r2, g2, b2, _ := color2.RGBA()
	var diff float64
	diff = math.Abs(float64(r1-r2)) + math.Abs(float64(b1-b2)) + math.Abs(float64(g1-g2))
	return (diff / 3) / 255
}

//returns a list of colors after iterating through the image and collecting
//all ciolors that are unique enough to not contrast more than the fuzz allowed
func getColorsListFuzzy(img image.Image, fuzz float64) *list.List {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	colors := list.New()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			cellColor := img.At(x, y)
			found := false
			for c := colors.Front(); c != nil; c = c.Next() {
				if colorDiff(c.Value.(color.Color), cellColor) > fuzz {
					found = true
					break
				}
			}
			if !found {
				colors.PushBack(cellColor)
			}
		}
	}
	fmt.Print(colors.Len(), " colors with a fuzz of ", fuzz)
	return colors
}

func getTopColors(img image.Image) *list.List {
	var color_totals map[color.NRGBA]uint32
	color_totals = make(map[color.NRGBA]uint32)
	colors := getColorsListFuzzy(img, 1.0)
	for c := colors.Front(); c != nil; c = c.Next() {
		color_totals[c.Value.(color.NRGBA)] = 0
	}
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			cell_color := img.At(x, y)
			color_totals[cell_color.(color.NRGBA)]++
		}
	}
	var highest uint32
	highest = 0
	for c, _ := range color_totals {
		if color_totals[c] > highest {
			highest = color_totals[c]
			fmt.Print("highest=")
			fmt.Print(highest)
			fmt.Print(c)
			fmt.Print("\n")
		}
	}
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

func decode(filename string) image.Image {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Print(err)
	}
	defer file.Close()
	src, _, err := image.Decode(file)
	if err != nil {
		fmt.Print(err)
		return nil
	} else {
		return src
	}
}

func test() {
	img := decode("jihadi-john.png")
	//getColorsListFuzzy(img, 0.5)
	getTopColors(img)
}

func main() {
	/*
		cleanEnv()
		shift()
	*/
	test()
}
