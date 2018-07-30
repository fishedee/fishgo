package command

import (
	"github.com/fishedee/web/fishcmd/modules"
	"strings"
)

func build(appName string) error {
	//安装文件
	var err error
	timer := modules.NewTimer()
	modules.Log.Debug("start building (" + appName + ")...")
	timer.Start()
	err = modules.BuildPackage(appName)
	if err != nil {
		modules.Log.Error("build fail! error: %v", err.Error())
		return err
	}
	timer.Stop()
	modules.Log.Debug("build success! time: %v", timer.Elapsed())
	return nil
}

func test(appName string, args string) error {
	timer := modules.NewTimer()
	modules.Log.Debug("start test ...")
	timer.Start()
	err := modules.TestPackage(appName, args)
	if err != nil {
		modules.Log.Error("test fail! error: %v", err.Error())
		return err
	}
	timer.Stop()
	modules.Log.Debug("test success! time: %v", timer.Elapsed())
	return nil
}

func generate() error {
	timer := modules.NewTimer()
	modules.Log.Debug("start generate ...")
	timer.Start()
	err := modules.GeneratePackage("./...")
	if err != nil {
		modules.Log.Error("generate fail! error: %v", err.Error())
		return err
	}
	timer.Stop()
	modules.Log.Debug("generate success! time: %v", timer.Elapsed())
	return nil
}

func run(appName string, isAsync bool) error {
	appNameArray := strings.Split(appName, "/")
	appBinName := appNameArray[len(appNameArray)-1]

	err := modules.RunPackage(appBinName, isAsync)
	if err != nil {
		modules.Log.Error("%v running fail! error: %v", appBinName, err.Error())
		return err
	}
	modules.Log.Debug("%v is running", appBinName)
	return nil
}
