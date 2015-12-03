package web

import (
	"strconv"
	"encoding/json"
	_ "github.com/fishedee/web/beego_session"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
)

var Session *session.Manager

func init() {
	sessiondirver := beego.AppConfig.String("fishsessiondriver")
	sessionname := beego.AppConfig.String("fishsessionname")
	sessiongclifttime := beego.AppConfig.String("fishsessiongclifttime")
	sessioncookielifetime := beego.AppConfig.String("fishsessioncookielifetime")
	sessionsavepath := beego.AppConfig.String("fishsessionsavepath")
	sessionsecure := beego.AppConfig.String("fishsessionsecure")
	sessiondomain := beego.AppConfig.String("fishsessiondomain")
	sessionlength := beego.AppConfig.String("fishsessionlength")

	if sessiondirver == ""{
		return
	}
	var sessionlink struct{
		CookieName      string `json:"cookieName"`
		EnableSetCookie bool   `json:"enableSetCookie,omitempty"`
		GcLifeTime      int  `json:"gclifetime"`
		Secure          bool   `json:"secure"`
		CookieLifeTime  int    `json:"cookieLifeTime"`
		ProviderConfig  string `json:"providerConfig"`
		Domain          string `json:"domain"`
		SessionIdLength int  `json:"sessionIdLength"`
	}
	sessionlink.CookieName = sessionname
	sessionlink.EnableSetCookie = true
	sessionlink.GcLifeTime,_ = strconv.Atoi(sessiongclifttime)
	sessionlink.Secure,_ = strconv.ParseBool(sessionsecure)
	sessionlink.CookieLifeTime,_ = strconv.Atoi(sessioncookielifetime)
	sessionlink.ProviderConfig = sessionsavepath
	sessionlink.Domain = sessiondomain
	sessionlink.SessionIdLength,_ = strconv.Atoi(sessionlength)

	result,err := json.Marshal(sessionlink)
	if err != nil{
		panic(err)
	}

	Session, err = session.NewManager(sessiondirver, string(result))
	if err != nil {
		panic(err)
	}
	go Session.GC()
}
