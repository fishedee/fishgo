package util

import (
	"bytes"
	"errors"

	"github.com/dchest/captcha"
	. "github.com/fishedee/encoding"
)

type Captcha struct {
	data []byte
}

func NewCaptchaFromDigit(dataStr string) (*Captcha, error) {
	if len(dataStr) == 0 {
		return nil, errors.New("invalid empty digitLength ")
	}
	var dataByte = []byte(dataStr)
	var result []byte
	for _, char := range dataByte {
		if char < '0' || char > '9' {
			return nil, errors.New("invalid dight char [" + dataStr + "]")
		}
		result = append(result, char-'0')
	}
	return &Captcha{
		data: result,
	}, nil
}

func (this *Captcha) GetBase64Image(width int, height int) (string, error) {
	data, err := this.GetImage(width, height)
	if err != nil {
		return "", err
	}

	base64Str, err := EncodeBase64(data)
	if err != nil {
		return "", err
	}
	base64Str = "data:image/png;base64," + base64Str
	return base64Str, nil
}

func (this *Captcha) GetImage(width int, height int) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	img := captcha.NewImage("", this.data, width, height)
	_, err := img.WriteTo(buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
