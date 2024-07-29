package main

import (
	//"github.com/victor247k/TerminalVideoViewer/internal/render"

	"fmt"
	"time"
	"github.com/victor247k/TerminalVideoViewer/internal/download"
	"github.com/victor247k/TerminalVideoViewer/internal/extractvideoframes"
	"github.com/victor247k/TerminalVideoViewer/internal/render"
)

var horizontal_scale int = 4
var vertical_scale int = 10

func main() {
	clean()

	link := "https://youtu.be/Nuanwn3v-2I?si=LlpSm_Qk7U_bDc3p"

	download.DownloadFromYoutubeLink(link)

	extractvideoframes.ExtractVideoFrames()

	// Prepare to render frames
	frameCount := extractvideoframes.GetFrameCount()
	fps := download.GetFps()
	//fmt.Println(fps)
	frameDuration := time.Second / time.Duration(fps)



	fmt.Println("preloading frames")
	for frameNumber := 1; frameNumber <= frameCount; frameNumber++ {
		src := fmt.Sprintf("temp/frames/out-%03d.jpg", frameNumber)
		render.PreLoadFrame(src, horizontal_scale, vertical_scale)
	}
	fmt.Println("finished preloading frames")
	fmt.Println("started rendereing frames")
	for frameNumber := 1; frameNumber <= frameCount; frameNumber++ {
		// Move cursor to the top-left corner without clearing the screen
		fmt.Print("\033[H")

		height := len(render.FramesArray[frameNumber-1])
		width := len(render.FramesArray[frameNumber-1][0])

		for j := 0; j < height; j++ { 
			for i := 0; i < width; i++ { 
				pixel := render.FramesArray[frameNumber-1][j][i]
				fmt.Printf("\033[48;2;%d;%d;%dm ", pixel[0], pixel[1], pixel[2])
			}
			fmt.Printf("\033[m\n")
		}

		// Wait for the next frame
		time.Sleep(frameDuration)
	}

	clean()
}

func clean() {
	download.DeleteTempFiles()
	extractvideoframes.RemoveFrames()
	render.ClearTerminal()
}