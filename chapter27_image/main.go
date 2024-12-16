package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	// 创建一个新的图像，尺寸为 200x200
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	// 绘制一个红色的正方形
	red := color.RGBA{R: 255, A: 255}
	for x := 50; x < 150; x++ {
		for y := 50; y < 150; y++ {
			img.Set(x, y, red)
		}
	}

	// 创建一个输出文件
	file, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 将图像写入输出文件
	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
}
