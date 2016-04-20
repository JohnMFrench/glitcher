package main

import (
	"container/list"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

func shift() {
	src := decode("smut")
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	out := image.NewRGBA(bounds)
	//top_color_list := getTopColors(src, 3)
	randC1 := randColor()
	randC2 := randColor()
	drip(src, out, 1, 0, 0.8)
	for ix := 0; ix < w; ix++ {
		for iy := 0; iy < h; iy++ {
			cont := contrast(src, ix, iy)
			if ix%100 == 0 && iy%100 == 0 {
				fmt.Print("contrast at ", ix, ", ",
					iy, " = ", cont, "\n")
			}
			if cont < 0.01 {
				out.Set(ix, iy, colorCloserTo(src, ix, iy, randC1))
			}
			if cont > 0.01 && cont < 0.018 {
				out.Set(ix, iy, colorCloserTo(src, ix, iy, randC2))
			}
			if iy < 20 {
				out.Set(ix, iy, randC1)
			}
			if iy < 40 && iy > 20 {
				out.Set(ix, iy, randC2)
			}
		}
	}
	//blockDrip(src, out, 0, -1, 40, 30, 0.0005)
	//blockDrip(src, out, -1, -1, 20, 10, 0.0008)
	transposeX(src, out, 100, randX(src), randX(src))
	transposeY(src, out, 90, randY(src), randY(src))
	transposeY(src, out, 60, randY(src), randY(src))
	/*
		p1 := image.Point{5, 5}
		p2 := image.Point{200, 178}
		c := src.At(1, 1)
		drawSolidRect(out, image.Rectangle{p1, p2}, c)
	*/
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

//p is the top left-most pixel where the rectangle will be drawn
func drawSolidRect(img *image.RGBA, rect image.Rectangle, c color.Color) {
	c = c.(color.RGBA)
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := rect.Min.X; x < w; x++ {
		for y := rect.Min.Y; y < h; y++ {
			if isInBounds(img, x, y) {
				img.Set(x, y, c)
			}
		}
	}
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

func iterate(out *image.RGBA, foo func(c color.RGBA)) {
	w, h := out.Bounds().Max.X, out.Bounds().Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			foo(out.At(x, y).(color.RGBA))
		}
	}
}

//returns a list of colors after iterating through the image and collecting
//all ciolors that are unique enough to not contrast more than the fuzz allowed
func getColorsListFuzzy(img image.Image, fuzz float64) *list.List {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	colors := list.New()
	found := false
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			cellColor := img.At(x, y)
			found = false
			for c := colors.Front(); c != nil; c = c.Next() {
				if colorDiff(c.Value.(color.Color), cellColor) < fuzz {
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
	var color_totals map[color.Color]uint32
	color_totals = make(map[color.Color]uint32)
	colors := getColorsListFuzzy(img, 0.1)
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

//where filename is the name of the file without it's type
func decode(filename string) image.Image {
	//concatenate the filename (not performant)
	filename = "img/" + filename + ".png"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Print(err)
		panic(err)
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

func decodeGif(filename string) *gif.GIF {
	filename = "img/" + filename + ".gif"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	defer file.Close()
	gif, err := gif.DecodeAll(file)
	if err != nil {
		fmt.Print("error!!!\n:")
		fmt.Print(err)
		return gif
	} else {
		return gif
	}
}

func test() {
	//getColorsListFuzzy(img, 0.5)
	//getTopColors(img, 5)
	img := decodeGif("pool")
	fmt.Println(len(img.Image), " frames")
	shiftGif(img)
	//fmt.Print("found ", colors.Len(), " colors")
	//avg := meanContrast(img)
	//fmt.Print("average contrast=", avg)
}

func shiftGif(g *gif.GIF) {
	out, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
	w, h := g.Image[0].Bounds().Max.X, g.Image[0].Bounds().Max.Y
	var images []*image.Paletted
	var delay []int
	vaporwave_img := decode("vaporwave_palette2")
	num_of_colors := 29
	colors := getTopColors(vaporwave_img, num_of_colors)
	fmt.Print(colors.Len(), " colors found", "\n")
	var palette = make([]color.Color, num_of_colors)
	c := colors.Front()
	for i := 0; i < num_of_colors; i++ {
		palette[i] = c.Value.(color.Color)
		fmt.Println("added color ", palette[i])
		c = c.Next()
	}
	fmt.Print("images array size of ", len(images), "\n")
	fmt.Print("g.Image array size of ", len(g.Image), "\n")
	fmt.Print("w=", w, " & h=", h, "\n")
	for i2 := 0; i2 < len(g.Image); i2++ {
		frame_in := g.Image[i2]
		frame_out := image.NewPaletted(image.Rect(0, 0, w, h), palette)
		for ix := 0; ix < w; ix++ {
			for iy := 0; iy < h; iy++ {
				clor := frame_in.At(ix, iy)
				//r, g, b, _ := clor.RGBA()
				//fmt.Print("color vals of (", r, ", ", b, ", ", g, ")\n")
				//r2, g2, b2, _ := images[i2].At(ix, iy).RGBA()
				//fmt.Print("out color vals of (", r2, ", ", b2, ", ", g2, ")\n")
				//fmt.Print("at index ", ix, ", ", iy, "\n")
				//fmt.Println(clor)
				frame_out.Set(ix, iy, clor)
			}
		}
		fmt.Print(meanContrast(g.Image[i2]))
		images = append(images, frame_out)
		delay = append(delay, 0)
	}
	defer out.Close()
	gif.EncodeAll(out, &gif.GIF{
		Image: images,
		Delay: delay,
	})
}

func main() {
	//cleanEnv()
	//shift()
	test()
}
