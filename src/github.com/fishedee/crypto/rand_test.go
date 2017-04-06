package crypto

import (
	"testing"
)

func TestRand(t *testing.T) {
	testCase := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		100, 500, 1000,
	}

	for _, singleTestCase := range testCase {
		result := CryptoRand(singleTestCase)
		if len(result) != singleTestCase {
			t.Errorf("assert false! %v != %v", len(result), singleTestCase)
		}
		for _, singleResult := range result {
			if singleResult >= '0' && singleResult <= '9' {
				continue
			}
			if singleResult >= 'a' && singleResult <= 'z' {
				continue
			}
			if singleResult >= 'A' && singleResult <= 'Z' {
				continue
			}
			t.Errorf("invalid crypto rand", singleTestCase)
		}
	}

	for _, singleTestCase := range testCase {
		result := CryptoRandDigit(singleTestCase)
		if len(result) != singleTestCase {
			t.Errorf("assert false! %v != %v", len(result), singleTestCase)
		}
		for _, singleResult := range result {
			if singleResult >= '0' && singleResult <= '9' {
				continue
			}
			t.Errorf("invalid crypto rand digit", singleTestCase)
		}
	}

	t.Log(CryptoRand(32))
}
