package web

import (
	"encoding/json"
	"github.com/astaxie/beego/session"
	_ "github.com/fishedee/web/util_session"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SessionStore interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
	SessionRelease()
	Flush() error
}

type Session interface {
	WithContext(ctx Context) Session
	SessionStart() (session SessionStore, err error)
}

type SessionConfig struct {
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

type sessionImplement struct {
	*session.Manager
	config SessionConfig
	ctx    Context
}

type sessionStoreImplement struct {
	session.Store
	responseWriter http.ResponseWriter
}

func NewSession(config SessionConfig) (Session, error) {
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

	return &sessionImplement{
		Manager: sessionManager,
		config:  config,
	}, nil
}

func NewSessionFromConfig(configName string) (Session, error) {
	sessiondirver := globalBasic.Config.GetString(configName + "driver")
	sessionname := globalBasic.Config.GetString(configName + "name")
	sessiongclifttime := globalBasic.Config.GetString(configName + "gclifttime")
	sessioncookielifetime := globalBasic.Config.GetString(configName + "cookielifetime")
	sessionsavepath := globalBasic.Config.GetString(configName + "savepath")
	sessionsecure := globalBasic.Config.GetString(configName + "secure")
	sessiondomain := globalBasic.Config.GetString(configName + "domain")
	sessionlength := globalBasic.Config.GetString(configName + "length")

	sessionlink := SessionConfig{}
	sessionlink.Driver = sessiondirver
	sessionlink.CookieName = sessionname
	sessionlink.EnableSetCookie = true
	sessionlink.GcLifeTime, _ = strconv.Atoi(sessiongclifttime)
	sessionlink.Secure, _ = strconv.ParseBool(sessionsecure)
	sessionlink.CookieLifeTime, _ = strconv.Atoi(sessioncookielifetime)
	sessionlink.ProviderConfig = sessionsavepath
	sessionlink.Domain = sessiondomain
	sessionlink.SessionIdLength, _ = strconv.Atoi(sessionlength)

	return NewSession(sessionlink)
}

func newSessionStoreImplement(store session.Store, responseWriter http.ResponseWriter) SessionStore {
	result := sessionStoreImplement{
		Store:          store,
		responseWriter: responseWriter,
	}
	return &result
}

func (manager *sessionImplement) WithContext(ctx Context) Session {
	result := *manager
	result.ctx = ctx
	return &result
}

func (manager *sessionImplement) SessionStart() (session SessionStore, err error) {
	w := manager.ctx.GetRawResponseWriter().(http.ResponseWriter)
	r := manager.ctx.GetRawRequest().(*http.Request)

	result, errOrgin := manager.Manager.SessionStart(w, r)
	if errOrgin != nil {
		return newSessionStoreImplement(result, w), errOrgin
	}
	//获取当前的cookie值
	cookie, err := r.Cookie(manager.config.CookieName)
	if err != nil || cookie.Value == "" {
		return newSessionStoreImplement(result, w), errOrgin
	}
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return newSessionStoreImplement(result, w), errOrgin
	}

	//补充延续session时间的逻辑
	cookieValue := w.Header().Get("Set-Cookie")
	cookieName := manager.config.CookieName
	if strings.Index(cookieValue, cookieName) != -1 {
		return newSessionStoreImplement(result, w), err
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
	return newSessionStoreImplement(result, w), errOrgin
}

func (this *sessionStoreImplement) SessionRelease() {
	this.Store.SessionRelease(this.responseWriter)
}
