package main

import (
//	"strconv"
	"fmt"
	"os"
	"io/ioutil"
	"image"
	"errors"
	"regexp"
	"flag"
//	"path/filepath"
	"image/jpeg"
//	"image/png"

	"github.com/oliamb/cutter"
)

var path string = "."
var split bool = false

func main() {
	flag.StringVar(&path, "path", ".", "工作路径")
	flag.BoolVar(&split, "split", false, "是否切割图片")
	flag.Parse()

	fmt.Println("work path:", path)

	listFile(path)
}

func listFile(folder string) {
	files, _ := ioutil.ReadDir(folder)
	for _,file := range files{
		name := folder + "/" + file.Name()
		if file.IsDir(){
			listFile(name)
		} else {
			if checkImage(name) {
				splitImage(name)
			}
		}
	}
}

func checkImage(name string) bool {
	pic := regexp.MustCompile(`-poster\.jpg$`)
	if !pic.MatchString(name) {
		return false
	}

	return true
}

func loadImage(name string) (image.Image, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, errors.New("open file failed:" + name)
	}
	defer file.Close()

	return jpeg.Decode(file)
}

func splitImage(name string) error {
	img, err := loadImage(name)
	if err != nil {
		return errors.New("Load image failed:" + name)
	}

	size := img.Bounds().Max
	if size.X <= size.Y {
		return nil
	}

	fmt.Println("Found image:", name, size)
	if !split {
		return nil
	}

	fmt.Println("Start split:", name, "size:", size)

	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:  size.X / 2,
		Height: size.Y,
		Anchor: image.Point{size.X / 2, 0},
		Mode:   cutter.TopLeft, // optional, default value
	})
	if err != nil {
		return errors.New("Crop image failed:" + err.Error())
	}
	err = newImage(croppedImg, name)
	if err != nil {
		return err
	}


	return nil
}

func newImage(img image.Image, file string) error {
	os.Rename(file, file + ".bk.jpg")

	fo, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fo.Close()

	err = jpeg.Encode(fo, img, &jpeg.Options{100})
	if err != nil {
		return err
	}
	return nil
}


