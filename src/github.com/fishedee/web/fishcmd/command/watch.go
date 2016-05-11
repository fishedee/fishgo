package command

import (
	"github.com/fishedee/web/fishcmd/modules"
)

func Watch(argv []string) (string, error) {
	//处理参数
	buildAll := getBuildAll(argv)

	//读取配置
	err := modules.InitConfig()
	if err != nil {
		return "", err
	}
	appName := modules.GetAppName()

	//设置watch的文件
	allFile := modules.GetAppAllDirectory()
	err = modules.Watch(allFile, func(singlePackage string) {
		err := modules.GeneratePackage("./...")
		if err != nil {
			return
		}

		err = buildAll(appName)
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
