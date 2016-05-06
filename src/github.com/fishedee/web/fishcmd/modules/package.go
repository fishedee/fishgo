package modules

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"sort"
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
	fileInfo, err := ioutil.ReadDir(GetGoPathSrc() + packageName)
	if err != nil {
		return nil, err
	}

	result := map[string]time.Time{}
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
		singleDependence, singleName, err := getSinglePackageDependence(GetGoPathSrc() + "/" + packageName + "/" + singleFile)
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

func getSinglePackageSymbol(fileName string) (string, error) {
	//读取.a文件的导出符号
	fileReader, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	fileBufReader := bufio.NewReader(fileReader)
	result := []string{}
	hasBegin := false
	for {
		singleLine, err := fileBufReader.ReadString('\n')
		if err != nil {
			return "", err
		}
		if singleLine == "$$" {
			if hasBegin == false {
				hasBegin = true
			} else {
				break
			}
		} else {
			if hasBegin {
				result = append(result, singleLine)
			}
		}
	}

	//排序并重建为md5
	resultSlice := sort.StringSlice(result)
	sort.Sort(resultSlice)
	dataBuffer := bytes.Buffer{}
	for _, singleResult := range result {
		dataBuffer.WriteString(singleResult)
	}

	hash := md5.New()
	hash.Write(dataBuffer.Bytes())
	md5Byte := hash.Sum(nil)
	return hex.EncodeToString(md5Byte), nil
}

func GetPackageLibraryInfo(packageName string) (PackageLibraryInfo, error) {
	//获取库代码的最新修改时间
	packageLibraryName := GetGoPathPkg() + "/" + packageName + ".a"
	modifyTime, err := GetFileModifyTime(packageLibraryName)
	if err != nil {
		return PackageLibraryInfo{}, err
	}

	//获取缓存数据
	libraryInfo, isExist := cachePackageLibraryInfo[packageName]
	if isExist && libraryInfo.ModifyTime == modifyTime {
		return libraryInfo, nil
	}
	symbol, err := getSinglePackageSymbol(packageLibraryName)
	if err != nil {
		return PackageLibraryInfo{}, err
	}
	result := PackageLibraryInfo{
		ModifyTime: modifyTime,
		Symbol:     symbol,
	}

	//写入缓存
	cachePackageLibraryInfo[packageName] = result
	return result, nil
}

func RefreshPackageLibrary(packageName string) error {
	//更新.a的时间，防止继续编译
	packageLibraryName := GetGoPathPkg() + "/" + packageName + ".a"
	now := time.Now()
	err := os.Chtimes(packageLibraryName, now, now)
	if err != nil {
		return err
	}

	//更新缓存
	curLibrary := cachePackageLibraryInfo[packageName]
	curLibrary.ModifyTime = now
	cachePackageLibraryInfo[packageName] = curLibrary
	return nil
}
