package number

import (
	. "github.com/fishedee/web"
	"testing"
)

type testFishGenStruct struct{}

func TestNumber(t *testing.T) {
	RunBeegoValidateTest(t, &testFishGenStruct{})
}
