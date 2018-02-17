package router

import (
	"fmt"
	"net/http"
	"testing"
)

type fakeLog struct {
}

func (this *fakeLog) Critical(format string, v ...interface{}) {
	fmt.Printf("[Critical] "+format+"\n", v...)
}

func (this *fakeLog) Error(format string, v ...interface{}) {
	fmt.Printf("[Error] "+format+"\n", v...)
}

func (this *fakeLog) Debug(format string, v ...interface{}) {
	fmt.Printf("[Debug] "+format+"\n", v...)
}

func TestLog(t *testing.T) {
	factory := NewRouterFactory()
	factory.Use(NewLogMiddleware(&fakeLog{}))
	factory.GET("/a", func(w http.ResponseWriter, r *http.Request) {

	})
	factory.GET("/b", func(w http.ResponseWriter, r *http.Request) {
		panic("Hello World")
	})

	router := factory.Create()
	testCase := []string{"/a", "/b", "/c"}
	for _, url := range testCase {
		r, _ := http.NewRequest("GET", url, nil)
		w := &fakeWriter{}
		router.ServeHttp(w, r)
	}

}
