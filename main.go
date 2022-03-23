package main

import (
	"fmt"
	"image"
	"math/rand"
	"os"

	"image/color"
	_ "image/jpeg"
	"image/png"
)

func rgbaToGray(img image.Image) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}
	return gray
}

func rgbaToRandom(img image.Image) *image.RGBA {
	var (
		bounds = img.Bounds()
		newImg = image.NewRGBA(bounds)
	)

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y)

			r, g, b, a := rgba.RGBA()

			var alphaPremultipliedArray = []uint32{r, g, b, a}

			if r > 0 {
				fmt.Printf("rgba: %+v\n", rgba)
				fmt.Printf("alphaPremultipliedArray: %+v\n\n", alphaPremultipliedArray)

				newImg.Set(x, y, rgba)
			} else {
				var randomInt1 = uint8(rand.Intn(255))
				var randomInt2 = uint8(rand.Intn(255))
				var randomInt3 = uint8(rand.Intn(255))

				newImg.Set(x, y, color.RGBA{randomInt1, randomInt2, randomInt3, 255})
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

func main() {
	var img, _ = loadImage("go.png")
	var newImg = rgbaToRandom(img)

	f, _ := os.Create("new.png")
	defer f.Close()
	png.Encode(f, newImg)
}
