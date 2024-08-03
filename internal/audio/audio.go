package audio

import (
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

var StartTime time.Time
var AudioDuration time.Duration
var PlaybackPosition time.Duration

func PlayAudio(src string) {
	file, err := os.Open(src)
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}

	// Decode file
	decodedMp3, err := mp3.NewDecoder(file)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.

	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

    StartTime = time.Now()

    fileInfo, err := file.Stat()
	if err != nil {
		panic("gettin gile stats failed: " + err.Error())
	}

    AudioDuration = time.Duration(fileInfo.Size()) * time.Second / time.Duration((op.SampleRate * 2 * 2)) // Estimate based on size and format

	// Create a new 'player' that will handle our sound. Paused by default.
	player := otoCtx.NewPlayer(decodedMp3)

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

    go func() {
		for player.IsPlaying() {
			PlaybackPosition = time.Since(StartTime)
			time.Sleep(100 * time.Millisecond) // Update playback position every 100ms
		}
	}()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// If you don't want the player/sound anymore simply close
	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}
