package input

import (

	"github.com/eiannone/keyboard"
	"github.com/victor247k/TerminalVideoViewer/internal/audio"
	"github.com/victor247k/TerminalVideoViewer/internal/render"
)

func HandleInput() {		
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyEsc || char == rune('q') {
			audio.Close()
			render.Running = false
			break
		} else if key == keyboard.KeySpace {
			if audio.Paused {
				audio.Play()
			} else {
				audio.Pause()
			}
		} else if char == rune('m') {
			audio.MuteVolume()
		}
	}	
}