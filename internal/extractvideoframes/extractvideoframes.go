package extractvideoframes

import (
	"os/exec"
)

func ExtractVideoFrames() {
	cmd := exec.Command("mkdir", "temp/frames")
	cmd.Run()

	cmd = exec.Command("ffmpeg", "-i", "temp/video.mp4", "temp/frames/out-%03d.jpg")
	cmd.Run()
}

func RemoveFrames() {
	cmd := exec.Command("rm", "-rf", "temp/frames")
	cmd.Run()
}