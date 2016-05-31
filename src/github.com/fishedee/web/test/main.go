package main

import (
	. "github.com/fishedee/encoding"
	"github.com/fishedee/web"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type User struct {
	UserId     int `xorm:"autoincr"`
	Name       string
	Password   string
	Type       int
	Email      string
	CreateTime time.Time `xorm:"created"`
	ModifyTime time.Time `xorm:"updated"`
}

type Users struct {
	Count int
	Data  []User
}
type testController struct {
	web.Controller
}

func (this *testController) Doing_Json() interface{} {
	return "Hello World"
}

func (this *testController) DbTask_Json() interface{} {
	result := Users{}

	db := this.DB.NewSession()
	defer db.Close()

	count, err := db.Count(&User{})
	if err != nil {
		panic(err)
	}
	result.Count = int(count)

	data := []User{}
	err = db.Find(&data)
	if err != nil {
		panic(err)
	}
	result.Data = data

	return result
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
	var result struct {
		Code int
		Msg  string
		Data interface{}
	}
	result.Data = data
	dataByte, err := EncodeJson(result)
	if err != nil {
		panic(err)
	}
	this.Ctx.Write(dataByte)
}

//go:generate fishgen ^./models/.*(ao|db)\.go$
func main() {
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			panic(err)
		}
	}()
	web.InitDaemon(func(this *testController) {
		this.Queue.Consume("/test", (*testController).receiveEvent)
	})
	web.InitRoute("", &testController{})
	web.Run()
}
