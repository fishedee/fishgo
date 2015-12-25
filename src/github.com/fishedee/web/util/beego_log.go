package util

import (
	"strconv"
	"strings"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	. "github.com/fishedee/util"
	"fmt"
)

type LogManagerConfig struct {
	Driver     string `json:driver`
	Filename   string `json:"filename"`
	Maxlines   int `json:"maxlines"`
	Maxsize    int `json:"maxsize"`
	Daily      bool  `json:"daily"`
	Maxdays    int `json:"maxdays"`
	Rotate 	   bool `json:"rotate"`
	Level 	   int `json:"level"`
}

type LogManager struct{
	*logs.BeeLogger
	logPrefix string
}

var newLogManagerMemory *MemoryFunc
var newLogManagerFromConfigMemory *MemoryFunc

func init(){
	var err error
	newLogManagerMemory,err = NewMemoryFunc(newLogManager,MemoryFuncCacheNormal)
	if err != nil{
		panic(err)
	}
	newLogManagerFromConfigMemory,err = NewMemoryFunc(newLogManagerFromConfig,MemoryFuncCacheNormal)
	if err != nil{
		panic(err)
	}
}

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

func newLogManager(config LogManagerConfig)(*logs.BeeLogger,error){
	if config.Driver == ""{
		return nil,nil
	}
	logConfigString,err := json.Marshal(config)
	if err != nil{
		panic(err)
	}
	Log := logs.NewLogger(10000)
	err = Log.SetLogger(config.Driver, string(logConfigString))
	if err != nil{
		return nil,err
	}
	beego.BeeLogger = Log
	return Log,nil
}

func NewLogManager(config LogManagerConfig)(*logs.BeeLogger,error){
	result,err := newLogManagerMemory.Call(config)
	if err != nil{
		return nil,err
	}
	return result.(*logs.BeeLogger),err
}

func newLogManagerFromConfig(configName string)(*logs.BeeLogger,error){
	fishlogdriver := beego.AppConfig.String(configName+"driver")
	fishlogfile := beego.AppConfig.String(configName+"file")
	fishlogmaxline := beego.AppConfig.String(configName+"maxline")
	fishlogmaxsize := beego.AppConfig.String(configName+"maxsize")
	fishlogdaily := beego.AppConfig.String(configName+"daily")
	fishlogmaxday := beego.AppConfig.String(configName+"maxday")
	fishlogrotate := beego.AppConfig.String(configName+"rotate")
	fishloglevel := beego.AppConfig.String(configName+"level")

	logConfig := LogManagerConfig{}
	logConfig.Driver = fishlogdriver
	logConfig.Filename = fishlogfile
	logConfig.Maxlines,_ = strconv.Atoi(fishlogmaxline)
	logConfig.Maxsize,_ = strconv.Atoi(fishlogmaxsize)
	logConfig.Daily,_ = strconv.ParseBool(fishlogdaily)
	logConfig.Maxdays,_ = strconv.Atoi(fishlogmaxday)
	logConfig.Rotate,_ = strconv.ParseBool(fishlogrotate)
	logConfig.Level = getLevel(fishloglevel)

	return NewLogManager(logConfig)
}

func NewLogManagerFromConfig(configName string)(*logs.BeeLogger,error){
	result,err := newLogManagerFromConfigMemory.Call(configName)
	if err != nil{
		return nil,err
	}
	return result.(*logs.BeeLogger),err
}

func NewLogManagerWithCtx(ctx *context.Context,logger *logs.BeeLogger)(*LogManager){
	if logger == nil{
		logger = beego.BeeLogger
	}
	logPrefix := ""
	if ctx == nil{
		logPrefix = " 0.0.0.0 * "
	}else{
		ip := ctx.Input.IP()
		url := ctx.Input.Url()
		realIP := ctx.Input.Header("X-Real-Ip")
		if ip == "127.0.0.1" && realIP != ""{
			ip = realIP
		}
		logPrefix = fmt.Sprintf(" %s %s ",ip,url)
	}
	return &LogManager{
		BeeLogger:logger,
		logPrefix:logPrefix,
	}
}

func (this *LogManager) Emergency(format string, v ...interface{}) {
	this.BeeLogger.Emergency(this.logPrefix+format,v...)
}

func (this *LogManager) Alert(format string, v ...interface{}) {
	this.BeeLogger.Alert(this.logPrefix+format,v...)
}

func (this *LogManager) Critical(format string, v ...interface{}) {
	this.BeeLogger.Critical(this.logPrefix+format,v...)
}

func (this *LogManager) Error(format string, v ...interface{}) {
	fmt.Println(this.logPrefix)
	this.BeeLogger.Error(this.logPrefix+format,v...)
}

func (this *LogManager) Warning(format string, v ...interface{}) {
	this.BeeLogger.Warning(this.logPrefix+format,v...)
}

func (this *LogManager) Notice(format string, v ...interface{}) {
	this.BeeLogger.Notice(this.logPrefix+format,v...)
}

func (this *LogManager) Informational(format string, v ...interface{}) {
	this.BeeLogger.Informational(this.logPrefix+format,v...)
}

func (this *LogManager) Debug(format string, v ...interface{}) {
	this.BeeLogger.Debug(this.logPrefix+format,v...)
}

func (this *LogManager) Warn(format string, v ...interface{}) {
	this.BeeLogger.Warn(this.logPrefix+format,v...)
}

func (this *LogManager) Info(format string, v ...interface{}) {
	this.BeeLogger.Info(this.logPrefix+format,v...)
}

func (this *LogManager) Trace(format string, v ...interface{}) {
	this.BeeLogger.Trace(this.logPrefix+format,v...)
}