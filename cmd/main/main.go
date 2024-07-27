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

	//download.DownloadFromYoutubeLink("https://www.youtube.com/watch?v=MFDwFOygx_g")
	download.DownloadFromYoutubeLink("https://youtu.be/dQw4w9WgXcQ?si=2CF_hhZ8OLTDtOiA")
	extractvideoframes.ExtractVideoFrames()

	// Prepare to render frames
	frameCount := extractvideoframes.GetFrameCount()
	fps := download.GetFps()
	//fmt.Println(fps)
	frameDuration := time.Second / time.Duration(fps)

	render.ClearTerminal()
	for frameNumber := 1; frameNumber <= frameCount; frameNumber++ {
		// Move cursor to the top-left corner without clearing the screen
		fmt.Print("\033[H")

		src := fmt.Sprintf("temp/frames/out-%03d.jpg", frameNumber)
		render.RenderImageFromSrc(src, horizontal_scale, vertical_scale)

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