package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"os"
)

func cleanEnv() {
	filename := "out.png"
	if _, err := os.Stat(filename); err == nil {
		err := os.Remove(filename)
		fmt.Print(err)
	}
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

func main() {
	//cleanEnv()
	//shift()
	test()
}
