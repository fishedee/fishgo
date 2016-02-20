package main

import (
	"os"
	"regexp"
)

func filterSingleDir(path string, data []os.FileInfo) ([]os.FileInfo, error) {
	reg, err := regexp.CompilePOSIX(Config.fileregex)
	if err != nil {
		return nil, err
	}

	dataMap := map[string]os.FileInfo{}
	for _, singleFileInfo := range data {
		dataMap[singleFileInfo.Name()] = singleFileInfo
	}

	//过滤出需要生成的文件
	newdata := []os.FileInfo{}
	for _, singleFileInfo := range data {
		if reg.Match([]byte(path+"/"+singleFileInfo.Name())) == false {
			continue
		}
		generateFileName := GetGenerateFileName(singleFileInfo.Name())
		singleFileInfo2, ok := dataMap[generateFileName]
		if ok && singleFileInfo2.ModTime().After(singleFileInfo.ModTime()) {
			continue
		}
		newdata = append(newdata, singleFileInfo)
	}

	//删除多余的生成文件
	for _, singleFileInfo := range data {
		if IsGenerateFileName(singleFileInfo.Name()) == false {
			continue
		}
		originFileName := GetOriginFileName(singleFileInfo.Name())
		_, ok := dataMap[originFileName]
		if ok {
			err := os.Remove(path + "/" + singleFileInfo.Name())
			if err != nil {
				return nil, err
			}
		}
	}
	return newdata, nil
}

func FilterDir(data map[string][]os.FileInfo) (map[string][]os.FileInfo, error) {
	result := map[string][]os.FileInfo{}
	for singleKey, singleDir := range data {
		singleDir, err := filterSingleDir(singleKey, singleDir)
		if err != nil {
			return nil, err
		}
		if len(singleDir) == 0 {
			continue
		}
		result[singleKey] = singleDir
	}
	return result, nil
}
