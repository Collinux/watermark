package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	// Parse arguments
	args := os.Args[1:]
	var files []string
	if len(args) == 0 {
		log.Print("Enter a picture source, '*' or '.' for all")
		os.Exit(1)
	}
	if strings.EqualFold(args[0], ".") || strings.EqualFold(args[0], "*") {
		// Get all files in current directory
		fileList, err := ioutil.ReadDir(".")
		for _, f := range fileList {
			// fmt.Println(f.Name())
			files = append(files, f.Name())
		}
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files = []string{args[0]}
		log.Print("Adding watermark to image " + files[0])
	}

	watermark := Watermark{Position: TOP_RIGHT}
	watermark.apply(files[0])
}

// Watermark quadrant position
const (
	CENTER       = 0
	TOP_LEFT     = 1
	TOP_RIGHT    = 2
	BOTTOM_LEFT  = 3
	BOTTOM_RIGHT = 4
)

type Watermark struct {
	// Margin around the watermark image
	PaddingTop    int
	PaddingLeft   int
	PaddingRight  int
	PaddingBottom int

	Position int    // Placement of watermark image
	Source   string // File path of image
}

func (w *Watermark) apply(file string) error {
	originalBytes, _ := os.Open(file)
	originalImage, _ := jpeg.Decode(originalBytes)
	originalImageConfig, _, err := image.DecodeConfig(originalBytes)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error getting original image's configuration.")
	}
	bounds := originalImage.Bounds()
	defer originalBytes.Close()

	watermarkBytes, err := os.Open("watermark.png")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error opening watermark image.")
	}
	watermark, err := png.Decode(watermarkBytes)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error decoding watermark.")
	}
	defer watermarkBytes.Close()

	offset := image.Pt(originalImageConfig.Width, originalImageConfig.Height)
	colorRange := image.NewRGBA(bounds)
	draw.Draw(colorRange, bounds, originalImage, image.ZP, draw.Src)
	draw.Draw(colorRange, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	baseName := strings.Split(file, ".")[0]
	extension := "." + strings.Split(file, ".")[1]
	newImage, err := os.Create(baseName + "_watermark" + extension)
	jpeg.Encode(newImage, colorRange, &jpeg.Options{jpeg.DefaultQuality})
	defer newImage.Close()
	log.Print("Finished")

	return nil
}
