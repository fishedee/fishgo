package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func combineDirInfo(a1 map[string][]os.FileInfo, a2 map[string][]os.FileInfo) map[string][]os.FileInfo {
	for key, value := range a2 {
		if len(value) != 0 {
			a1[key] = value
		}
	}
	return a1
}

func ReadDir(path string) (map[string][]os.FileInfo, error) {
	result := map[string][]os.FileInfo{}
	tempResult, err := ioutil.ReadDir(path)
	if err != nil {
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
