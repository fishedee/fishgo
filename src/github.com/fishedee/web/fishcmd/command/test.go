package command

import (
	"github.com/fishedee/web/fishcmd/modules"
)

func Test(argv []string) (string, error) {
	//读取参数
	isWatch := false
	appName := "."
	for _, singleArgv := range argv {
		if singleArgv == "--watch" {
			isWatch = true
		} else {
			appName = singleArgv
		}
	}

	//读取配置
	err := modules.InitConfig()
	if err != nil {
		return "", err
	}

	//运行
	handler := func(singlePackage string) {
		err := generate()
		if err != nil {
			return
		}

		err = build(appName)
		if err != nil {
			return
		}

		err = test(appName)
		if err != nil {
			return
		}
	}

	if isWatch {
		allFile := modules.GetAppAllDirectory()
		err = modules.Watch(allFile, handler)
		if err != nil {
			return "", err
		}
		err = modules.RunWatcher()
		if err != nil {
			return "", err
		}
		return "", nil
	} else {
		handler(".")
		return "", nil
	}
}
