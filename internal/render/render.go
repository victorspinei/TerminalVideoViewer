package render

import (
	"os"
	"image/jpeg"
	"fmt"
)

func RenderImageFromSrc(src string, scale int, vertical_scale int) {
	imageFile, err := os.Open(src)
	if err != nil {
		fmt.Println("couldnt open image")
		return
	}

	defer imageFile.Close()

	imageFile.Seek(0, 0)

	loadedImage, err := jpeg.Decode(imageFile)
	if err != nil {
		fmt.Println("couldnt decode image")
		return
	}
	width := loadedImage.Bounds().Dx()
	height := loadedImage.Bounds().Dy()

	for j := 0; j < height; j += vertical_scale { 
		for i := 0; i < width; i += scale { 
			avgR := 0
			avgG := 0
			avgB := 0
			count := 0
			for r := 0; r < vertical_scale; r++ {
				for c := 0; c < scale; c++ {
					ir := i + c 
					jc := j + r
					if ir < width && jc < height {
						r, g, b, _ := loadedImage.At(ir, jc).RGBA()
						avgR += int(r >> 8) 
						avgG += int(g >> 8)
						avgB += int(b >> 8)
						count++
					}
				}
			}
			avgR /= count
			avgG /= count
			avgB /= count
			fmt.Printf("\033[48;2;%d;%d;%dm ", avgR, avgG, avgB)
		}
		fmt.Printf("\033[m\n")
	}
}