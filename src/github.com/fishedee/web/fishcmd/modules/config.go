package modules

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

var goPath []string
var goPathSrc []string
var goPathBin []string
var goPathPkg []string
var workingDir string
var appAllDirectory []string

func InitConfig() error {
	var err error
	osGoPath := strings.Split(os.Getenv("GOPATH"), ":")
	for _, singleGoPath := range osGoPath {
		singleGoPath = strings.TrimRight(singleGoPath, "/")
		if singleGoPath == "" {
			continue
		}

		goPath = append(goPath, singleGoPath)
		if IsExistDir(singleGoPath) == false {
			return fmt.Errorf("$GOPATH [%v] is not exist", singleGoPath)
		}

		src := singleGoPath + "/src"
		goPathSrc = append(goPathSrc, src)
		if IsExistDir(src) == false {
			return fmt.Errorf("$GOPATH/src [%v] is not exist", src)
		}

		bin := singleGoPath + "/bin"
		goPathBin = append(goPathBin, bin)
		if IsExistDir(bin) == false {
			return fmt.Errorf("$GOPATH/bin [%v] is not exist", bin)
		}

		goPathPkg = append(goPathPkg,singleGoPath + "/pkg/" + runtime.GOOS + "_" + runtime.GOARCH)

	}


	workingDir, err = os.Getwd()
	if err != nil {
		return err
	}
	workingDir = strings.TrimRight(workingDir, "/")

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

func GetAppAllDirectory() []string {
	return appAllDirectory
}

func GetGoPathSrc() []string {
	return goPathSrc
}

