package web

import (
	"errors"
	"github.com/astaxie/beego/config"
	"os"
	"path"
	"strconv"
	"strings"
)

type Configure interface {
	GetString(key string) string
	GetFloat(key string) float64
	GetInt(key string) int
	GetBool(key string) bool
}

type configureImplement struct {
	runMode  string
	configer config.Configer
}

func checkFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

func findAppConfPath(file string) (string, bool, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", false, err
	}
	appPath := workingDir + "/" + file
	if checkFileExist(appPath) {
		return appPath, true, nil
	}

	for workingDir != "/" {
		workingDir = path.Dir(workingDir)
		appPath := workingDir + "/" + file
		if checkFileExist(appPath) {
			return appPath, false, nil
		}
	}
	return "", false, errors.New("can't not find conf")
}

func NewConfig(file string) (Configure, error) {
	appConfigPath, isCurrentDir, err := findAppConfPath(file)
	if err != nil {
		return nil, err
	}
	configer, err := config.NewConfig("ini", appConfigPath)
	if err != nil {
		return nil, err
	}

	var runMode string
	if isCurrentDir == false {
		runMode = "test"
	} else if envRunMode := os.Getenv("BEEGO_RUNMODE"); envRunMode != "" {
		runMode = envRunMode
	} else if configRunMode := configer.String("RunMode"); configRunMode != "" {
		runMode = configRunMode
	} else {
		runMode = "dev"
	}

	return &configureImplement{
		runMode:  runMode,
		configer: configer,
	}, nil
}

func (this *configureImplement) GetString(key string) string {
	if strings.ToLower(key) == "runmode" {
		return this.runMode
	}
	if v := this.configer.String(this.runMode + "::" + key); v != "" {
		return v
	}
	return this.configer.String(key)
}

func (this *configureImplement) GetFloat(key string) float64 {
	v := this.GetString(key)
	vF, _ := strconv.ParseFloat(v, 64)
	return vF
}

func (this *configureImplement) GetInt(key string) int {
	v := this.GetString(key)
	vI, _ := strconv.ParseInt(v, 10, 64)
	return int(vI)
}

func (this *configureImplement) GetBool(key string) bool {
	v := this.GetString(key)
	vB, _ := strconv.ParseBool(v)
	return bool(vB)
}
