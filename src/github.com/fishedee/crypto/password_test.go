package crypto

import (
	"testing"
)

func TestCryptoPassword(t *testing.T) {
	testCase := []struct {
		origin  string
		origin2 string
		target  bool
	}{
		{"password123456", "$2y$10$6shpZxgObxSefyVVfqWbguIwnfiRNST0n2RWrJoqrx1aF2eOlzz7a", true},
		{"password123456", "$2a$123", false},
		{"123", "$2y$123", false},
		{"你好", "$2y$123", false},
	}
	for singleTestCaseIndex, singleTestCase := range testCase {
		result, err := PasswordVerify([]byte(singleTestCase.origin), singleTestCase.origin2)
		if err != nil {
			t.Errorf("check fail! error : %v, case : %v", err, singleTestCaseIndex)
		}
		if result != singleTestCase.target {
			t.Errorf("check fail! case : %v", singleTestCaseIndex)
		}
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		hashResult, err := PasswordHash([]byte(singleTestCase.origin), BCRYPT)
		if err != nil {
			t.Errorf("hash fail! error : %v, case : %v", err, singleTestCaseIndex)
		}
		result, err := PasswordVerify([]byte(singleTestCase.origin), hashResult)
		if err != nil {
			t.Errorf("check fail! error : %v, case : %v", err, singleTestCaseIndex)
		}
		if result != true {
			t.Errorf("check fail! case : %v", singleTestCaseIndex)
		}
	}
}
