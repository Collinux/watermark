/*
* watermark_test.go
* Applies the "test_watermark.png" to a specified jpg source
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*
* Usage: go run example_program.go test_image.jpg
* 		 go run example_program.go .
 */

package main

import (
	"github.com/collinux/watermark"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Get input of files to watermark
	targetList := os.Args[1:]
	if len(targetList) == 0 {
		log.Fatal("An image file is required as an argument.")
	}

	// Prepare the logo
	logo := watermark.Watermark{
		Source:   "test_watermark.png",
		Position: watermark.BOTTOM_RIGHT,
	}

	// Look for a wildcard and add image files to the list
	if targetList[0] == "*" || targetList[0] == "." {
		targetList = targetList[:0]
		files, _ := ioutil.ReadDir("./")
		for _, f := range files {
			fileName := f.Name()
			targetList = append(targetList, fileName)
		}
	}

	// Only accept jpg or jpeg, must not have "_watemark" already added
	for _, image := range targetList {
		imageRegex := regexp.MustCompile(`.(jpg|jpeg)`)
		if !imageRegex.MatchString(image) || strings.Contains(image, "_watermark") {
			continue
		}
		err := logo.Apply(image)
		if err != nil {
			log.Fatal("Failed to apply watermark to image: ", image)
		} else {
			log.Print("Successfully applied watermark to image: ", image)
		}
	}
}
