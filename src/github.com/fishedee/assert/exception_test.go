package assert

import (
	. "github.com/fishedee/language"
	"testing"
)

func TestAssertException(t *testing.T) {
	//below test case should fail!
	AssertException(t, 1, "", func() {
		Throw(2, "")
	})
	AssertException(t, 1, "123", func() {
		Throw(1, "456")
	})
	AssertException(t, 1, "", func() {
	})
}
