package modules

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

var appName string
var appBinName string
var goPath string
var goPathSrc string
var goPathBin string
var goPathPkg string
var workingDir string
var appAllDirectory []string

func InitConfig() error {
	var err error
	goPath = os.Getenv("GOPATH")
	goPath = strings.TrimRight(goPath, "/")
	if IsExistDir(goPath) == false {
		return fmt.Errorf("$GOPATH [%v] is not exist", goPath)
	}

	goPathSrc = goPath + "/src"
	if IsExistDir(goPathSrc) == false {
		return fmt.Errorf("$GOPATH/src [%v] is not exist", goPathSrc)
	}

	goPathBin = goPath + "/bin"
	if IsExistDir(goPathBin) == false {
		return fmt.Errorf("$GOPATH/bin [%v] is not exist", goPathBin)
	}

	goPathPkg = goPath + "/pkg/" + runtime.GOOS + "_" + runtime.GOARCH

	workingDir, err = os.Getwd()
	if err != nil {
		return err
	}
	workingDir = strings.TrimRight(workingDir, "/")
	appName = ""

	workingDirArray := strings.Split(workingDir, "/")
	appBinName = workingDirArray[len(workingDirArray)-1]
	appAllDirectory, err = getAllAppDirectory(workingDir)
	if err != nil {
		return fmt.Errorf("GetAllAppDirectory Fail [%v]", err)
	}
	return nil
}

func getAllAppDirectory(dir string) ([]string, error) {
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := []string{}
	hasGo := false
	for _, singleFile := range fileInfo {
		if singleFile.IsDir() {
			singleDir := dir + "/" + singleFile.Name()
			singleResult, err := getAllAppDirectory(singleDir)
			if err != nil {
				return nil, err
			}
			result = append(result, singleResult...)
		} else {
			if strings.HasSuffix(singleFile.Name(), ".go") == false {
				continue
			}
			hasGo = true
		}
	}
	if hasGo {
		result = append(result, dir)
	}
	return result, nil
}

func GetAppName() string {
	return appName
}

func GetAppAllDirectory() []string {
	return appAllDirectory
}

func GetAppInstallPath() string {
	return goPathBin + "/" + appBinName
}

func GetAppCurrentPath() string {
	return workingDir + "/" + appBinName
}

func GetWorkginDir() string {
	return workingDir
}

func GetGoPathSrc() string {
	return goPathSrc
}

func GetGoPathPkg() string {
	return goPathPkg
}

func GetGoPathBin() string {
	return goPathBin
}
