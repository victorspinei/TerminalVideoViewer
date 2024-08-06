package audio

import (
	"bytes"
	"io"

	//"io"
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
	audioData []byte
	decodedMp3 *mp3.Decoder
	op *oto.NewContextOptions
	endTime int64
	Paused bool
) 

const (
	seekInterval = 10 * 176400 // 10 seconds, 176400 = 44100 * 2 * 2
)

func PlayAudio(src string) {
	fileBytes, err := os.ReadFile(src)
	if err != nil {
	 panic("reading audio.mp3 failed: " + err.Error())
	}
 
	audioData = fileBytes
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
	endTime, _ = player.Seek(0, io.SeekEnd)
	player.Seek(0, io.SeekStart)
	player.Play()
	Paused = false

    go func() {
		for player.IsPlaying() {
			PlaybackPosition = time.Since(StartTime)
			PlaybackPosition += PlaybackOffset
			time.Sleep(100 * time.Millisecond)
		}
	}()

}

func SeekForward() {
	currentTime, _ := player.Seek(0, io.SeekCurrent)
	if currentTime + seekInterval >= endTime {
		player.Seek(0, io.SeekEnd)
		PlaybackOffset = AudioDuration - PlaybackPosition
	} else {
		player.Seek(seekInterval, io.SeekCurrent)
		PlaybackOffset += time.Second * 10
	}
}

func SeekBackward() {
	currentTime, _ := player.Seek(0, io.SeekCurrent)
	if currentTime - seekInterval < 0 {
		player.Seek(0, io.SeekStart)
		PlaybackOffset = -PlaybackPosition
	} else {
		player.Seek(-seekInterval, io.SeekCurrent)
		PlaybackOffset -= time.Second * 10
	}
}

func Pause() {
	player.Pause()
	Paused = true
}

func Play() {
	player.Play()
	Paused = false
}