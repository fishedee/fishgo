package router

import (
	"fmt"
	. "github.com/fishedee/assert"
	"net/http"
	"testing"
)

type fakeWriter struct {
	result string
}

func (this *fakeWriter) Header() http.Header {
	return http.Header{}
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
			},
			[]interface{}{
				[]interface{}{"/", "/_3"},
				[]interface{}{"/a/b/c", "/a/b/c_1"},
				[]interface{}{"/a/e/f", "/a/e/:mcId_1"},
				[]interface{}{"/a/c/d", "/a/:userId/:mcId_1"},
				[]interface{}{"/a/d", "/a/d_3"},
				[]interface{}{"/b", "/_3"},
			},
		},
		//大小写不敏感
		{
			[]interface{}{
				[]interface{}{1, "/Ab/b"},
				[]interface{}{1, "/BC"},
			},
			[]interface{}{
				[]interface{}{"/aB/b", "/Ab/b_1"},
				[]interface{}{"/ab/b", "/Ab/b_1"},
				[]interface{}{"/Bc", "/BC_1"},
				[]interface{}{"/bc", "/BC_1"},
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
				[]interface{}{"///", "_1"},
				[]interface{}{" / / / ", "_1"},
				[]interface{}{"/b/c", "b/c_1"},
				[]interface{}{"b/c", "b/c_1"},
				[]interface{}{"b/c//", "b/c_1"},
				[]interface{}{"/b/c/a", "b/c//a_1"},
				[]interface{}{"/b//c/a/", "b/c//a_1"},
				[]interface{}{"/b/g", "//b///g//_1"},
				[]interface{}{"/b/g//", "//b///g//_1"},
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
				routerFactory.addRoute(routerMethod.GET, priority, path, handler)
			}(singleInsertData.([]interface{}))
		}

		router := routerFactory.Create()
		for findDataIndex, singleFindData := range singleTestCase.findData.([]interface{}) {
			func(singleFindData []interface{}) {
				r, _ := http.NewRequest("GET", singleFindData[0].(string), nil)
				w := &fakeWriter{}
				router.ServeHttp(w, r)
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
		entrys := routerMethod.Entrys()
		for i := routerMethod.HEAD; i <= routerMethod.PATCH; i++ {
			method := entrys[i]
			r, _ := http.NewRequest(method, "/a", nil)
			w := &fakeWriter{}
			router.ServeHttp(w, r)
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
		router.ServeHttp(w, r)
		AssertEqual(t, w.Read(), singleTestCase.data, index)
	}
}

func TestRouterUrlPrefixParam(t *testing.T) {
	insertData := []string{
		"/",
		"/a",
		"/a/:userId",
		"/a/:userId/:typeId",
		"/b",
		"/b/mc/:fishId",
		"/b/:typeId/:userId",
	}
	findData := []struct {
		url   string
		param map[string]string
	}{
		{"/", nil},
		{"/a", nil},
		{"/b", nil},
		{"/a/mc", map[string]string{
			"userId": "mc",
		}},
		{"/a/mc/jk", map[string]string{
			"userId": "mc",
			"typeId": "jk",
		}},
		{"/b/mc/jk", map[string]string{
			"fishId": "jk",
		}},
		{"/b/bj/jk", map[string]string{
			"typeId": "bj",
			"userId": "jk",
		}},
	}

	routerFactory := NewRouterFactory()
	check := make(chan map[string]string, 10)
	for _, data := range insertData {
		routerFactory.GET(data, func(w http.ResponseWriter, r *http.Request, param map[string]string) {
			check <- param
		})
	}
	router := routerFactory.Create()
	for _, data := range findData {
		r, _ := http.NewRequest("GET", data.url, nil)
		w := &fakeWriter{}
		router.ServeHttp(w, r)
		AssertEqual(t, len(check), 1)
		AssertEqual(t, <-check, data.param)
	}
}

func TestRouterMiddleware(t *testing.T) {
	newMiddleware := func(data string) RouterMiddleware {
		return func(handler []interface{}) interface{} {
			last := handler[len(handler)-1].(func(w http.ResponseWriter, r *http.Request, param map[string]string))
			return func(w http.ResponseWriter, r *http.Request, param map[string]string) {
				w.Write([]byte(data))
				last(w, r, param)
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
		router.ServeHttp(w, r)
		AssertEqual(t, w.Read(), singleTestCase.data)
	}
}
