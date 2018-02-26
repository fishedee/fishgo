package router

import (
	"fmt"
	. "github.com/fishedee/assert"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

type fakeWriter struct {
	header http.Header
	result string
}

func (this *fakeWriter) Header() http.Header {
	if this.header == nil {
		this.header = http.Header{}
	}
	return this.header
}

func (this *fakeWriter) WriteHeader(status int) {

}

func (this *fakeWriter) Write(data []byte) (int, error) {
	this.result += string(data)
	return len(data), nil
}

func (this *fakeWriter) Read() string {
	return this.result
}

func TestRouterUrl(t *testing.T) {
	testCase := []struct {
		insertData interface{}
		findData   interface{}
	}{
		//精确url匹配
		{
			[]interface{}{
				[]interface{}{1, "/a"},
				[]interface{}{1, "/b"},
				[]interface{}{1, "/ab"},
				[]interface{}{1, "/ab/c"},
				[]interface{}{1, "/ab/u"},
			},
			[]interface{}{
				[]interface{}{"/a", "/a_1"},
				[]interface{}{"/b", "/b_1"},
				[]interface{}{"/c", "404"},
				[]interface{}{"/ab", "/ab_1"},
				[]interface{}{"/abu", "404"},
				[]interface{}{"/ab/c", "/ab/c_1"},
				[]interface{}{"/ab/u", "/ab/u_1"},
				[]interface{}{"/ab/cu", "404"},
			},
		},
		//前缀url匹配
		{
			[]interface{}{
				[]interface{}{1, "/a"},
				[]interface{}{1, "/a/:userId"},
				[]interface{}{1, "/ab"},
				[]interface{}{1, "/ab/:userId"},
				[]interface{}{1, "/ab/:userId/:likeId"},
				[]interface{}{1, "/ab/10001/:likeId"},
				[]interface{}{1, "/ab/10001/10002"},
				[]interface{}{1, "/ab/10001"},
				[]interface{}{1, "/ab/u"},
			},
			[]interface{}{
				[]interface{}{"/a", "/a_1"},
				[]interface{}{"/a/mc", "/a/:userId_1"},
				[]interface{}{"/a/mc/jk", "404"},
				[]interface{}{"/ab", "/ab_1"},
				[]interface{}{"/ab/123", "/ab/:userId_1"},
				[]interface{}{"/ab/10001", "/ab/10001_1"},
				[]interface{}{"/ab/ck/mj", "/ab/:userId/:likeId_1"},
				[]interface{}{"/ab/10001/mj", "/ab/10001/:likeId_1"},
				[]interface{}{"/ab/10001/10002", "/ab/10001/10002_1"},
				[]interface{}{"/ab/10001/10002/mc", "404"},
				[]interface{}{"/ab/mc/jk/mc", "404"},
				[]interface{}{"/aj/123", "404"},
			},
		},
		//前缀静态匹配
		{
			[]interface{}{
				[]interface{}{3, "/"},
				[]interface{}{3, "/a"},
				[]interface{}{3, "/b"},
				[]interface{}{3, "/a/10001"},
				[]interface{}{3, "/a/10002"},
				[]interface{}{3, "/b/10003/10004"},
			},
			[]interface{}{
				[]interface{}{"/", "/_3"},
				[]interface{}{"/a", "/a_3"},
				[]interface{}{"/b", "/b_3"},
				[]interface{}{"/a/mc", "/a_3"},
				[]interface{}{"/a/mc/jk", "/a_3"},
				[]interface{}{"/a/10001", "/a/10001_3"},
				[]interface{}{"/a/10001/gj", "/a/10001_3"},
				[]interface{}{"/a/10002/cd", "/a/10002_3"},
				[]interface{}{"/b/10001", "/b_3"},
				[]interface{}{"/b/10003", "/b_3"},
				[]interface{}{"/b/10003/10004", "/b/10003/10004_3"},
			},
		},
		//混合匹配
		{
			[]interface{}{
				[]interface{}{4, "/"},
				[]interface{}{3, "/"},

				[]interface{}{1, "/a/b/c"},
				[]interface{}{1, "/a/b/c"},
				[]interface{}{3, "/a/b/c"},
				[]interface{}{4, "/a/b/c"},

				[]interface{}{1, "/a/e/:mcId"},
				[]interface{}{3, "/a/e/f"},
				[]interface{}{4, "/a/e/f"},

				[]interface{}{1, "/a/:userId/:mcId"},
				[]interface{}{3, "/a/c/d"},
				[]interface{}{4, "/a/c/d"},

				[]interface{}{3, "/a/d"},
				[]interface{}{4, "/a/d"},

				[]interface{}{1, "/mc/:id1/:id2"},
				[]interface{}{1, "/mc/10001/:id3"},

				[]interface{}{1, "/ck/:id"},
				[]interface{}{1, "/ck/cvhbu"},
				[]interface{}{1, "/ck/cvhbg"},
			},
			[]interface{}{
				[]interface{}{"/", "/_3"},
				[]interface{}{"/a/b/c", "/a/b/c_1"},
				[]interface{}{"/a/e/f", "/a/e/:mcId_1"},
				[]interface{}{"/a/c/d", "/a/:userId/:mcId_1"},
				[]interface{}{"/a/d", "/a/d_3"},
				[]interface{}{"/b", "/_3"},
				[]interface{}{"/mc/100014/5555", "/mc/:id1/:id2_1"},
				[]interface{}{"/ck/cvhbc", "/ck/:id_1"},
			},
		},
		//不正常url拼接
		{
			[]interface{}{
				[]interface{}{1, ""},
				[]interface{}{1, "b/c"},
				[]interface{}{1, "b/c//a"},
				[]interface{}{1, "//b///g//"},
			},
			[]interface{}{
				[]interface{}{"/", "_1"},
				[]interface{}{"", "_1"},
				[]interface{}{"//", "_1"},
				[]interface{}{"/b/c", "b/c_1"},
				[]interface{}{"b/c", "b/c_1"},
				[]interface{}{"b/c/", "b/c_1"},
				[]interface{}{"/b/c/a", "b/c//a_1"},
				[]interface{}{"/b/c/a/", "b/c//a_1"},
				[]interface{}{"/b/g", "//b///g//_1"},
				[]interface{}{"/b/g/", "//b///g//_1"},
				[]interface{}{"/b/g/?a=45", "//b///g//_1"},
			},
		},
	}
	for testCaseIndex, singleTestCase := range testCase {
		routerFactory := NewRouterFactory()

		routerFactory.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("404"))
		})

		for _, singleInsertData := range singleTestCase.insertData.([]interface{}) {
			func(singleInsertData []interface{}) {
				priority := singleInsertData[0].(int)
				path := singleInsertData[1].(string)
				handler := func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(path + fmt.Sprintf("_%v", priority)))
				}
				routerFactory.addRoute(RouterMethod.GET, priority, path, handler)
			}(singleInsertData.([]interface{}))
		}

		router := routerFactory.Create()
		for findDataIndex, singleFindData := range singleTestCase.findData.([]interface{}) {
			func(singleFindData []interface{}) {
				r, _ := http.NewRequest("GET", singleFindData[0].(string), nil)
				w := &fakeWriter{}
				router.ServeHTTP(w, r)
				AssertEqual(t, w.Read(), singleFindData[1].(string), fmt.Sprintf("%v-%v", testCaseIndex, findDataIndex))
			}(singleFindData.([]interface{}))
		}
	}
}

