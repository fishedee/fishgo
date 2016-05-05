package command

import (
	"github.com/fishedee/web/fishcmd/modules"
)

func Watch(argv []string) (string, error) {
	appName := "server"
	if len(argv) != 0 {
		appName = argv[0]
	}

	//初始运行
	mainPackage, subPackageInfo, err := buildAll(appName)
	if err != nil {
		return "", err
	}
	err = run(appName)
	if err != nil {
		return "", err
	}

	//设置watch的文件
	err = modules.Watch([]string{mainPackage}, func(singlePackage string) {
		err := buildSubPackage(singlePackage)
		if err != nil {
			return
		}

		err = run(appName)
		if err != nil {
			return
		}
	})
	if err != nil {
		return "", err
	}
	err = modules.Watch(subPackageInfo, func(singlePackage string) {
		err := buildMainPackage(singlePackage, appName)
		if err != nil {
			return
		}

		err = run(appName)
		if err != nil {
			return
		}
	})
	if err != nil {
		return "", err
	}

	//开始watch
	err = modules.RunWatcher()
	if err != nil {
		return "", err
	}
	return "", nil
}
