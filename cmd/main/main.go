package main

import (
	//	"fmt"
	//	"sync"
	//	"time"
	//
	//	"github.com/victor247k/TerminalVideoViewer/internal/audio"
	"fmt"
	"os"

	"github.com/victor247k/TerminalVideoViewer/internal/download"
	"github.com/victor247k/TerminalVideoViewer/internal/extractvideoframes"
	"github.com/victor247k/TerminalVideoViewer/internal/render"
	"github.com/victor247k/TerminalVideoViewer/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	//horizontal_scale int = 4
	//factor float64 = 2.5
	//vertical_scale int = int(float64(horizontal_scale) * factor)
	//numWorkers int = 8
)

func main() {
//  clean()
	p := tea.NewProgram(tui.InitialModel())	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
	}
//	var link string
//
//	fmt.Println("Paste your YouTube link:")
//	fmt.Scan(&link)
//
//	download.DownloadFromYoutubeLink(link)
//
//	extractvideoframes.ExtractVideoFrames()
//
//	var wg sync.WaitGroup
//	wg.Add(2)
//
//	go func() {
//		defer wg.Done()
//		frameCount := extractvideoframes.GetFrameCount()
//		fps := download.GetFps()
//		frameDuration := time.Second / time.Duration(fps)
//
//		render.Render(frameCount, frameDuration, horizontal_scale, vertical_scale, numWorkers, fps)
//	}()
//	go func() {
//		defer wg.Done()
//		audio.PlayAudio("temp/audio.mp3")
//	}()
//
//	wg.Wait()
//
	clean()
}

func clean() {
	download.DeleteTempFiles()
	extractvideoframes.RemoveFrames()
	render.ClearTerminal()
}
