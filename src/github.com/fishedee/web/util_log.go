package web

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	. "github.com/fishedee/language"
	"github.com/k0kubun/pp"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Log interface {
	WithContextAndMonitor(ctx Context, monitor Monitor) Log
	Emergency(format string, v ...interface{})
	Alert(format string, v ...interface{})
	Critical(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Notice(format string, v ...interface{})
	Informational(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type LogConfig struct {
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

type logImplement struct {
	*logs.BeeLogger
	monitor     Monitor
	logPrefix   string
	prettyPrint bool
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

func NewLog(config LogConfig) (Log, error) {
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
	return &logImplement{
		BeeLogger:   Log,
		prettyPrint: config.PrettyPrint,
	}, nil
}

func NewLogFromConfig(configName string) (Log, error) {
	fishlogdriver := globalBasic.Config.GetString(configName + "driver")
	fishlogfile := globalBasic.Config.GetString(configName + "file")
	fishlogmaxline := globalBasic.Config.GetString(configName + "maxline")
	fishlogmaxsize := globalBasic.Config.GetString(configName + "maxsize")
	fishlogdaily := globalBasic.Config.GetString(configName + "daily")
	fishlogmaxday := globalBasic.Config.GetString(configName + "maxday")
	fishlogrotate := globalBasic.Config.GetString(configName + "rotate")
	fishloglevel := globalBasic.Config.GetString(configName + "level")
	fishlogprettyprint := globalBasic.Config.GetString(configName + "prettyprint")

	logConfig := LogConfig{}
	logConfig.Driver = fishlogdriver
	logConfig.Filename = fishlogfile
	logConfig.Maxlines, _ = strconv.Atoi(fishlogmaxline)
	logConfig.Maxsize, _ = strconv.Atoi(fishlogmaxsize)
	logConfig.Daily, _ = strconv.ParseBool(fishlogdaily)
	logConfig.Maxdays, _ = strconv.Atoi(fishlogmaxday)
	logConfig.Rotate, _ = strconv.ParseBool(fishlogrotate)
	logConfig.Level = getLevel(fishloglevel)
	logConfig.PrettyPrint, _ = strconv.ParseBool(fishlogprettyprint)

	return NewLog(logConfig)
}

func (this *logImplement) WithContextAndMonitor(ctx Context, monitor Monitor) Log {
	logPrefix := ctx.GetRemoteAddr()
	newLogManager := *this
	newLogManager.logPrefix = logPrefix
	newLogManager.monitor = monitor
	return &newLogManager
}

func (this *logImplement) getTraceLineNumber(traceNumber int) string {
	_, filename, line, ok := runtime.Caller(traceNumber + 1)
	if !ok {
		return "???.go:???"
	} else {
		return path.Base(filename) + ":" + strconv.Itoa(line)
	}
}

func (this *logImplement) getLogFormat(format string, v []interface{}) string {
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
	return fmt.Sprintf(this.logPrefix+" "+this.getTraceLineNumber(2)+" "+format, v...)
}

func (this *logImplement) Emergency(format string, v ...interface{}) {
	this.BeeLogger.Emergency(this.getLogFormat(format, v))
}

func (this *logImplement) Alert(format string, v ...interface{}) {
	this.BeeLogger.Alert(this.getLogFormat(format, v))
}

func (this *logImplement) Critical(format string, v ...interface{}) {
	if this.monitor != nil {
		this.monitor.AscCriticalCount()
	}
	this.BeeLogger.Critical(this.getLogFormat(format, v))
}

func (this *logImplement) Error(format string, v ...interface{}) {
	if this.monitor != nil {
		this.monitor.AscErrorCount()
	}
	this.BeeLogger.Error(this.getLogFormat(format, v))
}

func (this *logImplement) Warning(format string, v ...interface{}) {
	this.BeeLogger.Warning(this.getLogFormat(format, v))
}

func (this *logImplement) Notice(format string, v ...interface{}) {
	this.BeeLogger.Notice(this.getLogFormat(format, v))
}

func (this *logImplement) Informational(format string, v ...interface{}) {
	this.BeeLogger.Informational(this.getLogFormat(format, v))
}

func (this *logImplement) Debug(format string, v ...interface{}) {
	this.BeeLogger.Debug(this.getLogFormat(format, v))
}

func (this *logImplement) Warn(format string, v ...interface{}) {
	this.BeeLogger.Warn(this.getLogFormat(format, v))
}

func (this *logImplement) Info(format string, v ...interface{}) {
	this.BeeLogger.Info(this.getLogFormat(format, v))
}

func (this *logImplement) Trace(format string, v ...interface{}) {
	this.BeeLogger.Trace(this.getLogFormat(format, v))
}
