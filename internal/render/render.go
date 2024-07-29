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
	var preloadWg, renderWg sync.WaitGroup
	var mu sync.Mutex

	preloadWg.Add(1)
	go func() {
		defer preloadWg.Done()
		for frameNumber := 1; frameNumber <= frameCount; frameNumber++ {
			src := fmt.Sprintf("temp/frames/out-%03d.jpg", frameNumber)
			frame := preLoadFrame(src, horizontal_scale, vertical_scale)
			framesChan <- frame
		}
		close(framesChan)
	}()

	renderWg.Add(1)
	go func() {
		defer renderWg.Done()
		for frame := range framesChan {
			mu.Lock()
			//ClearTerminal()
			fmt.Print("\033[H")

			height := len(frame)
			width := len(frame[0])
			groups := 2
			partHeight := height / groups // Split the frame into two parts

			var partRenderWg sync.WaitGroup
			partRenderWg.Add(groups)

			for part := 0; part < groups; part++ {
				go func(part int) {
					defer partRenderWg.Done()
					start := part * partHeight
					end := (part + 1) * partHeight
					if part == groups-1 {
						end = height
					}

					for j := start; j < end; j++ {
						for i := 0; i < width; i++ {
							pixel := frame[j][i]
							fmt.Printf("\033[%d;%dH\033[48;2;%d;%d;%dm ", j+1, i+1, pixel[0], pixel[1], pixel[2])
						}
						fmt.Print("\033[m\n")
					}
				}(part)
			}

			partRenderWg.Wait()
			mu.Unlock()
			time.Sleep(frameDuration)
		}
	}()

	preloadWg.Wait() // Ensure all frames are preloaded
	renderWg.Wait()  // Ensure rendering is complete
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