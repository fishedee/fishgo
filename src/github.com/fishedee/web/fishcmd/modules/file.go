package modules

import (
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func IsExistFile(filename string) bool {
	var exist bool
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		exist = true
	}
	return exist
}

func IsExistDir(dirName string) bool {
	_, err := ioutil.ReadDir(dirName)
	return err == nil
}

func GetFileModifyTime(fileName string) (time.Time, error) {
	fileInfo, err := os.Lstat(fileName)
	if err != nil {
		return time.Time{}, nil
	}
	return fileInfo.ModTime(), nil
}

func ReadDir(path string) (map[string][]os.FileInfo, error) {
	result := map[string][]os.FileInfo{}
	tempResult, err := ioutil.ReadDir(path)
	if err != nil {
		Log.Error("ReadDir fail! error: %v", err.Error())
		return nil, err
	}
	singleResult := []os.FileInfo{}
	for _, singleFileInfo := range tempResult {
		name := path + "/" + singleFileInfo.Name()
		if singleFileInfo.IsDir() {
			//发现dir
			result2, err := ReadDir(name)
			if err != nil {
				return nil, err
			}
			result = combineDirInfo(result, result2)
		} else {
			//发现源代码文件
			if strings.HasSuffix(name, ".go") == false {
				continue
			}
			singleResult = append(singleResult, singleFileInfo)
		}
	}
	if len(singleResult) != 0 {
		result[path] = singleResult
	}
	return result, nil
}

func combineDirInfo(a1 map[string][]os.FileInfo, a2 map[string][]os.FileInfo) map[string][]os.FileInfo {
	for key, value := range a2 {
		if len(value) != 0 {
			a1[key] = value
		}
	}
	return a1
}
