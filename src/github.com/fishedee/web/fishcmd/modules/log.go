package modules

import (
	"github.com/astaxie/beego/logs"
)

var Log *logs.BeeLogger

func init() {
	Log = logs.NewLogger(102400)
	Log.SetLogger("console", "")
	Log.EnableFuncCallDepth(true)
}
