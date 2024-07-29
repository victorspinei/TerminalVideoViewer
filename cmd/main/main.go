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
const factor float64 = 2.5
var vertical_scale int = int(float64(horizontal_scale) * factor)

func main() {
	clean()

	link := "https://youtu.be/Nuanwn3v-2I?si=3nKI9Td8SrHjJTYD"

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