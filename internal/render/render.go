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

var FramesArray []string

func Render(frameCount int, frameDuration time.Duration, horizontal_scale int, vertical_scale int, numWorkers int) {
	framesChan := make(chan string, 1)
	var preloadWg sync.WaitGroup

	preloadWg.Add(1)
	go func() {
		defer preloadWg.Done()
		for frameNumber := 1; frameNumber <= frameCount; frameNumber++ {
			src := fmt.Sprintf("temp/frames/out-%03d.jpg", frameNumber)
			frame := preLoadFrame(src, horizontal_scale, vertical_scale, numWorkers)
			framesChan <- frame
		}
		close(framesChan)
	}()

	go func() {
		for frame := range framesChan {
            fmt.Print(frame)
            time.Sleep(frameDuration)
        }
	}()

	preloadWg.Wait() 
}
func preLoadFrame(src string, scale int, vertical_scale int, numWorkers int) string {
    // Open the image file
    imageFile, err := os.Open(src)
    if err != nil {
        fmt.Println("couldn't open image")
        return ""
    }
    defer imageFile.Close()

    // Decode the image
    loadedImage, err := jpeg.Decode(imageFile)
    if err != nil {
        fmt.Println("couldn't decode image")
        return ""
    }

    // Get image dimensions
    width := loadedImage.Bounds().Dx()
    height := loadedImage.Bounds().Dy()

    var currentFrame string = "\033[H"
    rows := make([]string, height/vertical_scale)
    var wg sync.WaitGroup
    var mu sync.Mutex

    // Calculate the height per worker
    segmentHeight := height / numWorkers

    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for j := workerID * segmentHeight; j < (workerID+1)*segmentHeight; j += vertical_scale {
                if j >= height {
                    break
                }
                var row string
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
                    row += fmt.Sprintf("\033[48;2;%d;%d;%dm ", avgR, avgG, avgB)
                }
                // Store row in the correct position
                mu.Lock()
                rows[j/vertical_scale] = row
                mu.Unlock()
            }
        }(w)
    }

    wg.Wait()

    // Combine rows into the final frame
    for _, row := range rows {
        currentFrame += row
        currentFrame += "\n"
    }

    return currentFrame
}
