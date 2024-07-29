package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/kbinani/screenshot"
)

func captureTrayIcon() (image.Image, error) {
	bounds := screenshot.GetDisplayBounds(0)
	width, height := 313, 37
	x := bounds.Dx() - width
	y := bounds.Dy() - height
	rect := image.Rect(x, y, x+width, y+height)

	img, err := screenshot.CaptureRect(rect)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func detectIconChange(prevImage, currentImage image.Image, threshold float64) bool {
	bounds := prevImage.Bounds()
	totalPixels := bounds.Dx() * bounds.Dy()
	diffPixels := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := prevImage.At(x, y).RGBA()
			r2, g2, b2, _ := currentImage.At(x, y).RGBA()

			if diffColor(r1, r2) || diffColor(g1, g2) || diffColor(b1, b2) {
				diffPixels++
			}
		}
	}

	return float64(diffPixels)/float64(totalPixels) > threshold
}

func diffColor(c1, c2 uint32) bool {
	const tolerance = 10 // 可以调整这个值来改变灵敏度
	return abs(int(c1)-int(c2)) > tolerance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func writeToFile(message string) {
	path := "./qywx.txt"
	tempPath := path + ".tmp"

	// 写入临时文件
	file, err := os.Create(tempPath)
	if err != nil {
		log.Println("Error creating temp file:", err)
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err = file.WriteString(fmt.Sprintf("%s: %s\n", timestamp, message))
	if err != nil {
		log.Println("Error writing to temp file:", err)
		file.Close()
		return
	}
	file.Close()

	// 重命名临时文件为目标文件
	err = os.Rename(tempPath, path)
	if err != nil {
		log.Println("Error renaming file:", err)
	}
}

func main() {
	prevImage, err := captureTrayIcon()
	if err != nil {
		log.Println("Error capturing initial image:", err)
		return
	}

	fmt.Println("Monitoring WeChat Work icon...")

	for {
		time.Sleep(1 * time.Second)
		currentImage, err := captureTrayIcon()
		if err != nil {
			log.Println("Error capturing image:", err)
			continue
		}

		if detectIconChange(prevImage, currentImage, 0.1) {
			writeToFile("WeChat Work icon flashing detected")
			fmt.Println("WeChat Work icon flashing detected")
		} else {
			fmt.Print(".")
		}

		prevImage = currentImage
	}
}