func TestRouterMethod(t *testing.T) {
	testCase := []struct {
		insertData func(*RouterFactory, string, interface{}) *RouterFactory
		findData   string
	}{
		{
			(*RouterFactory).HEAD,
			"HEAD",
		},
		{
			(*RouterFactory).OPTIONS,
			"OPTIONS",
		},
		{
			(*RouterFactory).GET,
			"GET",
		},
		{
			(*RouterFactory).POST,
			"POST",
		},
		{
			(*RouterFactory).DELETE,
			"DELETE",
		},
		{
			(*RouterFactory).PUT,
			"PUT",
		},
		{
			(*RouterFactory).PATCH,
			"PATCH",
		},
		{
			(*RouterFactory).Any,
			"Any",
		},
	}
	for _, singleTestCase := range testCase {
		routerFactory := NewRouterFactory()
		routerFactory.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("404"))
		})
		singleTestCase.insertData(routerFactory, "/a", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("bingo"))
		})
		router := routerFactory.Create()
		entrys := RouterMethod.Entrys()
		for i := RouterMethod.HEAD; i <= RouterMethod.PATCH; i++ {
			method := entrys[i]
			r, _ := http.NewRequest(method, "/a", nil)
			w := &fakeWriter{}
			router.ServeHTTP(w, r)
			if singleTestCase.findData == "Any" ||
				method == singleTestCase.findData {
				AssertEqual(t, w.Read(), "bingo")
			} else {
				AssertEqual(t, w.Read(), "404")
			}
		}
	}
}

