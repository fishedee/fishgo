package test

import (
	. "github.com/fishedee/web"
	"testing"
)

type testFishGenStruct struct{}

func TestTest(t *testing.T) {
	RunTest(t, &testFishGenStruct{})
}
