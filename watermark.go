/*
* GoWatermark.go
* Add a watermark to multiple images and customize the placement.
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

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
		files = []string{args[0]}
		log.Print("Adding watermark to image: " + files[0])
	}

	watermark := Watermark{Position: TOP_RIGHT}
	watermark.Apply(files[0])
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

func getPosition(position int, widthBounds int, heightBounds int) (int, int) {
	width, height := 0, 0
	if position == CENTER {
		width = widthBounds / 2
		height = heightBounds / 2
	} else {
		width = 0
	}
	return width, height
}

func (w *Watermark) Apply(file string) error {
	originalBytes, err := os.Open(file)
	if err != nil {
		log.Fatal(err.Error() + ". Cannot find file: " + file)
	}
	originalImage, _, err := image.Decode(originalBytes)
	if err != nil {
		log.Fatal(err.Error() + ". Cannot decode image.")
	}

	filething, err := os.Open(file)
	if err != nil {
		log.Print(err.Error() + ". Cannot find image.")
	}
	width, height := getJpgDimensions(filething)
	log.Print(string(width), string(height))
	offset := image.Pt(getPosition(CENTER, width, height))
	bounds := originalImage.Bounds()
	defer originalBytes.Close()

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

// Credit: http://openmymind.net/Getting-An-Images-Type-And-Size/
func getDimensions(file *os.File) (width int, height int) {
	fi, _ := file.Stat()
	fileSize := fi.Size()

	position := int64(4)
	bytes := make([]byte, 4)
	file.ReadAt(bytes[:2], position)
	length := int(bytes[0]<<8) + int(bytes[1])
	for position < fileSize {
		position += int64(length)
		file.ReadAt(bytes, position)
		length = int(bytes[2])<<8 + int(bytes[3])
		if (bytes[1] == 0xC0 || bytes[1] == 0xC2) && bytes[0] == 0xFF && length > 7 {
			file.ReadAt(bytes, position+5)
			width = int(bytes[2])<<8 + int(bytes[3])
			height = int(bytes[0])<<8 + int(bytes[1])
			return
		}
		position += 2
	}
	return 0, 0
}
