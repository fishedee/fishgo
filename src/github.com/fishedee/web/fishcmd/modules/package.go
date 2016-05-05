package modules

import (
	"errors"
	"fmt"
	. "github.com/fishedee/language"
	"go/parser"
	"go/token"
	"strings"
)

func getSinglePackageDep(filename string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.New("file parse error " + err.Error())
	}
	result := []string{}
	for _, singleImport := range f.Imports {
		singleImportInfo := strings.Trim(singleImport.Path.Value, "\"")
		result = append(result, singleImportInfo)
	}
	return result, nil
}

func getPackageDep(packages []string) (map[string][]string, error) {
	result := map[string][]string{}
	packagesSet := map[string]bool{}
	for _, singlePackage := range packages {
		packagesSet[singlePackage] = true
	}
	for _, singlePackage := range packages {
		singlePackageFiles, err := getPackageFile(singlePackage)
		if err != nil {
			fmt.Println("uc1")
			return nil, err
		}
		singlePackageImportSet := map[string]bool{}
		for _, singleFile := range singlePackageFiles {
			singleImport, err := getSinglePackageDep(singleFile)
			if err != nil {
				return nil, err
			}
			for _, singleImportPackage := range singleImport {
				if packagesSet[singleImportPackage] == false {
					continue
				}
				singlePackageImportSet[singleImportPackage] = true
			}
		}
		singlePackageImportArray := []string{}
		for singleImport, _ := range singlePackageImportSet {
			singlePackageImportArray = append(singlePackageImportArray, singleImport)
		}
		result[singlePackage] = singlePackageImportArray
	}
	return result, nil
}

func dfs(singlePackage string, packageDep map[string][]string, prevPackage map[string]bool, hasPackage map[string]bool) ([]string, error) {
	var result []string
	prevPackage[singlePackage] = true
	packageDepList := packageDep[singlePackage]
	for _, singlePackageDep := range packageDepList {
		if prevPackage[singlePackageDep] {
			return nil, errors.New(fmt.Sprintf("invalid cycle dep %v,%v", singlePackage, singlePackageDep))
		}
		if hasPackage[singlePackageDep] {
			continue
		}
		singleResult, err := dfs(singlePackageDep, packageDep, prevPackage, hasPackage)
		if err != nil {
			return nil, err
		}
		result = append(result, singleResult...)
	}
	delete(prevPackage, singlePackage)
	hasPackage[singlePackage] = true
	result = append(result, singlePackage)
	return result, nil
}

func PackageList(dir string) (string, []string, error) {
	//获取包
	goListOutput, err := runCmdSync("go", "list", dir)
	if err != nil {
		return "", nil, err
	}

	packages := Explode(string(goListOutput), "\n")
	if len(packages) == 0 {
		return "", nil, errors.New("go list package is empty!")
	}
	mainPackage := packages[0]
	var subPackages []string
	if len(packages) > 1 {
		subPackages = packages[1:]
	}

	//计算包的依赖度
	packageDep, err := getPackageDep(subPackages)
	if err != nil {
		return "", nil, err
	}
	resultSubPackage := []string{}
	hasPackage := map[string]bool{}
	for _, singleSubPackage := range subPackages {
		if hasPackage[singleSubPackage] == true {
			continue
		}
		prevPackage := map[string]bool{}
		singleResult, err := dfs(singleSubPackage, packageDep, prevPackage, hasPackage)
		if err != nil {
			return "", nil, err
		}
		resultSubPackage = append(resultSubPackage, singleResult...)
	}
	return mainPackage, resultSubPackage, nil
}
