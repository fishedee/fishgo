package util

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	"github.com/k0kubun/pp"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type LogManagerConfig struct {
	Driver      string `json:driver`
	Filename    string `json:"filename"`
	Maxlines    int    `json:"maxlines"`
	Maxsize     int    `json:"maxsize"`
	Daily       bool   `json:"daily"`
	Maxdays     int    `json:"maxdays"`
	Rotate      bool   `json:"rotate"`
	Level       int    `json:"level"`
	PrettyPrint bool   `json:"prettyprint"`
}

type LogManager struct {
	*logs.BeeLogger
	monitor     *MonitorManager
	logPrefix   string
	prettyPrint bool
}

var newLogManagerMemory *MemoryFunc
var newLogManagerFromConfigMemory *MemoryFunc

func init() {
	var err error
	newLogManagerMemory, err = NewMemoryFunc(newLogManager, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
	newLogManagerFromConfigMemory, err = NewMemoryFunc(newLogManagerFromConfig, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
}

func getLevel(in string) int {
	levelString := map[string]int{
		"Emergency":     logs.LevelEmergency,
		"Alert":         logs.LevelAlert,
		"Critical":      logs.LevelCritical,
		"Error":         logs.LevelError,
		"Warning":       logs.LevelWarning,
		"Notice":        logs.LevelNotice,
		"Informational": logs.LevelInformational,
		"Debug":         logs.LevelDebug,
	}
	for key, value := range levelString {
		if strings.ToLower(in) == strings.ToLower(key) {
			return value
		}
	}
	return logs.LevelDebug
}

func newLogManager(config LogManagerConfig) (*LogManager, error) {
	if config.Driver == "" {
		return nil, nil
	}
	logConfigString, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	Log := logs.NewLogger(10000)
	err = Log.SetLogger(config.Driver, string(logConfigString))
	if err != nil {
		return nil, err
	}
	beego.BeeLogger = Log
	return &LogManager{
		BeeLogger:   Log,
		prettyPrint: config.PrettyPrint,
	}, nil
}

func NewLogManager(config LogManagerConfig) (*LogManager, error) {
	result, err := newLogManagerMemory.Call(config)
	if err != nil {
		return nil, err
	}
	return result.(*LogManager), err
}

func newLogManagerFromConfig(configName string) (*LogManager, error) {
	fishlogdriver := globalBasic.Config.String(configName + "driver")
	fishlogfile := globalBasic.Config.String(configName + "file")
	fishlogmaxline := globalBasic.Config.String(configName + "maxline")
	fishlogmaxsize := globalBasic.Config.String(configName + "maxsize")
	fishlogdaily := globalBasic.Config.String(configName + "daily")
	fishlogmaxday := globalBasic.Config.String(configName + "maxday")
	fishlogrotate := globalBasic.Config.String(configName + "rotate")
	fishloglevel := globalBasic.Config.String(configName + "level")
	fishlogprettyprint := globalBasic.Config.String(configName + "prettyprint")

	logConfig := LogManagerConfig{}
	logConfig.Driver = fishlogdriver
	logConfig.Filename = fishlogfile
	logConfig.Maxlines, _ = strconv.Atoi(fishlogmaxline)
	logConfig.Maxsize, _ = strconv.Atoi(fishlogmaxsize)
	logConfig.Daily, _ = strconv.ParseBool(fishlogdaily)
	logConfig.Maxdays, _ = strconv.Atoi(fishlogmaxday)
	logConfig.Rotate, _ = strconv.ParseBool(fishlogrotate)
	logConfig.Level = getLevel(fishloglevel)
	logConfig.PrettyPrint, _ = strconv.ParseBool(fishlogprettyprint)

	return NewLogManager(logConfig)
}

func NewLogManagerFromConfig(configName string) (*LogManager, error) {
	result, err := newLogManagerFromConfigMemory.Call(configName)
	if err != nil {
		return nil, err
	}
	return result.(*LogManager), err
}

func NewLogManagerWithCtxAndMonitor(request *http.Request, monitor *MonitorManager, logger *LogManager) *LogManager {
	var beeLogger *logs.BeeLogger
	if logger != nil {
		beeLogger = logger.BeeLogger
	} else {
		beeLogger = beego.BeeLogger
	}
	logPrefix := ""
	if request == nil {
		logPrefix = " 0.0.0.0 * "
	} else {
		ip := request.RemoteAddr
		url := request.RequestURI
		realIP := request.Header.Get("X-Real-Ip")
		if ip == "127.0.0.1" && realIP != "" {
			ip = realIP
		}
		logPrefix = fmt.Sprintf(" %s %s ", ip, url)
	}
	newLogManager := *logger
	newLogManager.BeeLogger = beeLogger
	newLogManager.logPrefix = logPrefix
	newLogManager.monitor = monitor
	return &newLogManager
}

func (this *LogManager) getTraceLineNumber(traceNumber int) string {
	_, filename, line, ok := runtime.Caller(traceNumber + 1)
	if !ok {
		return "???.go:???"
	} else {
		return path.Base(filename) + ":" + strconv.Itoa(line)
	}
}

func (this *LogManager) getLogFormat(format string, v []interface{}) string {
	if this.prettyPrint {
		format = strings.Replace(format, "%+v", "%v", -1)
		format = strings.Replace(format, "%#v", "%v", -1)
		for singleIndex, singleV := range v {
			singleVType := reflect.TypeOf(singleV)
			singleVTypeKind := GetTypeKind(singleVType)
			if singleVTypeKind == TypeKind.BOOL ||
				singleVTypeKind == TypeKind.INT ||
				singleVTypeKind == TypeKind.UINT ||
				singleVTypeKind == TypeKind.FLOAT ||
				singleVTypeKind == TypeKind.STRING {
				v[singleIndex] = singleV
			} else {
				v[singleIndex] = pp.Sprint(singleV)
			}
		}
	}
	return fmt.Sprintf(this.logPrefix+this.getTraceLineNumber(2)+" "+format, v...)
}

func (this *LogManager) Emergency(format string, v ...interface{}) {
	this.BeeLogger.Emergency(this.getLogFormat(format, v))
}

func (this *LogManager) Alert(format string, v ...interface{}) {
	this.BeeLogger.Alert(this.getLogFormat(format, v))
}

func (this *LogManager) Critical(format string, v ...interface{}) {
	if this.monitor != nil {
		this.monitor.AscCriticalCount()
	}
	this.BeeLogger.Critical(this.getLogFormat(format, v))
}

func (this *LogManager) Error(format string, v ...interface{}) {
	if this.monitor != nil {
		this.monitor.AscErrorCount()
	}
	this.BeeLogger.Error(this.getLogFormat(format, v))
}

func (this *LogManager) Warning(format string, v ...interface{}) {
	this.BeeLogger.Warning(this.getLogFormat(format, v))
}

func (this *LogManager) Notice(format string, v ...interface{}) {
	this.BeeLogger.Notice(this.getLogFormat(format, v))
}

func (this *LogManager) Informational(format string, v ...interface{}) {
	this.BeeLogger.Informational(this.getLogFormat(format, v))
}

func (this *LogManager) Debug(format string, v ...interface{}) {
	this.BeeLogger.Debug(this.getLogFormat(format, v))
}

func (this *LogManager) Warn(format string, v ...interface{}) {
	this.BeeLogger.Warn(this.getLogFormat(format, v))
}

func (this *LogManager) Info(format string, v ...interface{}) {
	this.BeeLogger.Info(this.getLogFormat(format, v))
}

func (this *LogManager) Trace(format string, v ...interface{}) {
	this.BeeLogger.Trace(this.getLogFormat(format, v))
}
