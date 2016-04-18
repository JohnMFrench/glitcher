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
	"sort"
	"time"
)

func shift() {
	src := decode("smut.png")
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	out := image.NewRGBA(bounds)
	//top_color_list := getTopColors(src, 3)
	//randC1 := randColor()
	//randC2 := randColor()
	drip(src, out, 1, 0, 0.8)
	for i := 0; i < 50; i++ {
		drip(out, out, 1, 0, 0.8)
	}
	for i := 0; i < 25; i++ {
		drip(out, out, -1, 3, 0.95)
	}
	for ix := 0; ix < w; ix++ {
		for iy := 0; iy < h; iy++ {
			cont := contrast(src, ix, iy)
			if cont > 0.001 && cont < 0.015 {
				//out.Set(ix, iy, randC1)
			}
			if ix%50 == 0 && iy%50 == 0 {
				//fmt.Print("contrast at ", ix, ", ", iy, " = ", cont, "\n")
				//fmt.Print("color diff of ", colorDiff(src.At(ix, iy), src.At(ix+1, iy)), "\n")
			}
			if iy < 20 {
				//out.Set(ix, iy, randC1)
			}
			if iy < 40 && iy > 20 {
				//out.Set(ix, iy, randC2)
			}
		}
	}
	transposeX(src, out, 70, 320, 15)
	transposeY(src, out, 100, randY(src), randY(src))
	//encode the outfile and close
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

func randColor() color.RGBA {
	rand.Seed(time.Now().UTC().UnixNano())
	r := uint8(rand.Intn(255))
	g := uint8(rand.Intn(255))
	b := uint8(rand.Intn(255))
	return color.RGBA{r, g, b, 255}
}

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
					if x%100 == 0 && y%100 == 0 {
						//fmt.Print("dripping pixel at ", x, ", ", y, "\n")
					}
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

//where n is the number of cols to transpose and offset1 is the number of
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

func drawSolidRect(img image.Image, rect Rectangle) {

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
	fmt.Println(colors.Len(), " colors with a fuzz of ", fuzz)
	return colors
}

type Pair struct {
	Key   color.Color
	Value uint32
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func getTopColors(img image.Image, n int) *list.List {
	var color_totals map[color.RGBA]uint32
	color_totals = make(map[color.RGBA]uint32)
	colors := getColorsListFuzzy(img, 1.0)
	for c := colors.Front(); c != nil; c = c.Next() {
		color_totals[c.Value.(color.RGBA)] = 0
	}
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			cell_color := img.At(x, y)
			color_totals[cell_color.(color.RGBA)]++
		}
	}
	p1 := make(PairList, len(color_totals))
	i := 0
	for k, v := range color_totals {
		p1[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p1))
	for i2 := 0; i2 < n; i2++ {
		fmt.Print("Color #", i2+1)
		fmt.Print(p1[i2].Key)
		fmt.Print("\n")
	}
	/*
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
	*/
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
	//img := decode("jihadi-john.png")
	//getColorsListFuzzy(img, 0.5)
	//getTopColors(img, 5)
}

func main() {
	cleanEnv()
	shift()
	//test()
}
