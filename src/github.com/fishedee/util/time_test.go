package util

import (
	. "github.com/fishedee/assert"
	"testing"
	"time"
)

//Author:Edward
func TestIsTimeZero(t *testing.T) {

	local, errData := time.LoadLocation("Local")
	new_York, errData := time.LoadLocation("America/New_York")

	if errData != nil {
		panic(errData)
	}
	testCase := []struct {
		origin time.Time
		target bool
		err    error
	}{
		{
			time.Now(),
			false,
			nil,
		},
		{
			time.Date(1, time.January, 1, 0, 0, 0, 0, local),
			true,
			nil,
		},
		{
			time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			true,
			nil,
		},
		{
			time.Date(1, time.January, 1, 0, 0, 0, 0, new_York),
			false,
			nil,
		},
	}

	for singleKey, singleData := range testCase {
		result, err := IsTimeZero(singleData.origin)
		AssertEqual(t, result, singleData.target, singleKey)
		AssertEqual(t, singleData.err, err, singleKey)
	}
}
