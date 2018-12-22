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

	//搜索出生成文件
	genfile := GetGenerateFileName(path)

	//过滤出需要生成的文件
	newdata := []os.FileInfo{}
	var genfiledata os.FileInfo
	for _, singleFileInfo := range data {
		if singleFileInfo.Name() == genfile {
			genfiledata = singleFileInfo
			continue
		}
		if reg.Match([]byte(path+"/"+singleFileInfo.Name())) == false {
			continue
		}
		newdata = append(newdata, singleFileInfo)
	}

	//判断是否需要生成文件
	for _, singleFileInfo := range newdata {
		if genfiledata == nil || singleFileInfo.ModTime().After(genfiledata.ModTime()) {
			return newdata, nil
		}
	}

	//FIXME 暂时不做增量更新策略
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
