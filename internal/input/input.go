package input

import (
	"time"

	"github.com/eiannone/keyboard"

	"github.com/victor247k/TerminalVideoViewer/internal/audio"
	"github.com/victor247k/TerminalVideoViewer/internal/render"
)

const debounceDelay = 500 * time.Millisecond 
var lastKeyPressTime time.Time

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

		// Debounce logic
		if time.Since(lastKeyPressTime) < debounceDelay {
			continue
		}
		lastKeyPressTime = time.Now()

		if char == rune('q') || !audio.Running {
			audio.Close()
			render.Running = false
			return
		}
		switch {
		case key == keyboard.KeySpace:
			if audio.Paused {
				audio.Play()
			} else {
				audio.Pause()
			}
		case char == rune('z'):
			audio.SeekBackward()
		case char == rune('x'):
			audio.SeekForward()
		case char == rune('m'):
			audio.MuteVolume()
		}
	}
}