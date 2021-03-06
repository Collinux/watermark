/*
* watermark.go
* Add a watermark to multiple images and customize the placement.
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package watermark

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
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

type Watermark struct { // REQUIRED: File path of watermark png (MUST BE PNG)
	Source string

	// Placement of watermark image (default: bottom right)
	Position int

	// Margin around the watermark image (all default to 0)
	PaddingTop    int
	PaddingLeft   int
	PaddingRight  int
	PaddingBottom int
}

// Gets the position of where the watermark will be applied
// by looking at watermark.Position, padding, and logo size offset.
// Return the width, height
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

// Return the width and height of an image given the absolute path
func getImageDimensions(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}

// Apply a watermark on a file given the absolute path.
// Requires watermark.Source specified
func (w *Watermark) Apply(file string) error {
	// Open the watermark file
	watermarkBytes, err := os.Open(w.Source)
	if err != nil {
		return err
	}
	watermark, err := png.Decode(watermarkBytes)
	if err != nil {
		return err
	}
	defer watermarkBytes.Close()

	// Open the original image
	originalBytes, err := os.Open(file)
	if err != nil {
		return errors.New(fmt.Sprintf("%s. Cannot find file '%s'", err.Error(), file))
	}
	originalImage, _, err := image.Decode(originalBytes)
	if err != nil {
		return errors.New(err.Error() + ". Cannot decode image.")
	}
	offset := image.Pt(w.getPosition(file))
	originalImageBounds := originalImage.Bounds()
	defer originalBytes.Close()

	// Apply the watermark on top of the original
	colorRange := image.NewRGBA(originalImageBounds)
	draw.Draw(colorRange, originalImageBounds, originalImage, image.ZP, draw.Src)
	draw.Draw(colorRange, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	// Save the new image with "_watermark", not overriding the original file
	baseName := strings.Split(file, ".")[0]
	extension := ".jpg"
	newImage, err := os.Create(baseName + "_watermark" + extension)
	jpeg.Encode(newImage, colorRange, &jpeg.Options{jpeg.DefaultQuality})
	defer newImage.Close()

	return nil
}
