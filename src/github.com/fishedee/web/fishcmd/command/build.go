package command

import (
	"github.com/fishedee/web/fishcmd/modules"
)

func Build(argv []string) (string, error) {
	//读取配置
	err := modules.InitConfig()
	if err != nil {
		return "", err
	}
	appName := modules.GetAppName()

	//安装
	err = buildAll(appName)
	if err != nil {
		return "", err
	}
	return "", nil
}
