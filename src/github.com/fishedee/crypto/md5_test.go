package crypto

import (
	"testing"
)

func TestCryptoMd5(t *testing.T) {
	testCase := []struct {
		origin string
		target string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"123", "202cb962ac59075b964b07152d234b70"},
		{"你好", "7eca689f0d3389d9dea66ae112e5cfd7"},
	}
	for _, singleTestCase := range testCase {
		result := CryptoMd5([]byte(singleTestCase.origin))
		if result != singleTestCase.target {
			t.Errorf("%v != %v,[%v]", result, singleTestCase.target, singleTestCase.origin)
			return
		}
	}
}
