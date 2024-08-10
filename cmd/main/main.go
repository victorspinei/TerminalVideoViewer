package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/victor247k/TerminalVideoViewer/internal/audio"
	"os"

	"github.com/victor247k/TerminalVideoViewer/internal/download"
	"github.com/victor247k/TerminalVideoViewer/internal/extractvideoframes"
	"github.com/victor247k/TerminalVideoViewer/internal/render"
	"github.com/victor247k/TerminalVideoViewer/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	horizontal_scale int = 4
	factor float64 = 2.5
	vertical_scale int = int(float64(horizontal_scale) * factor)
	numWorkers int = 8
)

func main() {
  	clean()

	p := tea.NewProgram(tui.InitialModel())	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	switch tui.SelectedOption {
	case tui.NotSelected:
		os.Exit(0)
	case tui.YouTubeOption:
		var link string
	
		fmt.Printf("Paste your YouTube link\n\t> ")
		fmt.Scan(&link)
	
		download.DownloadFromYoutubeLink(link)

	case tui.LocalOption:
		var path string
	
		fmt.Printf("Paste your Video path\n\t> ")
		fmt.Scan(&path)

		download.CopyFromVideoPath(path, "temp/video.mp4")
	} 

  	extractvideoframes.ExtractVideoFrames()
  
  	var wg sync.WaitGroup
  	wg.Add(2)
  
	render.ClearTerminal()
  	go func() {
  		defer wg.Done()
  		frameCount := extractvideoframes.GetFrameCount()
  		fps := download.GetFps()
  		frameDuration := time.Second / time.Duration(fps)
  
  		render.Render(frameCount, frameDuration, horizontal_scale, vertical_scale, numWorkers, fps)
  	}()
  	go func() {
  		defer wg.Done()
  		audio.PlayAudio("temp/audio.mp3")
  	}()
  
  	wg.Wait()
  
    clean()
}

func clean() {
	download.DeleteTempFiles()
	extractvideoframes.RemoveFrames()
	render.ClearTerminal()
}
