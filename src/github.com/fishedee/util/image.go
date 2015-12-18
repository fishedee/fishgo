package util

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type ImageSize struct {
	Width, Height int
}

func GetImageSize(data []byte) (ImageSize, error) {
	jpegReader := bytes.NewReader(data)
	config, _, err := image.DecodeConfig(jpegReader)
	if err != nil {
		return ImageSize{}, err
	}
	return ImageSize{
		Width:  config.Width,
		Height: config.Height,
	}, nil
}
