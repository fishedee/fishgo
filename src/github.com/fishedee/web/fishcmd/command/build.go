package command

import (
	"github.com/fishedee/web/fishcmd/modules"
)

func Build(argv []string) (string, error) {
	appName := "server"
	if len(argv) != 0 {
		appName = argv[0]
	}

	_, _, err := buildAll(appName)
	if err != nil {
		return "", err
	}
	return "", nil
}

func buildMainPackage(singlePackage string, appName string) error {
	timer := modules.NewTimer()
	modules.Log.Debug("start building (" + singlePackage + ")")
	timer.Start()
	err := modules.BuildPackage(singlePackage, appName)
	if err != nil {
		modules.Log.Error("build fail! error: %v", err.Error())
		return err
	}
	timer.Stop()
	modules.Log.Debug("build success! time: %v", timer.Elapsed())
	return nil
}

func buildSubPackage(singlePackage string) error {
	timer := modules.NewTimer()
	modules.Log.Debug("start building (" + singlePackage + ")...")
	timer.Start()
	err := modules.InstallPackage(singlePackage)
	if err != nil {
		modules.Log.Error("build fail! error: %v", err.Error())
		return err
	}
	timer.Stop()
	modules.Log.Debug("build success! time: %v", timer.Elapsed())
	return nil
}

func buildAll(appName string) (string, []string, error) {
	timer := modules.NewTimer()
	modules.Log.Debug("start building!")
	timer.Start()
	modules.Log.Debug("list package...")
	mainPackage, subPackageInfo, err := modules.PackageList("./...")
	if err != nil {
		modules.Log.Error("build fail! error: %v", err.Error())
		return "", nil, err
	}
	for _, singleSubPageInfo := range subPackageInfo {
		modules.Log.Debug("install package (" + singleSubPageInfo + ")...")
		err := modules.InstallPackage(singleSubPageInfo)
		if err != nil {
			modules.Log.Error("build fail! error: %v", err.Error())
			return "", nil, err
		}
	}
	modules.Log.Debug("build package (" + mainPackage + ")...")
	err = modules.BuildPackage(mainPackage, appName)
	if err != nil {
		modules.Log.Error("build fail! error: %v", err.Error())
		return "", nil, err
	}

	timer.Stop()
	modules.Log.Debug("build success! time: %v", timer.Elapsed())
	return mainPackage, subPackageInfo, nil
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
