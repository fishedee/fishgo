package web

import (
	"github.com/astaxie/beego"
	"fmt"
)

func init(){
	accessLogs, err := beego.AppConfig.Bool("AccessLogs");
	if err == nil {
		beego.AccessLogs = accessLogs
	}
}
