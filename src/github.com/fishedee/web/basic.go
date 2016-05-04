package web

import (
	"github.com/a"
	"github.com/astaxie/beego"
	. "github.com/fishedee/web/util"
	"net/http"
	"os"
	"path"
	"testing"
)

type Basic struct {
	Ctx      *Context
	Security *Security
	Session  *Session
	DB       *Database
	DB2      *Database
	DB3      *Database
	DB4      *Database
	DB5      *Database
	logger   *Log
	Log      *Log
	Monitor  *Monitor
	timer    *Timer
	Timer    *Timer
	queue    *Queue
	Queue    *Queue
	cache    *Cache
	Cache    *Cache
}

var globalBasic Basic

func checkFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

func findAppConfPath() string {
	workingDir := a.GetWorkingDir()
	if checkFileExist(workingDir + "/conf/app.conf") {
		return ""
	}

	for workingDir != "/" {
		workingDir = path.Dir(workingDir)
		appPath := workingDir + "/conf/app.conf"
		if checkFileExist(appPath) {
			return workingDir
		}
	}
	panic("could not found app.conf")
}

func init() {
	//确定appPath
	appPath := findAppConfPath()
	if appPath != "" {
		os.Setenv("BEEGO_RUNMODE", "test")
		err := beego.LoadAppConfig("ini", appPath+"/conf/app.conf")
		if err != nil {
			panic(err)
		}
		beego.TestBeegoInit(appPath)
	}

	var err error
	globalBasic.Security, err = NewSecurityManagerFromConfig("fishsecurity")
	if err != nil {
		panic(err)
	}
	globalBasic.Session, err = NewSessionManagerFromConfig("fishsession")
	if err != nil {
		panic(err)
	}
	globalBasic.DB, err = NewDatabaseManagerFromConfig("fishdb")
	if err != nil {
		panic(err)
	}
	globalBasic.DB2, err = NewDatabaseManagerFromConfig("fishdb2")
	if err != nil {
		panic(err)
	}
	globalBasic.DB3, err = NewDatabaseManagerFromConfig("fishdb3")
	if err != nil {
		panic(err)
	}
	globalBasic.DB4, err = NewDatabaseManagerFromConfig("fishdb4")
	if err != nil {
		panic(err)
	}
	globalBasic.DB5, err = NewDatabaseManagerFromConfig("fishdb5")
	if err != nil {
		panic(err)
	}
	globalBasic.logger, err = NewLogManagerFromConfig("fishlog")
	if err != nil {
		panic(err)
	}
	globalBasic.Monitor, err = NewMonitorManagerFromConfig("fishmonitor")
	if err != nil {
		panic(err)
	}
	globalBasic.timer, err = NewTimerManager()
	if err != nil {
		panic(err)
	}
	globalBasic.queue, err = NewQueueManagerFromConfig("fishqueue")
	if err != nil {
		panic(err)
	}
	globalBasic.cache, err = NewCacheManagerFromConfig("fishcache")
	if err != nil {
		panic(err)
	}
}
func initBasic(ctx *context.Context, t *testing.T) *Basic {
	result := globalBasic
	result.ctx = ctx
	result.t = t
	result.Log = NewLogManagerWithCtxAndMonitor(ctx, result.Monitor, result.logger)
	result.Timer = NewTimerManagerWithLog(result.Log, result.timer)
	result.Queue = NewQueueManagerWithLog(result.Log, result.queue)
	result.Cache = NewCacheManagerWithLog(result.Log, result.cache)
	return &result
}
