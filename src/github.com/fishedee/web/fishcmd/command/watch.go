package command

import (
	"github.com/fishedee/web/fishcmd/modules"
)

func Watch(argv []string) (string, error) {
	//读取配置
	err := modules.InitConfig()
	if err != nil {
		return "", err
	}
	appName := modules.GetAppName()

	//初始运行
	err = buildAll(appName)
	if err != nil {
		return "", err
	}
	err = run(appName)
	if err != nil {
		return "", err
	}

	//设置watch的文件
	allFile := modules.GetAppAllDirectory()
	err = modules.Watch(allFile, func(singlePackage string) {
		err := buildAll(appName)
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