func TestRouterStatic(t *testing.T) {
	routerFactory := NewRouterFactory()
	routerFactory.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404"))
	})
	routerFactory.Static("/", "./testdata")
	routerFactory.Static("/mj", "./testdata")
	routerFactory.Group("/mc", func(routerFactory2 *RouterFactory) {
		routerFactory2.Static("/", "./testdata")
		routerFactory2.Static("/jk", "./testdata")
		routerFactory2.Group("/mu", func(routerFactory3 *RouterFactory) {
			routerFactory3.Static("/", "./testdata")
		})
	})

	router := routerFactory.Create()

	testCase := []struct {
		url  string
		data string
	}{
		{"/a.html", "hello a"},
		{"/b.html", "hello b"},
		{"/c/d.html", "hello d"},
		{"/e.html", "404"},
		{"/mj/a.html", "hello a"},
		{"/mj/b.html", "hello b"},
		{"/mj/c/d.html", "hello d"},
		{"/mj/e.html", "404"},
		{"/mc/a.html", "hello a"},
		{"/mc/b.html", "hello b"},
		{"/mc/c/d.html", "hello d"},
		{"/mc/e.html", "404"},
		{"/mc/mu/a.html", "hello a"},
		{"/mc/mu/b.html", "hello b"},
		{"/mc/mu/c/d.html", "hello d"},
		{"/mc/mu/e.html", "404"},
		{"/mc/jk/a.html", "hello a"},
		{"/mc/jk/b.html", "hello b"},
		{"/mc/jk/c/d.html", "hello d"},
		{"/mc/jk/e.html", "404"},
	}
	for index, singleTestCase := range testCase {
		r, _ := http.NewRequest("GET", singleTestCase.url, nil)
		w := &fakeWriter{}
		router.ServeHTTP(w, r)
		AssertEqual(t, w.Read(), singleTestCase.data, index)
	}
}

func TestRouterUrlPrefixParam(t *testing.T) {
	insertData := []string{
		"/",
		"/:mmId",
		"/a",
		"/a/:userId",
		"/a/:userId/:typeId",
		"/b",
		"/b/mc/:fishId",
		"/b/:typeId/:userId",
	}
	findData := []struct {
		url   string
		param RouterParam
	}{
		{"/", nil},
		{"/a", nil},
		{"/b", nil},
		{"/c", RouterParam{
			{"mmId", "c"},
		}},
		{"/a/mc", RouterParam{
			{"userId", "mc"},
		}},
		{"/a/mc/jk", RouterParam{
			{"userId", "mc"},
			{"typeId", "jk"},
		}},
		{"/b/mc/jk", RouterParam{
			{"fishId", "jk"},
		}},
		{"/b/bj/jk", RouterParam{
			{"typeId", "bj"},
			{"userId", "jk"},
		}},
		{"/b/mc/jk/", RouterParam{
			{"fishId", "jk"},
		}},
		{"/b/bj/jk/", RouterParam{
			{"typeId", "bj"},
			{"userId", "jk"},
		}},
	}

	routerFactory := NewRouterFactory()
	check := make(chan RouterParam, 10)
	for _, data := range insertData {
		routerFactory.GET(data, func(w http.ResponseWriter, r *http.Request, param RouterParam) {
			check <- param
		})
	}
	router := routerFactory.Create()
	for _, data := range findData {
		r, _ := http.NewRequest("GET", data.url, nil)
		w := &fakeWriter{}
		router.ServeHTTP(w, r)
		select {
		case result := <-check:
			AssertEqual(t, result, data.param, data)
		default:
			AssertEqual(t, false, true, data)
		}
	}
}

