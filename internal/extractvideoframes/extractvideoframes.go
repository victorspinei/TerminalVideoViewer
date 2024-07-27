package extractvideoframes

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"regexp"
)

var frameCount int

func ExtractVideoFrames() {
	os.MkdirAll("temp/frames", os.ModePerm)

	cmd := exec.Command("ffmpeg", "-i", "temp/video.mp4", "temp/frames/out-%03d.jpg")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Run()
	frameCount = countFrames(stderr.String())	
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