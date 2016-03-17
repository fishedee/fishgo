package crypto

import (
	"testing"
)

func TestCryptoSha1(t *testing.T) {
	testCase := []struct {
		origin string
		target string
	}{
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"123", "40bd001563085fc35165329ea1ff5c5ecbdbbeef"},
		{"你好", "440ee0853ad1e99f962b63e459ef992d7c211722"},
	}
	for _, singleTestCase := range testCase {
		result := CryptoSha1([]byte(singleTestCase.origin))
		if result != singleTestCase.target {
			t.Errorf("%v != %v,[%v]", result, singleTestCase.target, singleTestCase.origin)
			return
		}
	}
}
