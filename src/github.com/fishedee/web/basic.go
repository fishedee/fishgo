package web

import (
	"net/http"
	"testing"
)

type Basic struct {
	Ctx      Context
	Config   Config
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
}

var globalBasic Basic

func init() {
	var err error
	globalBasic.Config, err = NewConfig("conf/app.conf")
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
	globalBasic.Log, err = NewLogFromConfig("log")
	if err != nil {
		panic(err)
	}
	globalBasic.Monitor, err = NewMonitorFromConfig("monitor")
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
}
func initBasic(request *http.Request, response http.ResponseWriter, t *testing.T) *Basic {
	result := globalBasic
	result.Ctx = NewContext(request, response, t)
	result.Log = result.Log.WithContextAndMonitor(result.Ctx, result.Monitor)
	if result.Timer != nil {
		result.Timer = result.Timer.WithLog(result.Log)
	}
	if result.Queue != nil {
		result.Queue = result.Queue.WithLog(result.Log)
	}
	if result.Cache != nil {
		result.Cache = result.Cache.WithLog(result.Log)
	}
	return &result
}
