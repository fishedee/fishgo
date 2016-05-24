package main

import (
	"github.com/fishedee/web"
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

func (this *testController) AutoRender(data interface{}, viewname string) {
	this.Ctx.Write([]byte(data.(string)))
}

//go:generate fishgen ^./models/.*(ao|db)\.go$
func main() {
	web.InitRoute("", &testController{})
	web.Run()
}
