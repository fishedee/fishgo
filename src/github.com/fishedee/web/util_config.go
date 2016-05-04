package web

import (
	"errors"
	"github.com/astaxie/beego/config"
	"os"
	"path"
	"strconv"
)

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
	appPath := workingDir + file
	if checkFileExist(appPath) {
		return appPath, true, nil
	}

	for workingDir != "/" {
		workingDir = path.Dir(workingDir)
		appPath := workingDir + file
		if checkFileExist(appPath) {
			return appPath, false, nil
		}
	}
	return "", false, errors.New("can't not find conf")
}

func NewConfigManager(file string) (*ConfigManager, error) {
	appConfigPath, isCurrentDir, err := findAppConfPath(file)
	if err != nil {
		return err
	}
	configer, err := config.NewConfig("ini", appConfigPath)
	if err != nil {
		return err
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

	return &ConfigManager{
		runMode:  runMode,
		configer: configer,
	}, nil
}

type ConfigManager struct {
	runMode  string
	configer config.Configer
}

func (this *ConfigManager) Get(key string) string {
	if v := configer.String(this.runMode + "::" + key); v != "" {
		return v
	}
	return configer.String(key)
}

func (this *ConfigManager) GetFloat(key string) float64 {
	v := this.Get(key)
	vF, _ := strconv.ParseFloat(v, 64)
	return vF
}

func (this *ConfigManager) GetInt(key string) int {
	v := this.Get(key)
	vI, _ := strconv.ParseInt(v, 10, 64)
	return int(vI)
}

func (this *ConfigManager) GetBool(key string) bool {
	v := this.Get(key)
	vB, _ := strconv.ParseBool(v)
	return bool(vB)
}
