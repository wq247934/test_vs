package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
)

func main() {
	_, err := GetCardImg("2_of_clubs.jpg", "2_of_hearts.jpg", "2_of_diamonds.jpg")
	if err != nil {
		log.Println(err)
		return
	}
}

const (
	dir    = "img/jpg"
	height = 244
	width  = 168
)

func GetCardImg(path ...string) (*os.File, error) {
	if len(path) == 1 {
		return os.Open(filepath.Join(dir, path[0]))
	}
	for i, p := range path {
		path[i] = filepath.Join(dir, p)
	}
	name, err := CreatePicture(path...)
	if err != nil {
		return nil, err
	}
	return os.Open(name)
}

func readImgData(filePath string) image.Image {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("图片decode失败", err)
		return nil
	}
	return img
}

// 计算图片缩放后的尺寸
func calculateRatioFit(srcWidth, srcHeight int, defaultWidth, defaultHeight float64) (int, int) {
	ratio := math.Min(defaultWidth/float64(srcWidth), defaultHeight/float64(srcHeight))
	return int(math.Ceil(float64(srcWidth) * ratio)), int(math.Ceil(float64(srcHeight) * ratio))
}

func CreatePicture(paths ...string) (string, error) {
	fmt.Println(paths)
	var images []image.Image
	//根据图片地址获取图片.
	for _, path := range paths {
		img := readImgData(path)
		if img == nil {
			continue
		}
		images = append(images, img)
	}
	for i, img := range images {
		//图片缩放
		b := img.Bounds()
		imgWidth := b.Max.X
		imgHeight := b.Max.Y
		w1, h1 := calculateRatioFit(imgWidth, imgHeight, width, height)
		images[i] = resize.Resize(uint(w1), uint(h1), img, resize.Lanczos3)
	}

	//创建源图
	fileName := filepath.Join("tmp", uuid.New().String()+"-dst.jpg")
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("CreateGoodsPicture:图片资源关闭错误", err)
		}
	}()
	//图片三合一绘图
	jpg := image.NewRGBA(image.Rect(0, 0, (width/2)*(len(images)+1), 244))
	for i, img := range images {
		draw.Draw(jpg, jpg.Bounds().Add(image.Pt(i*width/2, 0)), img, img.Bounds().Min, draw.Src)
	}
	//jpeg.Encode默认图片质量75%
	err1 := jpeg.Encode(file, jpg, nil)
	if err1 != nil {
		log.Println("CreateGoodsPicture:图片png.Encode错误", err1)
		return "", err
	}

	return file.Name(), nil
}
