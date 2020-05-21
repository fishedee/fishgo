package session

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type JwtTokenConfig struct {
	SecretKey      string `config:"secretkey"`
	CookieLifeTime int    `config:"cookielifetime"`
	CookieName     string `config:"cookiename"`
	Secure         bool   `config:"secure"`
	Domain         string `config:"domain"`
}

type jwtTokenFactory struct {
	config JwtTokenConfig
}

func NewJwtTokenFactory(config JwtTokenConfig) (SessionFactory, error) {
	if config.CookieName == "" {
		config.CookieName = "session"
	}
	if config.CookieLifeTime <= 0 {
		config.CookieLifeTime = 3600 * 24
	}
	return &jwtTokenFactory{
		config: config,
	}, nil
}

func (this *jwtTokenFactory) Create(w http.ResponseWriter, r *http.Request) Session {
	return newJwtToken(this.config, w, r)
}

type jwtToken struct {
	w         http.ResponseWriter
	r         *http.Request
	config    JwtTokenConfig
	hasModify bool
	claims    jwt.MapClaims
}

func newJwtToken(config JwtTokenConfig, w http.ResponseWriter, r *http.Request) Session {
	return &jwtToken{
		w:         w,
		r:         r,
		config:    config,
		claims:    nil,
		hasModify: false,
	}
}

func (this *jwtToken) Set(key string, value interface{}) error {
	if this.claims == nil {
		return errors.New("you should begin session first")
	}
	this.hasModify = true
	this.claims[key] = value
	return nil
}

func (this *jwtToken) MustSet(key string, value interface{}) {
	err := this.Set(key, value)
	if err != nil {
		panic(err)
	}
}

func (this *jwtToken) Get(key string) (interface{}, error) {
	if this.claims == nil {
		return nil, errors.New("you should begin session first")
	}
	return this.claims[key], nil
}

func (this *jwtToken) MustGet(key string) interface{} {
	value, err := this.Get(key)
	if err != nil {
		panic(err)
	}
	return value
}

func (this *jwtToken) Delete(key string) error {
	if this.claims == nil {
		return errors.New("you should begin session first")
	}
	this.hasModify = true
	delete(this.claims, key)
	return nil
}

func (this *jwtToken) MustDelete(key string) {
	err := this.Delete(key)
	if err != nil {
		panic(err)
	}
}

func (this *jwtToken) SessionId() string {
	return ""
}

func (this *jwtToken) getClaims(token *jwt.Token) jwt.MapClaims {
	if !token.Valid {
		return jwt.MapClaims{}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}
	}

	//校验过期时间
	expireUnixNanoStr, ok := claims["_expire"].(string)
	if !ok {
		return jwt.MapClaims{}
	}
	expireUnixNano, err := strconv.ParseInt(expireUnixNanoStr, 10, 64)
	if err != nil {
		return jwt.MapClaims{}
	}
	nowTimeNano := time.Now().UnixNano()
	if expireUnixNano < nowTimeNano {
		return jwt.MapClaims{}
	}

	//校验IP
	remoteIP, ok := claims["_remoteIP"].(string)
	if !ok {
		return jwt.MapClaims{}
	}
	if remoteIP != this.remoteIP() {
		return jwt.MapClaims{}
	}
	return claims

}

func (this *jwtToken) proxyAddr() []string {
	if ips := this.r.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func (this *jwtToken) remoteAddr() string {
	ips := this.proxyAddr()
	if len(ips) > 0 && ips[0] != "" {
		return ips[0]
	}
	return this.r.RemoteAddr
}

func (this *jwtToken) remoteIP() string {
	addr := this.remoteAddr()
	ip := strings.Split(addr, ":")
	if len(ip) > 0 {
		if ip[0] != "" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (this *jwtToken) Begin() error {
	//获取cookie中的数值
	cookieValue := ""
	cookie, err := this.r.Cookie(this.config.CookieName)
	if err == nil {
		cookieValue = cookie.Value
	}

	//读取数值中的map
	token, err := jwt.Parse(cookieValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(this.config.SecretKey), nil
	})
	if err == nil {
		this.claims = this.getClaims(token)
	} else {
		this.claims = jwt.MapClaims{}
	}
	this.hasModify = false
	return nil
}
func (this *jwtToken) MustBegin() {
	err := this.Begin()
	if err != nil {
		panic(err)
	}
}

func (this *jwtToken) Commit() error {
	if this.hasModify == false {
		this.claims = nil
		return nil
	}

	//将map转换为数值
	expires := time.Now().Add(time.Duration(this.config.CookieLifeTime) * time.Second)
	token := jwt.New(jwt.SigningMethodHS256)
	this.claims["_expire"] = strconv.FormatInt(expires.UnixNano(), 10)
	this.claims["_remoteIP"] = this.remoteIP()
	token.Claims = this.claims
	tokenString, err := token.SignedString([]byte(this.config.SecretKey))
	if err != nil {
		this.claims = nil
		return err
	}

	//将数值写入cookie
	cookie := &http.Cookie{
		Name:     this.config.CookieName,
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   this.config.Secure,
		Domain:   this.config.Domain,
	}
	cookie.MaxAge = this.config.CookieLifeTime
	cookie.Expires = expires
	http.SetCookie(this.w, cookie)
	this.claims = nil
	return nil

}
func (this *jwtToken) MustCommit() {
	err := this.Commit()
	if err != nil {
		panic(err)
	}
}
