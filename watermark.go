/*
* GoWatermark.go
* Add a watermark to multiple images and customize the placement.
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package main

import (
	"fmt"
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
		log.Fatal("Enter a picture source, '*' or '.' for all")
	}
	if strings.EqualFold(args[0], ".") || strings.EqualFold(args[0], "*") {
		// Get all files in current directory
		fileList, err := ioutil.ReadDir(".")
		for _, f := range fileList {
			if strings.HasSuffix(strings.ToLower(f.Name()), "png") ||
				strings.Contains(strings.ToLower(f.Name()), "jpg") ||
				strings.Contains(strings.ToLower(f.Name()), "jpeg") {
				files = append(files, f.Name())
			}
		}
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Only ues one image
		// TODO: support multiple file arguments
		files = []string{args[0]}
		log.Print("Adding watermark to image: " + files[0])
	}

	// Get position argument (default bottom right)
	positionArg := args[1]
	position := BOTTOM_RIGHT
	log.Println("positionArg: ", positionArg)
	if strings.EqualFold(positionArg, "center") {
		position = CENTER
	} else if strings.EqualFold(positionArg, "top_left") {
		position = TOP_LEFT
	} else if strings.EqualFold(positionArg, "top_right") {
		position = TOP_RIGHT
	} else if strings.EqualFold(positionArg, "bottom_left") {
		position = BOTTOM_LEFT
	}
	watermark := Watermark{Position: position}
	for index, _ := range files {
		watermark.Apply(files[index])
	}
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

func getPosition(position int, widthBounds int, heightBounds int, watermark string) (int, int) {
	placementWidth, placementHeight := 0, 0 // TOP_LEFT
	watermarkWidth, watermarkHeight := getImageDimensions(watermark)
	if position == CENTER {
		placementWidth = (widthBounds / 2) - (watermarkWidth / 2)
		placementHeight = (heightBounds / 2) - (watermarkHeight / 2)
	} else if position == TOP_RIGHT {
		placementWidth = widthBounds - watermarkWidth
	} else if position == BOTTOM_LEFT {
		placementHeight = heightBounds - watermarkHeight
	} else if position == BOTTOM_RIGHT {
		placementHeight = heightBounds - watermarkHeight
		placementWidth = widthBounds - watermarkWidth
	}
	return placementWidth, placementHeight
}

func getImageDimensions(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}

func (w *Watermark) Apply(file string) error {
	// Open the watermark file
	watermarkBytes, err := os.Open("watermark.png")
	if err != nil {
		log.Fatal(err.Error() + ". Error opening watermark image.")
	}
	watermark, err := png.Decode(watermarkBytes)
	if err != nil {
		log.Fatal(err.Error() + ". Error decoding watermark.")
	}
	defer watermarkBytes.Close()
	// TODO: make sure watermark is not bigger than target file

	// Open the original image
	originalBytes, err := os.Open(file)
	if err != nil {
		log.Fatal(err.Error() + ". Cannot find file: " + file)
	}
	originalImage, _, err := image.Decode(originalBytes)
	if err != nil {
		log.Fatal(err.Error() + ". Cannot decode image.")
	}
	originalWidth, originalHeight := getImageDimensions(file)
	offset := image.Pt(getPosition(w.Position, originalWidth, originalHeight, "watermark.png"))
	bounds := originalImage.Bounds()
	defer originalBytes.Close()

	// Apply the watermark on top of the original
	colorRange := image.NewRGBA(bounds)
	draw.Draw(colorRange, bounds, originalImage, image.ZP, draw.Src)
	draw.Draw(colorRange, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	// Save the new image with "_watermark", not overriding the original file
	baseName := strings.Split(file, ".")[0]
	extension := ".jpg" //"." + strings.Split(file, ".")[1]
	newImage, err := os.Create(baseName + "_watermark" + extension)
	jpeg.Encode(newImage, colorRange, &jpeg.Options{jpeg.DefaultQuality})
	defer newImage.Close()

	log.Print("Finished")
	return nil
}
