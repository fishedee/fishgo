package test

import (
	. "github.com/fishedee/web"
)

type clientAoTest struct {
	Test
	clientAo ClientLoginAoModel
}

func (this *ClientAoTest) TestBasic() {
	//没有登录
	this.AssertEqual(this.clientAo.IsLogin(), false)

	//错误登录
	this.clientAo.Login("fish", "123dd")
	this.AssertEqual(this.clientAo.IsLogin(), false)

	//正确登录
	this.clientAo.Login("fish", "123")
	this.AssertEqual(this.clientAo.IsLogin(), true)

	//登出
	this.clientAo.Logout()
	this.AssertEqual(this.clientAo.IsLogin(), false)

	//reset用法
	this.clientAo.Login("fish", "123")
	this.AssertEqual(this.clientAo.IsLogin(), true)
	this.RequestReset()
	this.AssertEqual(this.clientAo.IsLogin(), false)
	this.clientAo.Login("fish", "123")
	this.AssertEqual(this.clientAo.IsLogin(), true)
}
