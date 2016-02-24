package main

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func generateSingleField(data []FieldInfo) string {
	var result []string
	for _, singleData := range data {
		result = append(result, singleData.name+" "+singleData.tag)
	}
	return strings.Join(result, ",")
}

func generateSingleFieldName(data []FieldInfo) string {
	var result []string
	for _, singleData := range data {
		result = append(result, singleData.name)
	}
	return strings.Join(result, ",")
}

func generateSingleResult(data []FieldInfo) string {
	var result []string
	for index, singleData := range data {
		result = append(result, fmt.Sprintf("_fishgen%d", index+1)+" "+singleData.tag)
	}
	result = append(result, "_fishgenErr error")
	return strings.Join(result, ",")
}

func generateVariable(params int) string {
	var result []string
	for i := 1; i <= params; i++ {
		result = append(result, fmt.Sprintf("_fishgen%d", i))
	}
	return strings.Join(result, ",")
}

func generateSingleReturn(fun FunctionInfo) string {
	resultLen := len(fun.results)
	if resultLen == 0 {
		return fun.receiver[0].name + "." + fun.name + "(" + generateSingleFieldName(fun.params) + ")\n" +
			"return"
	} else {
		return generateVariable(resultLen) + " = " + fun.receiver[0].name + "." + fun.name + "(" + generateSingleFieldName(fun.params) + ")\n" +
			"return"
	}
}

func generateSingleFunction(fun FunctionInfo) string {
	if len(fun.receiver) == 0 {
		return ""
	}
	receiverName := fun.receiver[0].tag
	if strings.HasSuffix(receiverName, "Model") == false {
		return ""
	}
	if fun.name[0] < 'A' || fun.name[0] > 'Z' {
		return ""
	}
	return "func (" + generateSingleField(fun.receiver) + ")" +
		fun.name + "ForError" +
		"(" + generateSingleField(fun.params) + ")" +
		"(" + generateSingleResult(fun.results) + "){\n" +
		"defer Catch(func(exception Exception) {\n" +
		"_fishgenErr = exception\n" +
		"})\n" +
		generateSingleReturn(fun) +
		"}\n"

}

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
		if singleUseType == "string" || singleUseType == "int" ||
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
	result["github.com/fishedee/language"] = ImportInfo{
		name: ".",
		path: "github.com/fishedee/language",
	}
	resultArray := []string{}
	for _, singleImportInfo := range result {
		if singleImportInfo.path == "" {
			continue
		}
		resultArray = append(
			resultArray,
			"import "+singleImportInfo.name+" \""+singleImportInfo.path+"\"",
		)
	}
	return strings.Join(resultArray, "\n") + "\n", nil
}

func generateSingleFileFunction(data []ParserInfo) string {
	var result []string
	for _, singleParserInfo := range data {
		for _, singleFun := range singleParserInfo.functions {
			result = append(result, generateSingleFunction(singleFun))
		}
	}
	return strings.Join(result, "\n")
}

func generateSingleFileContent(data []ParserInfo) (string, error) {
	packageInfo := "package " + data[0].packname + "\n"

	contentInfo := generateSingleFileFunction(data)

	importInfo, err := generateSingleFileImport(data, packageInfo+contentInfo)
	if err != nil {
		return "", err
	}
	return packageInfo + importInfo + contentInfo, nil
}

func generateSingleFileFormat(filename string, data string) (string, error) {
	result, err := format.Source([]byte(data))
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func generateSingleFileWrite(filename string, data string) error {
	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func generateSingleFile(dirname string, data []ParserInfo) error {
	filename := dirname + "/" + GetGenerateFileName(dirname)

	result, err := generateSingleFileContent(data)
	if err != nil {
		return err
	}

	result, err = generateSingleFileFormat(filename, result)
	if err != nil {
		return err
	}

	err = generateSingleFileWrite(filename, result)
	if err != nil {
		return err
	}

	return nil
}

func Generator(data map[string][]os.FileInfo) error {
	for singleKey, singleDir := range data {
		singleResult := []ParserInfo{}
		for _, singleFile := range singleDir {
			singleFileResult, err := ParserSingleFile(singleKey + "/" + singleFile.Name())
			if err != nil {
				return err
			}
			singleResult = append(singleResult, singleFileResult)
		}
		err := generateSingleFile(singleKey, singleResult)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	dirDeclTypeCache = map[string]map[string]bool{}
}
