# GoWatermark
Go package "watermark" applies a specified logo in a configurable position and creates a new file from the original.

## Usage
#### See the 'example' directory
1. Run `git clone https://github.com/collinux/watermark`
2. Run `go run watermark.go <path_to_image>` with the path being jpg/jpeg file
3. A new image will be created from the original named "<image_name>_watermark.jpg" (does not override original)

#### Library
1. Run `go get github.com/collinux/watermark`
2. Add `import github.com/collinux/watermark`
2. Create a new Watermark struct with the configuration of your choice (only source is required by default)
3. Call `watermark.Apply(<filename>)` with the file that you want to place the watermark onto

##### Example
```
logo := watermark.Watermark{Source: "~/Pictures/my_logo.png"}
logo.Apply("~/Pictures/gopherart.jpg")
```

Remember to import `"github.com/collinux/watermark"`

## License
Copyright (C) 2016 Collinux
GPL version 2 or higher http://www.gnu.org/licenses/gpl.html  

## Contributing  
Pull requests happily accepted on GitHub!
