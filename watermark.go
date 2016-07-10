/*
* watermark.go
* Add a watermark to multiple images and customize the placement.
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package watermark

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

// Watermark quadrant position
const (
	CENTER       = 0
	TOP_LEFT     = 1
	TOP_RIGHT    = 2
	BOTTOM_LEFT  = 3
	BOTTOM_RIGHT = 4
)

type Watermark struct {
	// REQUIRED: File path of watermark png (MUST BE PNG)
	Source string

	// Placement of watermark image (default: bottom right)
	Position int

	// Margin around the watermark image (all default to 0")
	PaddingTop    int
	PaddingLeft   int
	PaddingRight  int
	PaddingBottom int
}

func (w *Watermark) getPosition(targetImage string) (int, int) {
	widthBounds, heightBounds := getImageDimensions(targetImage)
	placementWidth, placementHeight := 0, 0 // TOP_LEFT
	watermarkWidth, watermarkHeight := getImageDimensions(w.Source)
	if w.Position == CENTER {
		placementWidth = (widthBounds / 2) - (watermarkWidth / 2)
		placementHeight = (heightBounds / 2) - (watermarkHeight / 2)
	} else if w.Position == TOP_RIGHT {
		placementWidth = widthBounds - watermarkWidth
	} else if w.Position == BOTTOM_LEFT {
		placementHeight = heightBounds - watermarkHeight
	} else if w.Position == BOTTOM_RIGHT {
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
	log.Print("Adding watermark to image: " + file)

	// Open the watermark file
	watermarkBytes, err := os.Open(w.Source)
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
	offset := image.Pt(w.getPosition(file))
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

	log.Print("Successfully applied watermark.")
	return nil
}
