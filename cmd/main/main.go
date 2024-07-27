package main

import (
	//"github.com/victor247k/TerminalVideoViewer/internal/render"
	"github.com/victor247k/TerminalVideoViewer/internal/download"
)

//var horizontal_scale int = 16
//var vertical_scale int = 40

func main() {
	//render.RenderImageFromSrc("assets/fiona.jpg", horizontal_scale, vertical_scale);

	//download.DownloadFromYoutubeLink("https://www.youtube.com/watch?v=MFDwFOygx_g")
	//download.DownloadFromYoutubeLink("https://www.youtube.com/watch?v=gZXqy3PbM3Q&t")
	download.DeleteTempFiles()
}