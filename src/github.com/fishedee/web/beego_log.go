package web

import (
	"strconv"
	"strings"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var Log *logs.BeeLogger

func getLevel(in string)(int){
	levelString := map[string]int{
		"Emergency":logs.LevelEmergency,
		"Alert":logs.LevelAlert,
		"Critical":logs.LevelCritical,
		"Error":logs.LevelError,
		"Warning":logs.LevelWarning,
		"Notice":logs.LevelNotice,
		"Informational":logs.LevelInformational,
		"Debug":logs.LevelDebug,
	}
	for key,value := range levelString{
		if strings.ToLower(in) == strings.ToLower(key){
			return value
		}
	}
	return logs.LevelDebug
}

func init(){
	fishlogfile := beego.AppConfig.String("fishlogfile")
	fishlogmaxline := beego.AppConfig.String("fishlogmaxline")
	fishlogmaxsize := beego.AppConfig.String("fishlogmaxsize")
	fishlogdaily := beego.AppConfig.String("fishlogdaily")
	fishlogmaxday := beego.AppConfig.String("fishlogmaxday")
	fishlogrotate := beego.AppConfig.String("fishlogrotate")
	fishloglevel := beego.AppConfig.String("fishloglevel")

	if fishlogfile == ""{
		return
	}

	var logConfig struct{
		Filename string `json:"filename"`
		Maxlines          int `json:"maxlines"`
		Maxsize         int `json:"maxsize"`
		Daily          bool  `json:"daily"`
		Maxdays        int `json:"maxdays"`
		Rotate bool `json:"rotate"`
		Level int `json:"level"`
	}
	logConfig.Filename = fishlogfile
	logConfig.Maxlines,_ = strconv.Atoi(fishlogmaxline)
	logConfig.Maxsize,_ = strconv.Atoi(fishlogmaxsize)
	logConfig.Daily,_ = strconv.ParseBool(fishlogdaily)
	logConfig.Maxdays,_ = strconv.Atoi(fishlogmaxday)
	logConfig.Rotate,_ = strconv.ParseBool(fishlogrotate)
	logConfig.Level = getLevel(fishloglevel)

	logConfigString,err := json.Marshal(logConfig)
	if err != nil{
		panic(err)
	}

	Log = logs.NewLogger(10000)
	err = Log.SetLogger("file", string(logConfigString))
	if err != nil{
		panic(err)
	}
}