package test

import (
	"fmt"
	. "github.com/fishedee/web"
	"strconv"
)

type ClientLoginAoModel struct {
	Model
}

func (this *ClientLoginAoModel) IsLogin() bool {
	sess, err := this.Session.SessionStart()
	if err != nil {
		panic("session启动失败")
	}
	defer sess.SessionRelease()

	clientId := sess.Get("clientId")
	clientIdString := fmt.Sprintf("%v", clientId)
	clientIdInt, err := strconv.Atoi(clientIdString)
	if err == nil && clientIdInt >= 10000 {
		return true
	} else {
		return false
	}
}

func (this *ClientLoginAoModel) Logout() {
	sess, err := this.Session.SessionStart()
	if err != nil {
		panic("session启动失败！")
	}
	defer sess.SessionRelease()

	sess.Set("clientId", 0)
}

func (this *ClientLoginAoModel) Login(name string, password string) bool {
	if name != "fish" || password != "123" {
		return false
	}
	sess, err := this.Session.SessionStart()
	if err != nil {
		panic("session启动失败！")
	}
	defer sess.SessionRelease()

	sess.Set("clientId", 10001)
	return true
}
