package web

import (
	"github.com/astaxie/beego"
)

func init() {
	accessLogs, err := beego.AppConfig.Bool("AccessLogs")
	if err == nil {
		beego.BConfig.Log.AccessLogs = accessLogs
	}
}
