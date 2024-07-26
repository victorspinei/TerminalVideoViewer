package main

import (
	"github.com/victor247k/TerminalVideoViewer/internal/render"
)
func main() {

	var horizontal_scale int = 16
	var vertical_scale int = 40

	render.RenderImageFromSrc("assets/fiona.jpg", horizontal_scale, vertical_scale);
}