func TestRouterMiddleware(t *testing.T) {
	newMiddleware := func(data string) RouterMiddleware {
		return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
			handler := prev.Handler.(func(w http.ResponseWriter, r *http.Request, param RouterParam))
			return RouterMiddlewareContext{
				Data: prev.Data,
				Handler: func(w http.ResponseWriter, r *http.Request, param RouterParam) {
					w.Write([]byte(data))
					handler(w, r, param)
				},
			}
		}
	}
	doNothing := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fish"))
	}
	routerFactory := NewRouterFactory()
	routerFactory.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404"))
	})
	routerFactory.
		Use(newMiddleware("mid1_")).
		Use(newMiddleware("mid2_")).
		GET("/a", doNothing).
		GET("/b", doNothing).
		Group("/c", func(routerFactory2 *RouterFactory) {
			routerFactory2.
				Use(newMiddleware("mid3_")).
				GET("/d", doNothing).
				GET("/e", doNothing)
		}).
		Group("/", func(routerFactory2 *RouterFactory) {
			routerFactory2.
				Use(newMiddleware("mid4_")).
				GET("/f", doNothing).
				GET("/g", doNothing)
		})
	testCase := []struct {
		url  string
		data string
	}{
		{"/", "mid1_mid2_404"},
		{"/a", "mid1_mid2_fish"},
		{"/b", "mid1_mid2_fish"},
		{"/f", "mid1_mid2_mid4_fish"},
		{"/g", "mid1_mid2_mid4_fish"},
		{"/h", "mid1_mid2_404"},
		{"/c/d", "mid1_mid2_mid3_fish"},
		{"/c/e", "mid1_mid2_mid3_fish"},
	}

	router := routerFactory.Create()
	for _, singleTestCase := range testCase {
		r, _ := http.NewRequest("GET", singleTestCase.url, nil)
		w := &fakeWriter{}
		router.ServeHTTP(w, r)
		AssertEqual(t, w.Read(), singleTestCase.data)
	}
}

func benchmarkRouterBasic(b *testing.B, insertData []string, findData []string) {
	routerFactory := NewRouterFactory()
	doNothing := func(w http.ResponseWriter, r *http.Request, param RouterParam) {
	}
	routerFactory.NotFound(doNothing)
	for _, data := range insertData {
		routerFactory.GET(data, doNothing)
	}

	r, _ := http.NewRequest("GET", "", nil)
	w := &fakeWriter{}
	router := routerFactory.Create()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		single := findData[i%len(findData)]
		r.URL.Path = single
		router.ServeHTTP(w, r)
	}
}

func benchmarkGinBasic(b *testing.B, insertData []string, findData []string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	doNothing := func(c *gin.Context) {
	}
	for _, data := range insertData {
		router.GET(data, doNothing)
	}

	r, _ := http.NewRequest("GET", "", nil)
	w := &fakeWriter{}

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		single := findData[i%len(findData)]
		r.URL.Path = single
		router.ServeHTTP(w, r)
	}
}

func BenchmarkRouterShort(b *testing.B) {
	testUrl := "/abc"
	benchmarkRouterBasic(b, []string{testUrl}, []string{testUrl})
}

func BenchmarkGinShort(b *testing.B) {
	testUrl := "/abc"
	benchmarkGinBasic(b, []string{testUrl}, []string{testUrl})
}

func BenchmarkRouterLong(b *testing.B) {
	testUrl := "/abc/12312313/adf/asdf/asdf/asdf/sdaf/asdf/abc/12312313/adf/asdf/asdf/asdf/sdaf/asdf/"
	benchmarkRouterBasic(b, []string{testUrl}, []string{testUrl})
}

func BenchmarkGinLong(b *testing.B) {
	testUrl := "/abc/12312313/adf/asdf/asdf/asdf/sdaf/asdf/abc/12312313/adf/asdf/asdf/asdf/sdaf/asdf/"
	benchmarkGinBasic(b, []string{testUrl}, []string{testUrl})
}

func BenchmarkRouterParam(b *testing.B) {
	insertUrl := "/user/:userId"
	findUrl := "/user/123"
	benchmarkRouterBasic(b, []string{insertUrl}, []string{findUrl})
}

func BenchmarkGinParam(b *testing.B) {
	insertUrl := "/user/:userId"
	findUrl := "/user/123"
	benchmarkGinBasic(b, []string{insertUrl}, []string{findUrl})
}

func BenchmarkRouterParamLong(b *testing.B) {
	insertUrl := "/user/:userId"
	findUrl := "/user/123/adsfasdfadsfasdfasdfadsf/zcvczxcxzvzvcx"
	benchmarkRouterBasic(b, []string{insertUrl}, []string{findUrl})
}

func BenchmarkGinParamLong(b *testing.B) {
	insertUrl := "/user/:userId"
	findUrl := "/user/123/adsfasdfadsfasdfasdfadsf/zcvczxcxzvzvcx"
	benchmarkGinBasic(b, []string{insertUrl}, []string{findUrl})
}
