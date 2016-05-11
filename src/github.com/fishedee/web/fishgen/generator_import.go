package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var dirDeclTypeCache map[string]map[string]bool

func getDirDeclTypeInner(dir string) (map[string]bool, error) {
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := map[string]bool{}
	for _, singleFileInfo := range fileInfo {
		if strings.HasSuffix(singleFileInfo.Name(), ".go") == false {
			continue
		}
		parserInfo, err := ParserSingleFile(dir + "/" + singleFileInfo.Name())
		if err != nil {
			return nil, err
		}
		for _, singleDeclType := range parserInfo.declType {
			result[singleDeclType] = true
		}
	}
	return result, nil
}

func getDirDeclType(dir string) (map[string]bool, error) {
	data, ok := dirDeclTypeCache[dir]
	if ok {
		return data, nil
	}

	data, err := getDirDeclTypeInner(dir)
	if err != nil {
		return nil, err
	}

	dirDeclTypeCache[dir] = data
	return data, nil
}

func generateSingleFileImport(data []ParserInfo, source string) (string, error) {
	//解析源代码
	sourceParserInfo, err := ParserSingleSource(source)
	if err != nil {
		return "", err
	}

	//建立导入符号的映射
	nameImport := map[string]ImportInfo{}
	nameImport["InitController"] = ImportInfo{
		name: ".",
		path: "github.com/fishedee/web",
	}
	nameImport["InitModel"] = ImportInfo{
		name: ".",
		path: "github.com/fishedee/web",
	}
	nameImport["InitTest"] = ImportInfo{
		name: ".",
		path: "github.com/fishedee/web",
	}
	nameImport["Exception"] = ImportInfo{
		name: ".",
		path: "github.com/fishedee/language",
	}
	for _, singleParserInfo := range data {
		for _, singleImport := range singleParserInfo.imports {
			if singleImport.name == "_" {
				continue
			}
			if singleImport.name == "." {
				dirImportTypes, err := getDirDeclType(os.Getenv("GOPATH") + "/src/" + singleImport.path)
				if err != nil {
					return "", err
				}
				for singleImportType, _ := range dirImportTypes {
					if singleImportType[0] < 'A' || singleImportType[0] > 'Z' {
						continue
					}
					nameImport[singleImportType] = singleImport
				}
			} else {
				if singleImport.name != "" {
					nameImport[singleImport.name] = singleImport
				} else {
					nameImport[path.Base(singleImport.path)] = singleImport
				}
			}
		}
	}
	dirImportTypes, err := getDirDeclType(data[0].dir)
	if err != nil {
		return "", err
	}
	for singleImportType, _ := range dirImportTypes {
		nameImport[singleImportType] = ImportInfo{}
	}

	result := map[string]ImportInfo{}
	//遍历需要的命名导入
	for _, singleUseType := range sourceParserInfo.useType {
		if singleUseType == "string" ||
			singleUseType == "int" || singleUseType == "int8" || singleUseType == "int16" || singleUseType == "int32" || singleUseType == "int64" ||
			singleUseType == "float" || singleUseType == "float32" || singleUseType == "float64" ||
			singleUseType == "byte" || singleUseType == "bool" ||
			singleUseType == "error" {
			continue
		} else {
			singleNameImport, ok := nameImport[singleUseType]
			if ok {
				result[singleNameImport.path] = singleNameImport
			} else {
				return "", errors.New(fmt.Sprintf("%v can not handle the type %v import", data[0].dir, singleUseType))
			}
		}
	}

	//拼凑导入符号
	resultArray := []string{}
	for _, singleImportInfo := range result {
		if singleImportInfo.path == "" {
			continue
		}
		resultArray = append(
			resultArray,
			singleImportInfo.name+" \""+singleImportInfo.path+"\"",
		)
	}
	return "import (" + strings.Join(resultArray, "\n") + ")\n", nil
}
