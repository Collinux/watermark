# GoWatermark
Go package "Watermark" applies a specified logo in a configurable position and creates a new file from the original.

# Usage
### Standalone Script
1. Run `git clone https://github.com/Collinux/gowatermark`
2. Run `go get ./...` to get all dependencies
3. Open `watermark.yaml` and customize it to your own configuration (see inline comments)
4. Run `go run watermark.go <path_to_image>` with the path being jpg/jpeg file
5. A new image will be created from the original named "<image_name>_watermark.jpg" (does not override original)


### Library
1. Run `go get github.com/Collinux/gowatermark`
2. Add `import github.com/Collinux/gowatermark`
2. Create a new Watermark struct with the configuration of your choice (only source is required by default)
3. Call `watermark.Apply(<filename>)` with the file that you want to place the watermark onto

#### Example
```
watermark := Watermark{Source: "~/Pictures/my_logo.png"}
watermark.Apply("~/Pictures/photo_album/gopherart.jpg")
```
