package util

import (
	"encoding/json"
	"github.com/astaxie/beego/session"
	. "github.com/fishedee/util"
	_ "github.com/fishedee/util_session"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SessionManagerConfig struct {
	Driver          string `json:driver`
	CookieName      string `json:"cookieName"`
	EnableSetCookie bool   `json:"enableSetCookie,omitempty"`
	GcLifeTime      int    `json:"gclifetime"`
	Secure          bool   `json:"secure"`
	CookieLifeTime  int    `json:"cookieLifeTime"`
	ProviderConfig  string `json:"providerConfig"`
	Domain          string `json:"domain"`
	SessionIdLength int    `json:"sessionIdLength"`
}

type SessionManager struct {
	*session.Manager
	config SessionManagerConfig
}

var newSessionManagerMemory *MemoryFunc
var newSessionManagerFromConfigMemory *MemoryFunc

func init() {
	var err error
	newSessionManagerMemory, err = NewMemoryFunc(newSessionManager, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
	newSessionManagerFromConfigMemory, err = NewMemoryFunc(newSessionManagerFromConfig, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
}

func newSessionManager(config SessionManagerConfig) (*SessionManager, error) {
	if config.Driver == "" {
		return nil, nil
	}
	if config.CookieName == "" {
		config.CookieName = "beego_session"
	}
	if config.CookieLifeTime == 0 {
		config.CookieLifeTime = 3600
	}
	if config.GcLifeTime == 0 {
		config.GcLifeTime = 3600
	}
	result, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	sessionManager, err := session.NewManager(config.Driver, string(result))
	if err != nil {
		return nil, err
	}
	go sessionManager.GC()

	return &SessionManager{
		Manager: sessionManager,
		config:  config,
	}, nil
}

func NewSessionManager(config SessionManagerConfig) (*SessionManager, error) {
	result, err := newSessionManagerMemory.Call(config)
	return result.(*SessionManager), err
}

func newSessionManagerFromConfig(configName string) (*SessionManager, error) {
	sessiondirver := globalBasic.Config.String(configName + "driver")
	sessionname := globalBasic.Config.String(configName + "name")
	sessiongclifttime := globalBasic.Config.String(configName + "gclifttime")
	sessioncookielifetime := globalBasic.Config.String(configName + "cookielifetime")
	sessionsavepath := globalBasic.Config.String(configName + "savepath")
	sessionsecure := globalBasic.Config.String(configName + "secure")
	sessiondomain := globalBasic.Config.String(configName + "domain")
	sessionlength := globalBasic.Config.String(configName + "length")

	sessionlink := SessionManagerConfig{}
	sessionlink.Driver = sessiondirver
	sessionlink.CookieName = sessionname
	sessionlink.EnableSetCookie = true
	sessionlink.GcLifeTime, _ = strconv.Atoi(sessiongclifttime)
	sessionlink.Secure, _ = strconv.ParseBool(sessionsecure)
	sessionlink.CookieLifeTime, _ = strconv.Atoi(sessioncookielifetime)
	sessionlink.ProviderConfig = sessionsavepath
	sessionlink.Domain = sessiondomain
	sessionlink.SessionIdLength, _ = strconv.Atoi(sessionlength)

	return NewSessionManager(sessionlink)
}

func NewSessionManagerFromConfig(configName string) (*SessionManager, error) {
	result, err := newSessionManagerFromConfigMemory.Call(configName)
	return result.(*SessionManager), err
}

func (manager *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) (session session.Store, err error) {
	result, errOrgin := manager.Manager.SessionStart(w, r)
	if errOrgin != nil {
		return result, errOrgin
	}
	//获取当前的cookie值
	cookie, err := r.Cookie(manager.config.CookieName)
	if err != nil || cookie.Value == "" {
		return result, errOrgin
	}
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return result, errOrgin
	}

	//补充延续session时间的逻辑
	cookieValue := w.Header().Get("Set-Cookie")
	cookieName := manager.config.CookieName
	if strings.Index(cookieValue, cookieName) != -1 {
		return result, err
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
	return result, errOrgin
}
