package audio

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

var (
    StartTime time.Time
    AudioDuration time.Duration
    PlaybackPosition time.Duration
	PlaybackOffset time.Duration
    otoCtx *oto.Context
	readyChan chan struct{}
	player *oto.Player
	decodedMp3 *mp3.Decoder
	op *oto.NewContextOptions
	EndTime int64
	Paused bool
	Running bool = true
	volume float64
) 

const (
	seconds = 5
	seekInterval = seconds * 176400 // seconds, 176400 = 44100 * 2 * 2
)

func PlayAudio(src string) {
	fileBytes, err := os.ReadFile(src)
	if err != nil {
	 panic("reading audio.mp3 failed: " + err.Error())
	}
 
	fileBytesReader := bytes.NewReader(fileBytes)

	decodedMp3, err = mp3.NewDecoder(fileBytesReader)
	if err != nil {
	 panic("mp3.NewDecoder failed: " + err.Error())
	}

	op = &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err = oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	audioDuration := decodedMp3.Length() 
	if audioDuration <= 0 {
		panic("invalid audio length")
	}
	AudioDuration = time.Duration(audioDuration) * time.Second / time.Duration(op.SampleRate*op.ChannelCount*2)

    StartTime = time.Now()
    PlaybackOffset = 0

	player = otoCtx.NewPlayer(decodedMp3)
	EndTime, _ = player.Seek(-1, io.SeekEnd)
	EndTime++
	player.Seek(0, io.SeekStart)
	player.Play()
	Paused = false

    go func() {
		for player.IsPlaying() {
			PlaybackPosition = time.Since(StartTime) + PlaybackOffset
			time.Sleep(100 * time.Millisecond)
		}
	}()

}

func SeekForward() {
	Pause()
	currentTime, _ := player.Seek(0, io.SeekCurrent)
	if currentTime + seekInterval >= EndTime {
		Close()
	} else {
		player.Seek(seekInterval, io.SeekCurrent)
		PlaybackOffset += time.Second * seconds
	}
	Play()
}

func SeekBackward() {
	Pause()
	currentTime, _ := player.Seek(0, io.SeekCurrent)
	if currentTime - seekInterval < 0 {
		player.Seek(0, io.SeekStart)
		PlaybackOffset = -PlaybackPosition
	} else {
		player.Seek(-seekInterval, io.SeekCurrent)
		PlaybackOffset -= time.Second * seconds
	}
	Play()
}

func Pause() {
	if !Paused {
		player.Pause()
		Paused = true
		// Update PlaybackOffset to reflect the pause time
		PlaybackOffset += time.Since(StartTime) - PlaybackPosition
	}
}

func Play() {
	if Paused {
		player.Play()
		Paused = false
		// Adjust StartTime to account for the pause duration
		StartTime = time.Now().Add(-PlaybackPosition)
	}
}

func MuteVolume() {
	if player.Volume() == 0 {
		player.SetVolume(volume)
	} else {
		volume = player.Volume()
		player.SetVolume(0)
	}
}

func Close() {
	player.Close()
	Running = false
}