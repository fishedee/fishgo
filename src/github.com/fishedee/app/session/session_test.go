package session

import (
	. "github.com/fishedee/assert"
	. "github.com/fishedee/language"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getHeaderInfo(header http.Header, name string) map[string]string {
	cookieInfo := Explode(header.Get("Set-Cookie"), ";")
	result := map[string]string{}
	for _, singleCookieInfo := range cookieInfo {
		info := Explode(singleCookieInfo, "=")
		if len(info) == 2 {
			result[info[0]] = info[1]
		} else {
			result[info[0]] = ""
		}
	}
	return result
}

func firstRequest(t *testing.T, sessionFactory SessionFactory) string {
	r, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	w := httptest.NewRecorder()

	session := sessionFactory.Create(w, r)
	session.Begin()
	session.Set("mc", "123")
	session.End()

	headerInfo := getHeaderInfo(w.Header(), "Set-Cookie")
	AssertEqual(t, len(headerInfo["fishmm"]) != 0, true)
	return headerInfo["fishmm"]
}

func secondRequest(t *testing.T, sessionFactory SessionFactory, sessionId string) {
	r, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Cookie", "fishmm="+sessionId)
	session := sessionFactory.Create(w, r)
	session.Begin()
	data := session.Get("mc")
	session.End()

	headerInfo := getHeaderInfo(w.Header(), "Set-Cookie")
	AssertEqual(t, headerInfo["fishmm"], sessionId)
	AssertEqual(t, data, "123")
}

func TestSession(t *testing.T) {
	sessionFactory, _ := NewSessionFactory(SessionConfig{
		Driver:     "memory",
		CookieName: "fishmm",
	})

	sessionId := firstRequest(t, sessionFactory)
	secondRequest(t, sessionFactory, sessionId)
	sessionId2 := firstRequest(t, sessionFactory)

	AssertEqual(t, sessionId != sessionId2, true)
}
