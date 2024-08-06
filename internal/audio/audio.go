package audio

import (
	"bytes"
	"fmt"
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
) 

const (
	seekInterval = 10 * 176400 // 10 seconds, 176400 = 44100 * 2 * 2
)

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

    go func() {
		for player.IsPlaying() {
			PlaybackPosition = time.Since(StartTime)
			PlaybackPosition += PlaybackOffset
			time.Sleep(100 * time.Millisecond)
		}
	}()

}

func SeekAudio(targetTime time.Duration) {
	if player != nil {
		player.Close()
	}

	// **Ensure the target time is within the bounds**
	if targetTime < 0 || targetTime > AudioDuration {
		fmt.Printf("Error: seek time %v is out of bounds. AudioDuration is %v\n", targetTime, AudioDuration)
		panic("seek time is out of bounds")
	}

	// **Create a new decoder from the buffered audio data**
	audioDataReader := bytes.NewReader(audioData)
	decodedMp3, _ = mp3.NewDecoder(audioDataReader)

	seekOffset := int64(time.Duration(targetTime)) / int64(time.Second) * int64(op.SampleRate*op.ChannelCount*2) / 10
	fmt.Printf("SeekOffset calculated as %d bytes\n", seekOffset)

	n, err := decodedMp3.Seek(int64(seekOffset), 0)
  	if err != nil {
    	fmt.Printf("Error seeking in audio data: %v\n", err)
    	panic("failed to seek in audio data")
  	}
  	fmt.Printf("Seeked to byte position %d\n", n)

	//_ = player.Close()
	//player = otoCtx.NewPlayer(decodedMp3)
	//player.Play()
	player.Seek(-2000, 0)

	StartTime = time.Now()
	PlaybackPosition = targetTime

	fmt.Println("check 1")
	PlaybackPosition = time.Since(StartTime) + targetTime
	fmt.Println("check 2")
	time.Sleep(100 * time.Millisecond) 
}