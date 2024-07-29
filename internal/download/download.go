package download

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"

	"github.com/kkdai/youtube/v2"
)

func findFormatByItag(formats youtube.FormatList, itag int) *youtube.Format {
	for _, format := range formats {
		if format.ItagNo == itag {
			return &format
		}
	}
	return nil
}

func DownloadFromYoutubeLink(link string) {
	DeleteTempFiles()
	client := youtube.Client{}

	video, err := client.GetVideo(link)
	if err != nil {
        log.Fatalf("Error getting video: %v", err)
    }

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		downloadVideoFile(client, *video)
	}()
	go func() {
		defer wg.Done()
		downloadAudioFile(client, *video)
	}()

	wg.Wait()

	fmt.Println("Downloaded 720p video and audio streams")

	return 
	// Merge video and audio
	err = mergeVideoAndAudio("temp/video.mp4", "temp/audio.m4a", "temp/output.mp4")
	if err != nil {
		log.Fatalf("Error merging video and audio: %v", err)
	}

	fmt.Println("Merged video and audio into output.mp4")
}

func mergeVideoAndAudio(videoFilePath, audioFilePath, outputFilePath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoFilePath,
		"-i", audioFilePath,
		"-c:v", "copy",
		"-c:a", "aac",
		"-strict", "experimental",
		outputFilePath,
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func DeleteTempFiles() {
    cmd := exec.Command("rm", "temp/audio.m4a", "temp/output.mp4", "temp/video.mp4",)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

    cmd.Run()
}

func GetFps() float64 {
	cmd := exec.Command("ffmpeg", "-i", "temp/video.mp4")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Run()

	// Print the ffmpeg output for debugging
	output := stderr.String()
	fmt.Println("FFmpeg Output:", output)

	// Extract the FPS from ffmpeg output
	re := regexp.MustCompile(`(\d+(?:\.\d+)?) fps`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		log.Fatalf("FPS not found in ffmpeg output")
	}

	fps, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		log.Fatalf("Error parsing FPS: %v", err)
	}

	return fps
}

func downloadVideoFile(client youtube.Client, video youtube.Video) {
	videoFormat := findFormatByItag(video.Formats, 136) // 136 is for 720p video
	if videoFormat == nil {
		log.Fatal("720p video format not found")
	}

	videoStream, _, err := client.GetStream(&video, videoFormat)
	if err != nil {
		log.Fatalf("Error getting video stream: %v", err)
	}

	videoFile, err := os.Create("temp/video.mp4")
	if err != nil {
		log.Fatalf("Error creating video file: %v", err)
	}
	defer videoFile.Close()

	_, err = videoFile.ReadFrom(videoStream)
	if err != nil {
		log.Fatalf("Error downloading video: %v", err)
	}
}

func downloadAudioFile(client youtube.Client, video youtube.Video) {
	audioFormat := findFormatByItag(video.Formats, 140) // 140 is for audio
	if audioFormat == nil {
		log.Fatal("Audio format not found")
	}


	audioStream, _, err := client.GetStream(&video, audioFormat)
	if err != nil {
		log.Fatalf("Error getting audio stream: %v", err)
	}

	audioFile, err := os.Create("temp/audio.m4a")
	if err != nil {
		log.Fatalf("Error creating audio file: %v", err)
	}
	defer audioFile.Close()

	_, err = audioFile.ReadFrom(audioStream)
	if err != nil {
		log.Fatalf("Error downloading audio: %v", err)
	}
	convertAudioFile()
}

func convertAudioFile() {
	cmd := exec.Command("ffmpeg", "-i", "temp/audio.m4a", "temp/audio.mp3",)
	cmd.Run()
}