package web

import (
	_ "github.com/a"
	"github.com/astaxie/beego"
)

func init() {
	accessLogs, err := beego.AppConfig.Bool("AccessLogs")
	if err == nil {
		beego.BConfig.Log.AccessLogs = accessLogs
	}
	beego.BConfig.CopyRequestBody = true
}
