package test

import (
	. "github.com/fishedee/web"
)

type ClientAoTest struct {
	Test
	ClientAo ClientLoginAoModel
}

func (this *ClientAoTest) TestBasic() {
	//没有登录
	this.AssertEqual(this.ClientAo.IsLogin(), false)

	//错误登录
	this.ClientAo.Login("fish", "123dd")
	this.AssertEqual(this.ClientAo.IsLogin(), false)

	//正确登录
	this.ClientAo.Login("fish", "123")
	this.AssertEqual(this.ClientAo.IsLogin(), true)

	//登出
	this.ClientAo.Logout()
	this.AssertEqual(this.ClientAo.IsLogin(), false)

	//reset用法
	this.ClientAo.Login("fish", "123")
	this.AssertEqual(this.ClientAo.IsLogin(), true)
	this.RequestReset()
	this.AssertEqual(this.ClientAo.IsLogin(), false)
	this.ClientAo.Login("fish", "123")
	this.AssertEqual(this.ClientAo.IsLogin(), true)
}

func init() {
	InitTest(&ClientAoTest{})
}
