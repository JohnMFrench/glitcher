package main

import (
	"container/list"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"sort"
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
