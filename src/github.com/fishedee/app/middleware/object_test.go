package middleware

import (
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/assert"
	"net/http"
	"testing"
)

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

func (this *testObject) POST_do6_Json(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("do6"))
}

func TestRouterObject(t *testing.T) {
	routerFactory := NewRouterFactory()
	routerFactory.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404"))
	})
	ObjectRouter(routerFactory, "/", &testObject{})
	ObjectRouter(routerFactory, "/mc", &testObject{})
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
		{"ANY", "/mc/do1", "do1"},
		{"ANY", "/mc/do2", "do2"},
		{"ANY", "/mc/do3", "do3"},
		{"ANY", "/mc", "do4"},
		{"GET", "/mc/do5", "do5"},
		{"POST", "/mc/do6", "do6"},
	}
	router := routerFactory.Create()
	for index, singleTestCase := range testCase {
		entrys := RouterMethod.Entrys()
		for i := RouterMethod.HEAD; i <= RouterMethod.PATCH; i++ {
			r, _ := http.NewRequest(entrys[i], singleTestCase.url, nil)
			w := &fakeWriter{}
			router.ServeHttp(w, r)
			if singleTestCase.method == "ANY" ||
				singleTestCase.method == entrys[i] {
				AssertEqual(t, w.Read(), singleTestCase.data, index)
			} else {
				AssertEqual(t, w.Read(), "do4", index)
			}
		}
	}
}
