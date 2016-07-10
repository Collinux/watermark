package watermark

import (
	"github.com/collinux/gowatermark"
	"log"
	"testing"
)

func TestApply(t *testing.T) {
	watermark := watermark.Watermark{Source: "test_watermark.png"}
	err := watermark.Apply("test_image.jpeg")
	if err != nil {
		log.Fatal(err.Error() + ". Watermark was not applied.")
	}
	log.Println("Success! Open 'test_image_watermark.jpg' to see the result.")
}
