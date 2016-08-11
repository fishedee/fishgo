package util

import (
	. "github.com/fishedee/assert"
	. "github.com/fishedee/crypto"
	// "reflect"
	"testing"
)

// func AssertEqual(t *testing.T, left interface{}, right interface{}) {
// 	if reflect.DeepEqual(left, right) == false {
// 		t.Errorf("assert fail: %+v != %+v", left, right)
// 	}
// }

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
		AssertEqual(t, err, nil)

		_, err = captcha.GetBase64Image(singleTestCase.width, singleTestCase.height)
		AssertEqual(t, err, nil)

		imageData, err := captcha.GetImage(singleTestCase.width, singleTestCase.height)
		AssertEqual(t, err, nil)

		image, err := NewImageFromString(imageData)
		AssertEqual(t, err, nil)

		imageSize, err := image.GetSize()
		AssertEqual(t, err, nil)
		AssertEqual(t, imageSize, ImageSize{
			Width:  singleTestCase.width,
			Height: singleTestCase.height,
		})
	}
}
