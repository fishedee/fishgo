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

		// 支持传递不定参数
		uncertainParamStr := ""
		if strings.Contains(singleData.tag, "...") {
			uncertainParamStr = "..."
		}

		result = append(result, singleData.name+uncertainParamStr)
	}
	return strings.Join(result, ",")
}

func generateSingleResult(data []FieldInfo) string {
	var result []string
	for index, singleData := range data {
		result = append(result, fmt.Sprintf("_fishgen%d", index+1)+" "+singleData.tag)
	}
	result = append(result, "_fishgenErr Exception")
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

func generateSingleFunction(fun FunctionInfo) (string, error) {
	if len(fun.receiver) == 0 {
		return "", nil
	}
	receiverName := fun.receiver[0].tag
	if strings.HasSuffix(receiverName, "Model") == false {
		return "", nil
	}
	if fun.name[0] < 'A' || fun.name[0] > 'Z' {
		return "", nil
	}
	for _, singleParam := range fun.params {
		if strings.Index(singleParam.tag, "interface{") != -1 {
			return "", errors.New("invalid has interface type![" + singleParam.tag + "]")
		}
		if strings.Index(singleParam.tag, "func(") != -1 {
			return "", errors.New("invalid has function type![" + singleParam.tag + "]")
		}
	}
	return "func (" + generateSingleField(fun.receiver) + ")" +
		fun.name + "_WithError" +
		"(" + generateSingleField(fun.params) + ")" +
		"(" + generateSingleResult(fun.results) + "){\n" +
		"defer Catch(func(exception Exception) {\n" +
		"_fishgenErr = exception\n" +
		"})\n" +
		generateSingleReturn(fun) +
		"}\n", nil

}

var dirDeclTypeCache map[string]map[string]bool

func getDirDeclTypeInner(dir string) (map[string]bool, error) {
	if IsExistDir(dir) == false {
		return nil,nil
	}
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
	nameImport["Exception"] = ImportInfo{
		name: ".",
		path: "github.com/fishedee/language",
	}
	osGoPath := strings.Split(os.Getenv("GOPATH"), ":")
	for _, singleGoPath := range osGoPath {
		singleGoPath = strings.TrimRight(singleGoPath, "/")
		if singleGoPath == "" {
			continue
		}
		for _, singleParserInfo := range data {
			for _, singleImport := range singleParserInfo.imports {
				if singleImport.name == "_" {
					continue
				}
				if singleImport.name == "." {
					dirImportTypes, err := getDirDeclType(singleGoPath + "/src/" + singleImport.path)
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

func generateSingleFileFunction(data []ParserInfo) (string, error) {
	var result []string
	for _, singleParserInfo := range data {
		for _, singleFun := range singleParserInfo.functions {
			singleFunStr, err := generateSingleFunction(singleFun)
			if err != nil {
				return "", errors.New(singleFun.name + ":" + err.Error())
			}
			result = append(result, singleFunStr)
		}
	}
	return strings.Join(result, "\n"), nil
}

func generateSingleFileContent(data []ParserInfo) (string, error) {
	packageInfo := "package " + data[0].packname + "\n"

	contentInfo, err := generateSingleFileFunction(data)
	if err != nil {
		return "", err
	}

	importInfo, err := generateSingleFileImport(data, packageInfo+contentInfo)
	if err != nil {
		return "", err
	}
	return packageInfo + importInfo + contentInfo, nil
}

func generateSingleTestFileContent(data []ParserInfo) (string, error) {
	packageName := data[0].packname
	packageInfo := "package " + packageName + "\n"
	packageImport := `
		import (
			"testing"
			. "github.com/fishedee/web"
		)`
	packageFunc := "\ntype testFishGenStruct struct{}\n" +
		"func Test" + strings.ToUpper(packageName[0:1]) + packageName[1:] + "(t *testing.T){\n" +
		"RunTest(t,&testFishGenStruct{})\n" +
		"}\n"

	return packageInfo + packageImport + packageFunc, nil
}

func generateSingleFileFormat(filename string, data string) (string, error) {
	result, err := format.Source([]byte(data))
	if err != nil {
		return "", errors.New(err.Error() + "," + data)
	}
	return string(result), nil
}

func generateSingleFileWrite(filename string, data string) error {
	oldData, err := ioutil.ReadFile(filename)
	if err == nil && string(oldData) == data {
		return nil
	}
	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func generateSingleFileTest(dirname string, data []ParserInfo) error {
	filename := dirname + "/" + GetGenerateTestFileName(dirname)
	result, err := generateSingleTestFileContent(data)
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

func generateSingleFileNormal(dirname string, data []ParserInfo) error {
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

func generateSingleFile(dirname string, data []ParserInfo) error {
	err := generateSingleFileNormal(dirname, data)
	if err != nil {
		return err
	}

	err = generateSingleFileTest(dirname, data)
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
			return errors.New(singleKey + ":" + err.Error())
		}
	}
	return nil
}

func IsExistDir(dirName string) bool {
	_, err := ioutil.ReadDir(dirName)
	return err == nil
}

func init() {
	dirDeclTypeCache = map[string]map[string]bool{}
}
