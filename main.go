package main

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"image/color"
	"image/gif"
	_ "image/jpeg"
	"image/png"

	"github.com/andybons/gogif"
)

func rgbaToGray(img image.Image) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			rgba := img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}
	return gray
}

func zeroToRandom(img image.Image) *image.RGBA {
	var (
		bounds = img.Bounds()
		newImg = image.NewRGBA(bounds)
	)

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			rgba := img.At(x, y)
			r, g, b, a := rgba.RGBA()
			alphaPremultipliedArray := []uint32{r, g, b, a}

			if r+g+b+a == 0 {
				randomInt1 := uint8(rand.Intn(1))
				randomInt2 := uint8(rand.Intn(255))
				randomInt3 := uint8(rand.Intn(255))

				newImg.Set(x, y, color.RGBA{randomInt1, randomInt2, randomInt3, 255})
			} else {
				fmt.Printf("rgba: %+v\n", rgba)
				fmt.Printf("alphaPremultipliedArray: %+v\n\n", alphaPremultipliedArray)

				newImg.Set(x, y, rgba)
			}
		}
	}

	return newImg
}

func loadImage(filepath string) (image.Image, error) {
	infile, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return img, nil
}

func generateGif() {
	outGif := &gif.GIF{}

	for imageNum := 1; imageNum < 10; imageNum++ {
		name := strings.Join([]string{"gen/new_", strconv.Itoa(imageNum), ".png"}, "")
		inPng, _ := loadImage(name)

		bounds := inPng.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: 64}
		quantizer.Quantize(palettedImage, bounds, inPng, image.ZP)

		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 0)
	}

	f, _ := os.OpenFile("gen/out.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}

func main() {
	img, _ := loadImage("go.png")
	newImg := zeroToRandom(img)

	f, _ := os.Create("gen/new_1.png")
	defer f.Close()
	png.Encode(f, newImg)

	generateGif()
}
