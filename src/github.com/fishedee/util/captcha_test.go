package util

import (
	"reflect"
	"testing"

	. "github.com/fishedee/crypto"
)

func assertCaptchaEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}

func TestCaptchaOutOfBound(t *testing.T) {
	for i := 0; i != 1000; i++ {
		word := CryptoRandDigit(4)
		width := 100
		height := 50
		captcha, err := NewCaptchaFromDigit(word)
		if err != nil {
			panic(err)
		}
		_, err = captcha.GetBase64Image(width, height)
		if err != nil {
			panic(err)
		}
	}
}

func TestCaptcha(t *testing.T) {
	testCase := []struct {
		origin string
		width  int
		height int
	}{
		{"123", 100, 100},
		{"1227", 100, 50},
		{"67", 200, 50},
	}

	for _, singleTestCase := range testCase {
		captcha, err := NewCaptchaFromDigit(singleTestCase.origin)
		assertCaptchaEqual(t, err, nil)

		_, err = captcha.GetBase64Image(singleTestCase.width, singleTestCase.height)
		assertCaptchaEqual(t, err, nil)

		imageData, err := captcha.GetImage(singleTestCase.width, singleTestCase.height)
		assertCaptchaEqual(t, err, nil)

		image, err := NewImageFromString(imageData)
		assertCaptchaEqual(t, err, nil)

		imageSize, err := image.GetSize()
		assertCaptchaEqual(t, err, nil)
		assertCaptchaEqual(t, imageSize, ImageSize{
			Width:  singleTestCase.width,
			Height: singleTestCase.height,
		})
	}
}
