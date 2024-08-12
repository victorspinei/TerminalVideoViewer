package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/victor247k/TerminalVideoViewer/internal/audio"
	"github.com/victor247k/TerminalVideoViewer/internal/download"
	"github.com/victor247k/TerminalVideoViewer/internal/extractvideoframes"
	"github.com/victor247k/TerminalVideoViewer/internal/input"
	"github.com/victor247k/TerminalVideoViewer/internal/menu"
	"github.com/victor247k/TerminalVideoViewer/internal/message"
	"github.com/victor247k/TerminalVideoViewer/internal/render"

	"github.com/victor247k/TerminalVideoViewer/internal/progressbar"

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

	menuProgram := tea.NewProgram(menu.InitialModel())	
	if _, err := menuProgram.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	switch menu.SelectedOption {
	case menu.NotSelected:
		os.Exit(0)
	case menu.YouTubeOption:
		var link string
	
		fmt.Printf("Paste your YouTube link\n  > ")
		fmt.Scan(&link)
	
		var downloadingWg sync.WaitGroup
		downloadingWg.Add(2)

		progressModel := progressbar.InitialModel("Downloading Video from YouTube")
		go func ()  {
  			defer downloadingWg.Done()
			p := tea.NewProgram(progressModel)
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running progress bar: %v\n", err)
				os.Exit(1)
			}
		}()
		go func ()  {
  			defer downloadingWg.Done()
			download.DownloadFromYoutubeLink(link)
		}()

  		downloadingWg.Wait()

	case menu.LocalOption:
		var path string
	
		fmt.Printf("Paste the path to your video file:\n  > ")
		fmt.Scan(&path)

		var downloadingWg sync.WaitGroup
		downloadingWg.Add(2)

		progressModel := progressbar.InitialModel("Preparing Local Video")
		go func ()  {
  			defer downloadingWg.Done()
			p := tea.NewProgram(progressModel)
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running progress bar: %v\n", err)
				os.Exit(1)
			}
		}()
		go func ()  {
  			defer downloadingWg.Done()
			download.CopyFromVideoPath(path, "temp/video.mp4")
		}()

  		downloadingWg.Wait()
	} 

	var exctractingWg sync.WaitGroup
	exctractingWg.Add(2)

	progressModel := progressbar.InitialModel("Extracting Video frames using FFmpeg")
	go func ()  {
		defer exctractingWg.Done()
		p := tea.NewProgram(progressModel)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running progress bar: %v\n", err)
			os.Exit(1)
		}
	}()
	go func ()  {
		defer exctractingWg.Done()
  		extractvideoframes.ExtractVideoFrames()
	}()

	exctractingWg.Wait()

	messageProgram := tea.NewProgram(message.InitialModel())	
	if _, err := messageProgram.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
  
  	var wg sync.WaitGroup
  	wg.Add(3)
  
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
  	go func() {
  		defer wg.Done()
		input.HandleInput()
  	}()
  
  	wg.Wait()

    clean()
}

func clean() {
	download.DeleteTempFiles()
	extractvideoframes.RemoveFrames()
	render.ClearTerminal()
}