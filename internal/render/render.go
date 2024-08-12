package render

import (
	"fmt"
	"image/jpeg"
	"os"
	"os/exec"
	"time"
	"sync"
    "strings"

	"github.com/victor247k/TerminalVideoViewer/internal/audio"
)

func ClearTerminal() {
	clear(FramesArray)
	cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

var FramesArray []string
var Running bool

var barWidth = 0

func renderProgressBar(percent int) string {
	progress := (percent * barWidth) / 100
	return fmt.Sprintf("[%s%s] %d%%", 
		strings.Repeat("█", progress), 
		strings.Repeat(" ", barWidth-progress), 
		percent)
}

func Render(frameCount int, frameDuration time.Duration, horizontal_scale int, vertical_scale int, numWorkers int, fps float64) {
    Running = true
	framesChan := make(chan string, 1)
	var preloadWg sync.WaitGroup

	preloadWg.Add(1)
    previousSecond := 0
	go func() {
		defer preloadWg.Done()
        frameNumber := 1
		for frameNumber <= frameCount{
			src := fmt.Sprintf("temp/frames/out-%03d.jpg", frameNumber)
			frame := preLoadFrame(src, horizontal_scale, vertical_scale, numWorkers)
			framesChan <- frame

            // Update progress bar
            percent := int(frameNumber * 100 / frameCount)
			fmt.Printf("\r%s", renderProgressBar(percent))

            currentSecond := int(audio.PlaybackPosition / time.Second)
            if previousSecond != currentSecond {
                previousSecond = currentSecond
                frameNumber = int(fps * float64(currentSecond))
            }
            if !Running || !audio.Running {
                break
            }
            if !audio.Paused {
                frameNumber++ 
            }
		}
		close(framesChan)
	}()

	go func() {
		for frame := range framesChan {
            fmt.Print(frame)
            //time.Sleep(frameDuration)
        }
	}()

	preloadWg.Wait() 

    fmt.Print("\n\n")
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

    // Create and initialize the summed-area table
    summedAreaTable := make([][][3]int, height+1)
    for i := range summedAreaTable {
        summedAreaTable[i] = make([][3]int, width+1)
    }

    // Compute the summed-area table
    for y := 1; y <= height; y++ {
        for x := 1; x <= width; x++ {
            r, g, b, _ := loadedImage.At(x-1, y-1).RGBA()
            r, g, b = uint32(r>>8), uint32(g>>8), uint32(b>>8)
            summedAreaTable[y][x][0] = int(r) + summedAreaTable[y-1][x][0] + summedAreaTable[y][x-1][0] - summedAreaTable[y-1][x-1][0]
            summedAreaTable[y][x][1] = int(g) + summedAreaTable[y-1][x][1] + summedAreaTable[y][x-1][1] - summedAreaTable[y-1][x-1][1]
            summedAreaTable[y][x][2] = int(b) + summedAreaTable[y-1][x][2] + summedAreaTable[y][x-1][2] - summedAreaTable[y-1][x-1][2]
        }
    }

    // Function to get average color from the summed-area table
    getAverageColor := func(x1, y1, x2, y2 int) (int, int, int) {
        if x1 < 0 || y1 < 0 || x2 > width || y2 > height {
            return 0, 0, 0
        }
        totalArea := (y2 - y1) * (x2 - x1)
        if totalArea <= 0 {
            return 0, 0, 0
        }
        r := (summedAreaTable[y2][x2][0] - summedAreaTable[y1][x2][0] - summedAreaTable[y2][x1][0] + summedAreaTable[y1][x1][0]) / totalArea
        g := (summedAreaTable[y2][x2][1] - summedAreaTable[y1][x2][1] - summedAreaTable[y2][x1][1] + summedAreaTable[y1][x1][1]) / totalArea
        b := (summedAreaTable[y2][x2][2] - summedAreaTable[y1][x2][2] - summedAreaTable[y2][x1][2] + summedAreaTable[y1][x1][2]) / totalArea
        return r, g, b
    }

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
                    avgR, avgG, avgB := getAverageColor(i, j, i+scale, j+vertical_scale)
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
        if barWidth == 0 {
            barWidth = len(row) / 17
        }
    }

    return currentFrame
}