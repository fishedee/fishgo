package log

import (
	"encoding/json"
	"fmt"
	"github.com/fishedee/app/log/logs"
	. "github.com/fishedee/language"
	"github.com/k0kubun/pp"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Log interface {
	Emergency(format string, v ...interface{})
	Alert(format string, v ...interface{})
	Critical(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Notice(format string, v ...interface{})
	Informational(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Close()
}

type LogConfig struct {
	Driver      string `json:"driver" config:"driver"`
	Filename    string `json:"filename" config:"file"`
	Maxlines    int    `json:"maxlines" config:"maxline"`
	Maxsize     int    `json:"maxsize" config:"maxsize"`
	Daily       bool   `json:"daily" config:"daily"`
	Maxdays     int    `json:"maxdays" config:"maxday"`
	Rotate      bool   `json:"rotate" config:"rotate"`
	Level       string `json:"level" config:"level"`
	PrettyPrint bool   `json:"prettyprint" config:"prettyprint"`
	Async       bool   `json:"async" config:"async"`
}

type logImplement struct {
	*logs.BeeLogger
	logPrefix   string
	prettyPrint bool
}

func getLevel(in string) int {
	levelString := map[string]int{
		"emergency":     logs.LevelEmergency,
		"alert":         logs.LevelAlert,
		"critical":      logs.LevelCritical,
		"error":         logs.LevelError,
		"warning":       logs.LevelWarning,
		"notice":        logs.LevelNotice,
		"informational": logs.LevelInformational,
		"debug":         logs.LevelDebug,
	}
	value, isExist := levelString[strings.ToLower(in)]
	if isExist {
		return value
	}
	return logs.LevelDebug
}

func NewLog(config LogConfig) (Log, error) {
	//配置处理
	if config.Driver == "" {
		config.Driver = "console"
	}
	if config.Level == "" {
		config.Level = "debug"
	}
	configMap := ArrayToMap(config, "json").(map[string]interface{})
	configMap["level"] = getLevel(config.Level)
	logConfigString, err := json.Marshal(configMap)
	if err != nil {
		panic(err)
	}

	//建立config
	Log := logs.NewLogger(10000)
	err = Log.SetLogger(config.Driver, string(logConfigString))
	if err != nil {
		return nil, err
	}
	if config.Async {
		Log = Log.Async()
	}

	return &logImplement{
		BeeLogger:   Log,
		prettyPrint: config.PrettyPrint,
	}, nil
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
	v = append([]interface{}{this.logPrefix, this.getTraceLineNumber(2)}, v...)
	return fmt.Sprintf("%s %s "+format, v...)
}

func (this *logImplement) Emergency(format string, v ...interface{}) {
	this.BeeLogger.Emergency("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Alert(format string, v ...interface{}) {
	this.BeeLogger.Alert("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Critical(format string, v ...interface{}) {
	this.BeeLogger.Critical("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Error(format string, v ...interface{}) {
	this.BeeLogger.Error("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Warning(format string, v ...interface{}) {
	this.BeeLogger.Warning("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Notice(format string, v ...interface{}) {
	this.BeeLogger.Notice("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Informational(format string, v ...interface{}) {
	this.BeeLogger.Informational("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Debug(format string, v ...interface{}) {
	this.BeeLogger.Debug("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Warn(format string, v ...interface{}) {
	this.BeeLogger.Warn("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Info(format string, v ...interface{}) {
	this.BeeLogger.Info("%s", this.getLogFormat(format, v))
}

func (this *logImplement) Trace(format string, v ...interface{}) {
	this.BeeLogger.Trace(this.getLogFormat(format, v))
}

func (this *logImplement) Close() {
	this.BeeLogger.Close()
}
