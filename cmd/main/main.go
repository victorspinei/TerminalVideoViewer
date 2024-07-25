package main

import (
	//"fmt"
	//"image/jpeg"
	//"os"
)

func main() {
//	imageFile, err := os.Open("assets/fiona.jpg")
//	if err != nil {
//		fmt.Println("couldnt open image")
//		return
//	}
//
//	defer imageFile.Close()
//
//	imageFile.Seek(0, 0)
//
//	loadedImage, err := jpeg.Decode(imageFile)
//	if err != nil {
//		fmt.Println("couldnt decode image")
//		return
//	}
//	width := loadedImage.Bounds().Dx()
//	height := loadedImage.Bounds().Dy()
//	scale := 8
//	verticalScale := scale * 2
//
//	fmt.Println(width / scale, height / verticalScale)
//
//	for j := 0; j < height; j += verticalScale { // Iterate through height first
//		for i := 0; i < width; i += scale { // Iterate through width
//			avgR := 0
//			avgG := 0
//			avgB := 0
//			count := 0
//			for r := 0; r < verticalScale; r++ {
//				for c := 0; c < scale; c++ {
//					ir := i + c // Corrected iteration for width
//					jc := j + r // Corrected iteration for height
//					if ir < width && jc < height {
//						r, g, b, _ := loadedImage.At(ir, jc).RGBA()
//						avgR += int(r >> 8) // Scale down from 0-65535 to 0-255
//						avgG += int(g >> 8)
//						avgB += int(b >> 8)
//						count++
//					}
//				}
//			}
//			avgR /= count
//			avgG /= count
//			avgB /= count
//			fmt.Printf("\033[48;2;%d;%d;%dm ", avgR, avgG, avgB)
//		}
//		fmt.Printf("\033[m\n")
//	}


}