package render

import (
	"fmt"
	"image/jpeg"
	"os"
	"os/exec"
	"time"
	"sync"
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

func Render(frameCount int, frameDuration time.Duration, horizontal_scale int, vertical_scale int) {
	framesChan := make(chan [][][]int, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for frameNumber := 1; frameNumber <= frameCount; frameNumber++ {
			src := fmt.Sprintf("temp/frames/out-%03d.jpg", frameNumber)
			frame := preLoadFrame(src, horizontal_scale, vertical_scale)
			framesChan <- frame
		}
		close(framesChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for frame := range framesChan {
			//ClearTerminal()
			fmt.Print("\033[H")

			height := len(frame)
			width := len(frame[0])

			for j := 0; j < height; j++ {
				for i := 0; i < width; i++ {
					pixel := frame[j][i]
					fmt.Printf("\033[48;2;%d;%d;%dm ", pixel[0], pixel[1], pixel[2])
				}
				fmt.Printf("\033[m\n")
			}

			// Wait for the next frame
			time.Sleep(frameDuration)
		}
	}()

	wg.Wait() // Ensure both goroutines complete
}

func preLoadFrame(src string, scale int, vertical_scale int) [][][]int {
	// Open the image file
	imageFile, err := os.Open(src)
	if err != nil {
		fmt.Println("couldn't open image")
		return nil
	}
	defer imageFile.Close()

	// Decode the image
	loadedImage, err := jpeg.Decode(imageFile)
	if err != nil {
		fmt.Println("couldn't decode image")
		return nil
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

	return currentFrame
}