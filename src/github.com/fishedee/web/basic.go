package web

import (
	"bytes"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

type Basic struct {
	Ctx      Context
	Config   Configure
	Security Security
	Session  Session
	DB       Database
	DB2      Database
	DB3      Database
	DB4      Database
	DB5      Database
	Log      Log
	Monitor  Monitor
	Timer    Timer
	Queue    Queue
	Cache    Cache
	Grace    Grace
}

var globalBasic Basic

func init() {
	//初始化组件
	var err error
	globalBasic.Config, err = NewConfig("conf/app.conf")
	if err != nil {
		panic(err)
	}
	globalBasic.Log, err = NewLogFromConfig("log")
	if err != nil {
		panic(err)
	}
	globalBasic.Monitor, err = NewMonitorFromConfig("monitor")
	if err != nil {
		panic(err)
	}
	globalBasic.Grace, err = NewGraceFromConfig("grace")
	if err != nil {
		panic(err)
	}
	globalBasic.Security, err = NewSecurityFromConfig("security")
	if err != nil {
		panic(err)
	}
	globalBasic.Session, err = NewSessionFromConfig("session")
	if err != nil {
		panic(err)
	}
	globalBasic.DB, err = NewDatabaseFromConfig("db")
	if err != nil {
		panic(err)
	}
	globalBasic.DB2, err = NewDatabaseFromConfig("db2")
	if err != nil {
		panic(err)
	}
	globalBasic.DB3, err = NewDatabaseFromConfig("db3")
	if err != nil {
		panic(err)
	}
	globalBasic.DB4, err = NewDatabaseFromConfig("db4")
	if err != nil {
		panic(err)
	}
	globalBasic.DB5, err = NewDatabaseFromConfig("db5")
	if err != nil {
		panic(err)
	}
	globalBasic.Timer, err = NewTimer()
	if err != nil {
		panic(err)
	}
	globalBasic.Queue, err = NewQueueFromConfig("queue")
	if err != nil {
		panic(err)
	}
	globalBasic.Cache, err = NewCacheFromConfig("cache")
	if err != nil {
		panic(err)
	}

	//初始化随机数
	rand.Seed(time.Now().Unix())
}

type memoryResponseWriter struct {
	header     http.Header
	headerCode int
	data       []byte
}

func (this *memoryResponseWriter) Header() http.Header {
	if this.header == nil {
		this.header = http.Header{}
	}
	return this.header
}

func (this *memoryResponseWriter) Write(in []byte) (int, error) {
	this.data = append(this.data, in...)
	return len(this.data), nil
}

func (this *memoryResponseWriter) WriteHeader(headerCode int) {
	this.headerCode = headerCode
}

func initEmptyBasic(t *testing.T) *Basic {
	request, err := http.NewRequest("get", "/", bytes.NewReader([]byte("")))
	if err != nil {
		panic(err)
	}
	response := &memoryResponseWriter{}
	return initBasic(request, response, t)
}

func initBasic(request *http.Request, response http.ResponseWriter, t *testing.T) *Basic {
	result := globalBasic
	result.Ctx = NewContext(request, response, t)
	result.Log = result.Log.WithContextAndMonitor(result.Ctx, result.Monitor)
	if result.Session != nil {
		result.Session = result.Session.WithContext(result.Ctx)
	}
	if result.Timer != nil {
		result.Timer = result.Timer.WithLog(result.Log)
	}
	if result.Queue != nil {
		result.Queue = result.Queue.WithLogAndContext(result.Log, result.Ctx)
	}
	if result.Cache != nil {
		result.Cache = result.Cache.WithLog(result.Log)
	}
	return &result
}

func destroyBasic() {
	if globalBasic.Log != nil {
		globalBasic.Log.Close()
	}
	if globalBasic.Timer != nil {
		globalBasic.Timer.Close()
	}
	if globalBasic.Queue != nil {
		globalBasic.Queue.Close()
	}
}

func GetAppBasic() Basic {
	return globalBasic
}
