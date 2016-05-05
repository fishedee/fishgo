package modules

import (
	"io/ioutil"
	"os"
	"strings"
)

func getPackageFile(packageName string) ([]string, error) {
	dirName := os.Getenv("GOPATH") + "/src/" + packageName
	fileInfo, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, singleFileInfo := range fileInfo {
		fileName := singleFileInfo.Name()
		if strings.HasSuffix(fileName, ".go") == false {
			continue
		}
		result = append(result, dirName+"/"+fileName)
	}
	return result, nil
}
