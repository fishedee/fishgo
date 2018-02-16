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
				[]interface{}{"/", "/_4"},
				[]interface{}{"/a/b/c", "/a/b/c_1"},
				[]interface{}{"/a/e/f", "/a/e/:mcId_1"},
				[]interface{}{"/a/c/d", "/a/:userId/:mcId_1"},
				[]interface{}{"/a/d", "/a/d_3"},
				[]interface{}{"/b", "/_4"},
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

func TestRouterObject(t *testing.T) {
	routerFactory := NewRouterFactory()
	routerFactory.GET("/a", &testObject{})
	routerFactory.GET("/", &testObject{})
	testCase := []struct {
		url  string
		data string
	}{
		{"/do1", "do1"},
		{"/do2", "do2"},
		{"/do3", "do3"},
		{"/a/do1", "do1"},
		{"/a/do2", "do2"},
		{"/a/do3", "do3"},
	}
	router := routerFactory.Create()
	for _, singleTestCase := range testCase {
		r, _ := http.NewRequest("GET", singleTestCase.url, nil)
		w := &fakeWriter{}
		router.ServeHttp(w, r)
		AssertEqual(t, w.Read(), singleTestCase.data)
	}
}

func TestRouterStatic404(t *testing.T) {

}

func TestRouterUrlPrefixParam(t *testing.T) {

}

func TestRouterGroup(t *testing.T) {
	/*
		routerFactory := NewRouterFactory()
		addGroupRoute := func(group []string, path string) {
			routerFactoryHandler.Group
		}
		routerFactory.addRoute(routerMethod.GET, priority, path, handler)
	*/
	//group操作，路由与中间件
	//http方法指向
	//object路由生成
}
