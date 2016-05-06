package command

import (
	"github.com/fishedee/web/fishcmd/modules"
)

func buildAll(appName string) error {
	//安装文件
	var err error
	timer := modules.NewTimer()
	modules.Log.Debug("start building (" + appName + ")...")
	timer.Start()
	err = modules.InstallPackage(appName)
	if err != nil {
		modules.Log.Error("build fail! error: %v", err.Error())
		return err
	}

	//复制文件
	err = modules.CopyFile(modules.GetAppInstallPath(), modules.GetAppCurrentPath())
	if err != nil {
		modules.Log.Error("copy fail! error: %v", err.Error())
		return err
	}
	timer.Stop()
	modules.Log.Debug("build success! time: %v", timer.Elapsed())
	return nil
}

func run(appName string) error {
	err := modules.RunPackage(appName)
	if err != nil {
		modules.Log.Error("%v running fail! error: %v", appName, err.Error())
		return err
	}
	modules.Log.Debug("%v is running", appName)
	return nil
}
