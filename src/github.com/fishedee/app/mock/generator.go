package main

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
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
	for _, singleData := range data {
		result = append(result, singleData.tag)
	}
	return strings.Join(result, ",")
}

func generateSingleInterface(typeName string, methods []FunctionInfo) string {
	funResult := []string{}
	for _, method := range methods {
		single := method.name + "" +
			"(" + generateSingleField(method.params) + ")" +
			"(" + generateSingleResult(method.results) + ")"
		funResult = append(funResult, single)
	}
	typeName = strings.ToUpper(typeName[0:1]) + typeName[1:]
	return "type I" + typeName + " interface{\n" + strings.Join(funResult, "\n") + "\n}\n"
}

func generateSingleMock(typeName string, methods []FunctionInfo) string {
	funResult := []string{}
	for _, method := range methods {
		methodFieldName := method.name + "Handler"
		single := methodFieldName + " func" +
			"(" + generateSingleField(method.params) + ")" +
			"(" + generateSingleResult(method.results) + ")"
		funResult = append(funResult, single)
	}
	mockTypeStruct := "type " + typeName + "Mock struct{\n" + strings.Join(funResult, "\n") + "\n}\n"

	methodResult := []string{}
	for _, method := range methods {
		single := "func ( this *" + typeName + "Mock)" +
			method.name +
			"(" + generateSingleField(method.params) + ")" +
			"(" + generateSingleResult(method.results) + "){\n"
		paramNames := []string{}
		for _, param := range method.params {
			paramNames = append(paramNames, param.name)
		}
		methodFieldName := method.name + "Handler"
		returnSingle := "this." + methodFieldName + "(" + strings.Join(paramNames, ",") + ")"
		if len(method.results) == 0 {
			single = single + returnSingle + "\n}\n"
		} else {
			single = single + "return " + returnSingle + "\n}\n"
		}
		methodResult = append(methodResult, single)
	}
	return mockTypeStruct + strings.Join(methodResult, "\n")
}

func generateSingleType(typeName string, methodInfo generateTypeInfo) string {
	methods := []FunctionInfo{}
	for _, singleMethodInfo := range methodInfo {
		methods = append(methods, singleMethodInfo)
	}
	sort.Slice(methods, func(i int, j int) bool {
		return methods[i].name < methods[j].name
	})
	result := []string{}
	result = append(result, generateSingleInterface(typeName, methods))
	result = append(result, generateSingleMock(typeName, methods))
	return strings.Join(result, "\n")
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
	for _, singleDeclType := range sourceParserInfo.declType {
		nameImport[singleDeclType] = ImportInfo{}
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

type generateTypeInfo map[string]FunctionInfo

func generateGetTypeInfo(data ParserInfo, result map[string]generateTypeInfo) error {
	reg, err := regexp.CompilePOSIX(Config.typeregex)
	if err != nil {
		return err
	}
	for _, fun := range data.functions {
		if len(fun.receiver) == 0 {
			continue
		}
		receiverName := fun.receiver[0].tag
		if receiverName[0] == '*' {
			receiverName = receiverName[1:]
		}
		if reg.Match([]byte(receiverName)) == false {
			continue
		}
		if fun.name[0] < 'A' || fun.name[0] > 'Z' {
			continue
		}
		for _, singleParam := range fun.params {
			if strings.Index(singleParam.tag, "interface{") != -1 {
				return errors.New("invalid has interface type![" + singleParam.tag + "]")
			}
			if strings.Index(singleParam.tag, "func(") != -1 {
				return errors.New("invalid has function type![" + singleParam.tag + "]")
			}
		}
		functionName := fun.name
		typeInfo, isExist := result[receiverName]
		if isExist == false {
			typeInfo = generateTypeInfo{}
			result[receiverName] = typeInfo
		}
		_, isExist = typeInfo[functionName]
		if isExist == false {
			typeInfo[functionName] = fun
		} else {
			return errors.New("duplicate method![" + receiverName + "." + functionName + "]")
		}
	}
	return nil
}

func generateSingleFileInterface(data []ParserInfo) (string, error) {
	typeInfo := map[string]generateTypeInfo{}
	for _, singleParserInfo := range data {
		err := generateGetTypeInfo(singleParserInfo, typeInfo)
		if err != nil {
			return "", errors.New(singleParserInfo.file + ":" + err.Error())
		}
	}
	typeInfoList := []struct {
		name string
		info generateTypeInfo
	}{}
	for singleTypeName, singleTypeInfo := range typeInfo {
		typeInfoList = append(typeInfoList, struct {
			name string
			info generateTypeInfo
		}{singleTypeName, singleTypeInfo})
	}
	sort.Slice(typeInfoList, func(i int, j int) bool {
		return typeInfoList[i].name < typeInfoList[j].name
	})
	result := []string{}
	for _, singleTypeInfo := range typeInfoList {
		singleResult := generateSingleType(singleTypeInfo.name, singleTypeInfo.info)
		result = append(result, singleResult)
	}
	return strings.Join(result, "\n"), nil
}

func generateSingleFileContent(data []ParserInfo) (string, error) {
	packageInfo := "package " + data[0].packname + "\n"

	contentInfo, err := generateSingleFileInterface(data)
	if err != nil {
		return "", err
	}

	importInfo, err := generateSingleFileImport(data, packageInfo+contentInfo)
	if err != nil {
		return "", err
	}
	return packageInfo + importInfo + contentInfo, nil
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

func init() {
	dirDeclTypeCache = map[string]map[string]bool{}
}
