package middleware

import (
	. "github.com/fishedee/app/cors"
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/assert"
	"net/http"
	"testing"
)

func TestCors(t *testing.T) {
	cors, _ := NewCors(CorsConfig{
		//AllowedOrigins:     []string{"*"},
		OptionsPassthrough: true,
	})

	factory := NewRouterFactory()
	factory.Use(NewCorsMiddleware(cors))
	factory.GET("/a", func(w http.ResponseWriter, r *http.Request) {

	})

	router := factory.Create()
	testCase := []struct {
		method    string
		reqHeader map[string]string
		resHeader http.Header
	}{
		{"OPTIONS", map[string]string{
			"Origin":                         "http://foobar.com",
			"Access-Control-Request-Method":  "GET",
			"Access-Control-Request-Headers": "X-Header-2, X-HEADER-1",
		}, map[string][]string{
			"Vary": []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		}},
		{"GET", map[string]string{
			"Origin": "http://foobar.com",
		}, map[string][]string{
			"Vary": []string{"Origin"},
			"Access-Control-Allow-Origin": []string{"*"},
		}},
	}
	for _, singleTestCase := range testCase {

		r, _ := http.NewRequest(singleTestCase.method, "http://example.com/foo", nil)
		for name, value := range singleTestCase.reqHeader {
			r.Header.Add(name, value)
		}
		w := &fakeWriter{}
		router.ServeHTTP(w, r)
		AssertEqual(t, w.Header(), singleTestCase.resHeader)
	}
}
