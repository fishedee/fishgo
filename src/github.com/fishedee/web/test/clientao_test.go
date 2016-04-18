package test

import (
	"testing"
)

type ClientAoTest struct {
	ClientLoginAoModel
	ClientAo ClientLoginAoModel
}

func (this *ClientAoTest) TestBasic() {
	//没有
	this.AssertEqual(this.RequestGetCookie("clientId"), "")
}

func TestClientAo(t *testing.T) {
	InitTest(t, &ClientAoTest{})
}
