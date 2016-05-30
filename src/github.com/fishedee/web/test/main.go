package main

import (
	"github.com/fishedee/web"
	"net/http"
	"time"
)

type testController struct {
	web.Controller
}

func (this *testController) Doing_Json() interface{} {
	return "Hello World"
}

func (this *testController) LongTask_Json() interface{} {
	this.Log.Debug("task begin: %v", time.Now())
	time.Sleep(time.Second * 5)
	this.Log.Debug("task end: %v", time.Now())
	return ""
}

func (this *testController) receiveEvent(data int) {
	request := this.Ctx.GetRawRequest().(*http.Request)
	this.Log.Debug("%v", request)
	this.Log.Debug("%v", data)
}

func (this *testController) AutoRender(data interface{}, viewname string) {
	this.Ctx.Write([]byte(data.(string)))
}

//go:generate fishgen ^./models/.*(ao|db)\.go$
func main() {
	web.InitDaemon(func(this *testController) {
		this.Queue.Consume("/test", (*testController).receiveEvent)
	})
	web.InitRoute("", &testController{})
	web.Run()
}
