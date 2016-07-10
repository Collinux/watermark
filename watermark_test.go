/*
* watermark_test.go
* Demo that applies test_watermark.png on top of test_image.jpeg
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package watermark

import (
	"github.com/collinux/gowatermark"
	"log"
	"testing"
)

// Uses the included test_watermark.png and test_image.jpeg files
func TestApply(t *testing.T) {
	logo := watermark.Watermark{Source: "test_watermark.png"}
	err := logo.Apply("test_image.jpeg")
	if err != nil {
		log.Fatal(err.Error() + ". Watermark was not applied.")
	}
	log.Println("Success! Open 'test_image_watermark.jpg' to see the result.")
}
