package main

import (
	"github.com/fishedee/web"
)

type testController struct {
	web.Controller
}

func (this *testController) Doing_Json() interface{} {
	return "Hello World"
}

func (this *testController) AutoRender(data interface{}, viewname string) {
	this.Ctx.Write([]byte(data.(string)))
}

//go:generate fishgen ^./models/.*(ao|db)\.go$
func main() {
	web.InitRoute("", &testController{})
	web.Run()
}
