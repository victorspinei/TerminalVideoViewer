package render

import (
	"fmt"
	"image/jpeg"
	"os"
	"os/exec"
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

func ClearTerminal() {
	clear(FramesArray)
	cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

var FramesArray [][][][]int

func PreLoadFrame(src string, scale int, vertical_scale int) {
	// Open the image file
	imageFile, err := os.Open(src)
	if err != nil {
		fmt.Println("couldn't open image")
		return
	}
	defer imageFile.Close()

	// Decode the image
	loadedImage, err := jpeg.Decode(imageFile)
	if err != nil {
		fmt.Println("couldn't decode image")
		return
	}

	// Get image dimensions
	width := loadedImage.Bounds().Dx()
	height := loadedImage.Bounds().Dy()

	var currentFrame [][][]int

	// Loop over the image with the specified scales
	for j := 0; j < height; j += vertical_scale {
		var row [][]int
		for i := 0; i < width; i += scale {
			avgR, avgG, avgB, count := 0, 0, 0, 0

			// Calculate the average color for the scaled pixel
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

			// Average the color values
			if count > 0 {
				avgR /= count
				avgG /= count
				avgB /= count
			}

			// Create pixel and add to the row
			pixel := []int{avgR, avgG, avgB}
			row = append(row, pixel)
		}
		// Add row to current frame
		currentFrame = append(currentFrame, row)
	}

	// Append the current frame to the FramesArray
	FramesArray = append(FramesArray, currentFrame)
}