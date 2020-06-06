package modules

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
	"time"
)

type PackageCodeInfo struct {
	ModifyTime time.Time
	Dependence []string
	Name       string
}

type PackageLibraryInfo struct {
	ModifyTime time.Time
	Symbol     string
}

var (
	cachePackageCodeInfo    = map[string]PackageCodeInfo{}
	cachePackageLibraryInfo = map[string]PackageLibraryInfo{}
)

func getSinglePackageDependence(filename string) ([]string, string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ImportsOnly)
	if err != nil {
		return nil, "", errors.New("file parse error " + err.Error())
	}
	result := []string{}
	for _, singleImport := range f.Imports {
		singleImportInfo := strings.Trim(singleImport.Path.Value, "\"")
		result = append(result, singleImportInfo)
	}
	return result, f.Name.Name, nil
}

func getPackageDirInfo(packageName string) (map[string]time.Time, error) {
	result := map[string]time.Time{}

	goPathSrc := GetGoPathSrc()
	for _, singleGoPath := range goPathSrc {
		fileInfo, err := ioutil.ReadDir(singleGoPath + packageName)
		if err != nil {
			return nil, err
		}

		for _, singleFile := range fileInfo {
			fileName := singleFile.Name()
			if singleFile.IsDir() == true {
				continue
			}
			if strings.HasSuffix(fileName, ".go") == false {
				continue
			}
			if strings.HasSuffix(fileName, "_test.go") == true {
				continue
			}
			result[fileName] = singleFile.ModTime()
		}
	}

	return result, nil
}

func GetPackageCodeInfo(packageName string) (PackageCodeInfo, error) {
	//获取包代码的最新修改时间
	dirInfo, err := getPackageDirInfo(packageName)
	if err != nil {
		return PackageCodeInfo{}, err
	}
	newestModifyTime := time.Time{}
	for _, singleTime := range dirInfo {
		if singleTime.After(newestModifyTime) {
			newestModifyTime = singleTime
		}
	}

	//获取缓存数据
	packageInfo, isExist := cachePackageCodeInfo[packageName]
	if isExist && packageInfo.ModifyTime.Equal(newestModifyTime) {
		return packageInfo, nil
	}

	//获取数据
	dependence := map[string]bool{}
	name := ""
	for singleFile, _ := range dirInfo {
		goPathSrc := GetGoPathSrc()
		for _, singleGoPath := range goPathSrc {
			singleDependence, singleName, err := getSinglePackageDependence(singleGoPath + "/" + packageName + "/" + singleFile)
			if err != nil {
				return PackageCodeInfo{}, err
			}
			for _, singleDe := range singleDependence {
				dependence[singleDe] = true
			}
			if name == "" || name == singleName {
				name = singleName
			} else {
				return PackageCodeInfo{}, fmt.Errorf("invalid %v package has two packageName %v,%v", packageName, name, singleName)
			}
		}
	}
	resultDep := []string{}
	for singleDe, _ := range dependence {
		resultDep = append(resultDep, singleDe)
	}
	result := PackageCodeInfo{
		ModifyTime: newestModifyTime,
		Dependence: resultDep,
		Name:       name,
	}

	//写入缓存
	cachePackageCodeInfo[packageName] = result
	return result, nil
}
