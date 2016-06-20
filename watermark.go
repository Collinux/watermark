package main

import (
	// "fmt"
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
	args := os.Args[1:]

	var files []string
	if len(args) == 0 {
		log.Print("Enter a picture source, '*' or '.' for all")
		os.Exit(1)
	}
	if strings.EqualFold(args[0], ".") || strings.EqualFold(args[0], "*") {
		// Get all files in current directory
		files, err := ioutil.ReadDir(".")
		for _, f := range files {
			files = append(files, f)
		}
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files = []string{args[0]}
	}
	watermark(files)
}

func watermark(files []string) {
	for _, file := range files {
		originalBytes, _ := os.Open(file)
		originalImage, _ := jpeg.Decode(originalBytes)
		defer originalBytes.Close()

		// _ = original
		watermarkBytes, _ := os.Open("watermark.jpeg")
		watermark, _ := png.Decode(watermarkBytes)
		defer watermarkBytes.Close()

		offsetWidth := 200
		offsetHeight := 200
		offset := image.Pt(offsetWidth, offsetHeight)
		bounds := originalImage.Bounds()
		colorRange := image.NewRGBA(bounds)
		draw.Draw(colorRange, bounds, originalImage, image.ZP, draw.Src)
		draw.Draw(colorRange, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

		newImageName := strings.SplitAfter(file, ".")
		_ = newImageName
		// newImage, _ := os.Create(file + "_watermark"
	}

}

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
