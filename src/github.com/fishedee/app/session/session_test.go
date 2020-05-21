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
	session.MustBegin()
	session.MustSet("mc", "123")
	session.MustCommit()

	headerInfo := getHeaderInfo(w.Header(), "Set-Cookie")
	AssertEqual(t, len(headerInfo["fishmm"]) != 0, true)
	return headerInfo["fishmm"]
}

func secondRequest(t *testing.T, sessionFactory SessionFactory, sessionId string) {
	r, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Cookie", "fishmm="+sessionId)
	session := sessionFactory.Create(w, r)
	session.MustBegin()
	data := session.MustGet("mc")
	session.MustCommit()

	//sessionId在jwtToken中不是每次都会改变的
	//headerInfo := getHeaderInfo(w.Header(), "Set-Cookie")
	//AssertEqual(t, headerInfo["fishmm"], sessionId)
	AssertEqual(t, data, "123")
}

func ATestSession(t *testing.T) {
	sessionFactory, _ := NewSessionFactory(SessionConfig{
		Driver:     "memory",
		CookieName: "fishmm",
	})
	jwtTokenFactory, _ := NewJwtTokenFactory(JwtTokenConfig{
		SecretKey:  "123",
		CookieName: "fishmm",
	})
	testCase := []SessionFactory{
		sessionFactory,
		jwtTokenFactory,
	}

	for _, sessionFactory := range testCase {
		sessionId := firstRequest(t, sessionFactory)
		secondRequest(t, sessionFactory, sessionId)
		sessionId2 := firstRequest(t, sessionFactory)

		t.Logf("%v,%v", sessionId, sessionId2)
		AssertEqual(t, sessionId != sessionId2, true)
	}

}

func TestSession2(t *testing.T) {
	jwtTokenFactory, _ := NewJwtTokenFactory(JwtTokenConfig{
		SecretKey:  "123",
		CookieName: "fishmm",
	})

	//第一次测试
	r1, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	r1.RemoteAddr = "192.168.5.2"
	w1 := httptest.NewRecorder()

	session := jwtTokenFactory.Create(w1, r1)
	session.MustBegin()
	session.MustSet("mc", "456")
	session.MustCommit()

	cookie := getHeaderInfo(w1.Header(), "Set-Cookie")
	sessionId := cookie["fishmm"]

	//第二次测试，不同IP无法获取到登录态
	r2, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	r2.RemoteAddr = "192.168.5.3"
	w2 := httptest.NewRecorder()

	r2.Header.Set("Cookie", "fishmm="+sessionId)
	session2 := jwtTokenFactory.Create(w2, r2)
	session2.MustBegin()
	data2 := session2.MustGet("mc")
	session2.MustCommit()

	AssertEqual(t, data2, nil)

	//第三次测试，相同IP无法获取到登录态
	r3, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	r3.RemoteAddr = "192.168.5.2"
	w3 := httptest.NewRecorder()

	r3.Header.Set("Cookie", "fishmm="+sessionId)
	session3 := jwtTokenFactory.Create(w3, r3)
	session3.MustBegin()
	data3 := session3.MustGet("mc")
	session3.MustCommit()

	AssertEqual(t, data3, "456")

	//第四次测试，取的是X-Forwarded-For而不是RemoteIP
	r4, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	r4.RemoteAddr = "192.168.5.4"
	r4.Header.Set("X-Forwarded-For", "192.168.5.2")
	w4 := httptest.NewRecorder()

	r4.Header.Set("Cookie", "fishmm="+sessionId)
	session4 := jwtTokenFactory.Create(w4, r4)
	session4.MustBegin()
	data4 := session4.MustGet("mc")
	session4.MustCommit()

	AssertEqual(t, data4, "456")
}
