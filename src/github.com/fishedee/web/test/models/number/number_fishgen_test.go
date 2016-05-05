package number

import (
	. "github.com/fishedee/web"
	"testing"
)

type testFishGenStruct struct{}

func TestNumber(t *testing.T) {
	RunTest(t, &testFishGenStruct{})
}
