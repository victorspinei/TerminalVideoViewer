package extractvideoframes

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"regexp"
	"github.com/victor247k/TerminalVideoViewer/internal/progressbar"
)

var frameCount int

func ExtractVideoFrames() {
	progressbar.SetProgressMeter(0)
	os.MkdirAll("temp/frames", os.ModePerm)
	progressbar.SetProgressMeter(0.05)

	cmd := exec.Command("ffmpeg", "-i", "temp/video.mp4", "temp/frames/out-%03d.jpg")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Run()
	progressbar.SetProgressMeter(0.9)
	frameCount = countFrames(stderr.String())	
	progressbar.SetProgressMeter(1)
}

func countFrames(ffmpegOutput string) int {
	re := regexp.MustCompile(`frame=\s*(\d+)`)
	matches := re.FindAllStringSubmatch(ffmpegOutput, -1)
	if len(matches) == 0 {
		return 0
	}

	lastMatch := matches[len(matches)-1]
	frameNumber, err := strconv.Atoi(lastMatch[1])
	if err != nil {
		return 0
	}

	return frameNumber
}

func RemoveFrames() {
	cmd := exec.Command("rm", "-rf", "temp/frames",)
	cmd.Run()
}

func GetFrameCount() int {
	return frameCount
}