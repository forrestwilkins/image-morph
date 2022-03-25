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
	"github.com/disintegration/imaging"
)

const ImageCount = 5

func rgbaToGray(img image.Image) *image.NRGBA {
	gray := imaging.Grayscale(img)
	gray = imaging.AdjustContrast(gray, 20)
	gray = imaging.Sharpen(gray, 2)
	return gray
}

func withImaging(img image.Image, factor float64) *image.NRGBA {
	// dstImage := imaging.AdjustBrightness(img, factor*0.5)
	dstImage := imaging.Blur(img, factor*0.75)
	// dstImage = imaging.AdjustSaturation(img, factor*-8)
	// dstImage = imaging.AdjustGamma(img, factor*0.1)
	return dstImage
}

func makeEmbossed(img image.Image, factor float64) *image.NRGBA {
	adjustedFactor := factor*0.05 + 1
	embossed := imaging.Convolve3x3(
		img,
		[9]float64{
			-adjustedFactor, -adjustedFactor, 0,
			-adjustedFactor, adjustedFactor, adjustedFactor,
			0, adjustedFactor, adjustedFactor,
		},
		nil,
	)
	return embossed
}

func zeroToRandom(img image.Image) *image.NRGBA {
	var (
		bounds = img.Bounds()
		newImg = image.NewNRGBA(bounds)
	)

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			rgba := img.At(x, y)
			r, g, b, a := deMultiply(rgba)

			if r+g+b+a == 0 {
				randomR := uint8(rand.Intn(1))
				randomG := uint8(rand.Intn(255))
				randomB := uint8(rand.Intn(255))

				newImg.Set(x, y, color.RGBA{randomR, randomG, randomB, 255})
			} else {
				newImg.Set(x, y, color.RGBA{r, g, b, a})
			}
		}
	}

	return newImg
}

func deMultiply(preMultiplied color.Color) (uint8, uint8, uint8, uint8) {
	divideBy := uint32(257)
	r, g, b, a := preMultiplied.RGBA()
	deMultiplied := color.RGBA{uint8(r / divideBy), uint8(g / divideBy), uint8(b / divideBy), uint8(a / divideBy)}
	return deMultiplied.R, deMultiplied.G, deMultiplied.B, deMultiplied.A
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
	fmt.Println("Generating GIF...")

	outGif := &gif.GIF{}

	for imageNum := 1; imageNum < ImageCount; imageNum++ {
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
	for imageNum := 1; imageNum <= ImageCount; imageNum++ {
		fmt.Printf("Morphing image %d...\n\n", imageNum)

		img, _ := loadImage("assets/source/go_1.png")
		newImg := zeroToRandom(img)
		// newImg := withImaging(img, float64(imageNum))
		name := strings.Join([]string{"gen/new_", strconv.Itoa(imageNum), ".png"}, "")

		f, _ := os.Create(name)
		defer f.Close()
		png.Encode(f, newImg)
	}

	generateGif()
}
