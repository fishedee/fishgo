package web

import (
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"time"
	"encoding/json"
	_ "github.com/fishedee/web/beego_session"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
)

type MySessionManagerConfig struct{
	CookieName      string `json:"cookieName"`
	EnableSetCookie bool   `json:"enableSetCookie,omitempty"`
	GcLifeTime      int  `json:"gclifetime"`
	Secure          bool   `json:"secure"`
	CookieLifeTime  int    `json:"cookieLifeTime"`
	ProviderConfig  string `json:"providerConfig"`
	Domain          string `json:"domain"`
	SessionIdLength int  `json:"sessionIdLength"`
}

type MySessionManager struct{
	*session.Manager
	config *MySessionManagerConfig
}

var Session *MySessionManager

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
	
	sessionlink := &MySessionManagerConfig{}
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

	sessionManager, err := session.NewManager(sessiondirver, string(result))
	if err != nil {
		panic(err)
	}
	go sessionManager.GC()

	Session = &MySessionManager{
		Manager:sessionManager,
		config:sessionlink,
	}
}

func (manager *MySessionManager) SessionStart(w http.ResponseWriter, r *http.Request) (session session.SessionStore, err error) {
	result,errOrgin := manager.Manager.SessionStart(w,r)
	if errOrgin != nil{
		return result,errOrgin
	}
	//获取当前的cookie值
	cookie, err := r.Cookie(manager.config.CookieName)
	if err != nil || cookie.Value == ""{
		return result,errOrgin
	}
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil{
		return result,errOrgin
	}

	//补充延续session时间的逻辑
	cookieValue := w.Header().Get("Set-Cookie")
	cookieName := manager.config.CookieName
	if strings.Index(cookieValue,cookieName) != -1{
		return result,err
	}
	cookie = &http.Cookie{
		Name:     manager.config.CookieName,
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		Secure:   manager.config.Secure,
		Domain:   manager.config.Domain,
	}
	if manager.config.CookieLifeTime > 0 {
		cookie.MaxAge = manager.config.CookieLifeTime
		cookie.Expires = time.Now().Add(time.Duration(manager.config.CookieLifeTime) * time.Second)
	}
	if manager.config.EnableSetCookie {
		http.SetCookie(w, cookie)
	}
	return result,errOrgin
}
