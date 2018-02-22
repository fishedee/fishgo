package session

import (
	"encoding/json"
	"github.com/astaxie/beego/session"
	. "github.com/fishedee/language"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionId() string
	Begin() error
	End()
}

type SessionFactory interface {
	Create(w http.ResponseWriter, r *http.Request) Session
}

type SessionConfig struct {
	Driver           string `json:"driver" config:"dirver"`
	CookieName       string `json:"cookieName" config:"name"`
	DisableSetCookie bool   `json:"-" config:"disablesetcookie"`
	GcLifeTime       int    `json:"gclifetime" config:"gclifttime"`
	Secure           bool   `json:"secure" config:"secure"`
	CookieLifeTime   int    `json:"cookieLifeTime" config:"cookielifetime"`
	ProviderConfig   string `json:"providerConfig" config:"savepath"`
	Domain           string `json:"domain" config:"domain"`
	SessionIdLength  int    `json:"sessionIdLength" config:"length"`
}

type sessionFactoryImplement struct {
	manager *session.Manager
	config  SessionConfig
}

func NewSessionFactory(config SessionConfig) (SessionFactory, error) {
	if config.Driver == "" {
		return nil, nil
	}
	if config.CookieName == "" {
		config.CookieName = "session"
	}
	if config.CookieLifeTime == 0 {
		config.CookieLifeTime = 3600
	}
	if config.GcLifeTime == 0 {
		config.GcLifeTime = 36000
	}
	configMap := ArrayToMap(config, "json").(map[string]interface{})
	configMap["enableSetCookie"] = !config.DisableSetCookie
	result, err := json.Marshal(configMap)
	if err != nil {
		return nil, err
	}

	sessionManager, err := session.NewManager(config.Driver, string(result))
	if err != nil {
		return nil, err
	}
	go sessionManager.GC()

	return &sessionFactoryImplement{
		manager: sessionManager,
		config:  config,
	}, nil
}

func (this *sessionFactoryImplement) Create(w http.ResponseWriter, r *http.Request) Session {
	return newSession(this.manager, this.config, w, r)
}

type sessionImplement struct {
	manager *session.Manager
	config  SessionConfig
	store   session.Store
	w       http.ResponseWriter
	r       *http.Request
}

func newSession(manager *session.Manager, config SessionConfig, w http.ResponseWriter, r *http.Request) Session {
	result := &sessionImplement{
		manager: manager,
		config:  config,
		w:       w,
		r:       r,
	}
	return result
}

func (this *sessionImplement) Set(key, value interface{}) error {
	if this.store == nil {
		panic("you should begin session first")
	}
	return this.store.Set(key, value)
}

func (this *sessionImplement) Get(key interface{}) interface{} {
	if this.store == nil {
		panic("you should begin session first")
	}
	return this.store.Get(key)
}

func (this *sessionImplement) Delete(key interface{}) error {
	if this.store == nil {
		panic("you should begin session first")
	}
	return this.store.Delete(key)
}

func (this *sessionImplement) SessionId() string {
	if this.store == nil {
		panic("you should begin session first")
	}
	return this.store.SessionID()
}

func (this *sessionImplement) Begin() error {
	if this.store != nil {
		panic("you should begin session already")
	}
	result, errOrgin := this.manager.SessionStart(this.w, this.r)
	if errOrgin != nil {
		return errOrgin
	}
	this.store = result

	//获取当前的cookie值
	cookie, err := this.r.Cookie(this.config.CookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil
	}

	//补充延续session时间的逻辑
	cookieValue := this.w.Header().Get("Set-Cookie")
	cookieName := this.config.CookieName
	if strings.Index(cookieValue, cookieName) != -1 {
		//已经设置过了
		return nil
	}
	cookie = &http.Cookie{
		Name:     this.config.CookieName,
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		Secure:   this.config.Secure,
		Domain:   this.config.Domain,
	}
	if this.config.CookieLifeTime > 0 {
		cookie.MaxAge = this.config.CookieLifeTime
		cookie.Expires = time.Now().Add(time.Duration(this.config.CookieLifeTime) * time.Second)
	}
	if this.config.DisableSetCookie == false {
		http.SetCookie(this.w, cookie)
	}
	return nil
}

func (this *sessionImplement) End() {
	if this.store == nil {
		panic("you should begin session first")
	}
	this.store.SessionRelease(this.w)
	this.store = nil
}
