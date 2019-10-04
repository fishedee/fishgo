package middleware

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/app/render"
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/app/session"
	. "github.com/fishedee/app/validator"
	. "github.com/fishedee/assert"
	"net/http"
	"testing"
)

type testInterface interface {
	Do1(w http.ResponseWriter, r *http.Request)
	Do2_Json(w http.ResponseWriter, r *http.Request)
	Do3_Html_Go(w http.ResponseWriter, r *http.Request)
	Any(w http.ResponseWriter, r *http.Request)
	GET_do5(w http.ResponseWriter, r *http.Request)
	POST_Do6_Json(w http.ResponseWriter, r *http.Request)
	Do7_Text(v Validator, s Session) interface{}
	NONE_Do8_Text(v Validator, s Session) interface{}
}
type testObject struct {
}

func (this *testObject) Do1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("do1"))
}

func (this *testObject) Do2_Json(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("do2"))
}

func (this *testObject) Do3_Html_Go(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("do3"))
}

func (this *testObject) Any(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("do4"))
}

func (this *testObject) GET_do5(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("do5"))
}

func (this *testObject) POST_Do6_Json(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("do6"))
}

func (this *testObject) Do7_Text(v Validator, s Session) interface{} {
	return "do7"
}

func (this *testObject) NONE_Do8_Text(v Validator, s Session) interface{} {
	return "do8"
}

func TestRouterObject(t *testing.T) {
	var testObjectInterface testInterface
	testObjectInterface = &testObject{}

	log, _ := NewLog(LogConfig{Driver: "console"})
	renderFactory, _ := NewRenderFactory(RenderConfig{})
	validatorFactory, _ := NewValidatorFactory(ValidatorConfig{})
	sessionFactory, _ := NewSessionFactory(SessionConfig{Driver: "memory", CookieName: "fishmm"})
	middleware := NewEasyMiddleware(log, validatorFactory, sessionFactory, renderFactory, nil)

	routerFactory := NewRouterFactory()
	routerFactory.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404"))
	})
	routerFactory.Use(middleware)

	ObjectRouter(routerFactory, "/", &testObject{})
	ObjectRouter(routerFactory, "/mc", &testObject{})
	ObjectRouter(routerFactory, "/mj", testObjectInterface)
	testCase := []struct {
		method string
		url    string
		data   string
	}{
		{"ANY", "/do1", "do1"},
		{"ANY", "/do2", "do2"},
		{"ANY", "/do3", "do3"},
		{"ANY", "/", "do4"},
		{"GET", "/do5", "do5"},
		{"POST", "/do6", "do6"},
		{"ANY", "/do7", "do7"},
		{"ANY", "/none", "404"},
		{"ANY", "/do8", "404"},
		{"ANY", "/mc/do1", "do1"},
		{"ANY", "/mc/do2", "do2"},
		{"ANY", "/mc/do3", "do3"},
		{"ANY", "/mc", "do4"},
		{"GET", "/mc/do5", "do5"},
		{"POST", "/mc/do6", "do6"},
		{"ANY", "/mc/do7", "do7"},
		{"ANY", "/mc/none", "404"},
		{"ANY", "/mc/do8", "404"},
		{"ANY", "/mj/do1", "do1"},
		{"ANY", "/mj/do2", "do2"},
		{"ANY", "/mj/do3", "do3"},
		{"ANY", "/mj", "do4"},
		{"GET", "/mj/do5", "do5"},
		{"POST", "/mj/do6", "do6"},
		{"ANY", "/mj/do7", "do7"},
		{"ANY", "/mj/none", "404"},
		{"ANY", "/mj/do8", "404"},
	}
	router := routerFactory.Create()
	for index, singleTestCase := range testCase {
		entrys := RouterMethod.Entrys()
		for i := RouterMethod.HEAD; i <= RouterMethod.PATCH; i++ {
			r, _ := http.NewRequest(entrys[i], singleTestCase.url, nil)
			w := &fakeWriter{}
			router.ServeHTTP(w, r)
			if singleTestCase.method == "ANY" ||
				singleTestCase.method == entrys[i] {
				AssertEqual(t, w.Read(), singleTestCase.data, index)
			} else {
				AssertEqual(t, w.Read(), "404", index)
			}
		}
	}
}
