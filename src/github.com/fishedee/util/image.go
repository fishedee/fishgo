package util

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type ImageSize struct {
	Width, Height int
}

type Image struct {
	data []byte
}

func NewImageFromString(data []byte) (Image, error) {
	return Image{
		data: data,
	}, nil
}

func NewImageFromFile(fileName string) (Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return Image{}, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return Image{}, err
	}

	fileSize := int(fileInfo.Size())
	data := make([]byte, fileSize)
	size, err := file.Read(data)
	if err != nil {
		return Image{}, err
	}
	if size != fileSize {
		return Image{}, errors.New("文件未读写完全！")
	}

	return Image{
		data: data,
	}, nil
}

func (this *Image) GetSize() (ImageSize, error) {
	jpegReader := bytes.NewReader(this.data)
	config, _, err := image.DecodeConfig(jpegReader)
	if err != nil {
		return ImageSize{}, err
	}
	return ImageSize{
		Width:  config.Width,
		Height: config.Height,
	}, nil
}
