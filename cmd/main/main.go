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

	link := "https://youtu.be/Qi_-rfVdx5w?si=kO28ij8TpuMJB5TA"

	download.DownloadFromYoutubeLink(link)

	extractvideoframes.ExtractVideoFrames()

	// Prepare to render frames
	frameCount := extractvideoframes.GetFrameCount()
	fps := download.GetFps()
	frameDuration := time.Second / time.Duration(fps)

	fmt.Printf("Begining to load frames and render")
	render.Render(frameCount, frameDuration, horizontal_scale, vertical_scale)

	clean()
}

func clean() {
	download.DeleteTempFiles()
	extractvideoframes.RemoveFrames()
	render.ClearTerminal()
}