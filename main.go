package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const inDir string = "./input/"
const outDir string = "./output/"

func main() {

	if _, err := os.Stat(inDir); os.IsNotExist(err) {
		err = os.Mkdir(inDir, 0775)
		if err != nil {
			panic("Can not create dir " + inDir)
		}
		fmt.Println("Input dir created")
		os.Exit(0)
	}

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.Mkdir(outDir, 0775)
		if err != nil {
			panic("Can not create dir " + outDir)
		}
	}

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	input := listDir()

	var wg sync.WaitGroup
	wg.Add(len(input))

	for _, f := range input {
		go goConvert(f, &wg)
	}

	wg.Wait()
}

func goConvert(f string, wg *sync.WaitGroup) {

	defer wg.Done()

	file, err := os.Open("./" + f)

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	pixels, size, err := getPixels(file)

	if err != nil {
		fmt.Print("Error: Image could not be decoded: ")
		fmt.Println(err)
		os.Exit(1)
	}

	edName := restoreImage(pixels, size, f)
	err = file.Close()

	if err != nil {
		panic("Can not close file")
	}

	fmt.Println("Image ready: ", edName)
}

func restoreImage(pixel [][]Pixel, size ImgSize, file string) string {
	upLeft := image.Point{}
	lowRight := image.Point{X: size.W, Y: size.H}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	for x := 0; x < size.W; x++ {
		for y := 0; y < size.H; y++ {
			px := toGray(pixel[y][x])
			clr := pixelToRGBA(px)
			img.Set(x, y, clr)
		}
	}

	name := filepath.Base(file)

	f, err := os.Create(outDir + name)
	if err != nil {
		panic("Can not restore image")
	}

	err = png.Encode(f, img)

	return f.Name()
}

func toGray(px Pixel) Pixel {
	var sum = int(px.R) + int(px.G) + int(px.B)
	var out = sum / 3
	if out > 255 { // -_-
		out = 255
	}
	return Pixel{uint8(out), uint8(out), uint8(out), px.A}
}

func pixelToRGBA(pixel Pixel) color.RGBA {
	return color.RGBA{R: pixel.R, G: pixel.G, B: pixel.B, A: pixel.A}
}

func getPixels(file io.Reader) ([][]Pixel, ImgSize, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, ImgSize{0, 0}, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var size = ImgSize{width, height}

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, size, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)}
}

func listDir() []string {
	files, err := filepath.Glob(inDir + "*.png")

	if err != nil {
		log.Fatal(err)
	}

	return files
}

type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type ImgSize struct {
	W int
	H int
}
