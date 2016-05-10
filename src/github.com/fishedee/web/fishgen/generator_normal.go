package main

import (
	"errors"
	"fmt"
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

func generateSingleResultWithException(data []FieldInfo) string {
	var result []string
	for index, singleData := range data {
		result = append(result, fmt.Sprintf("_fishgen%d", index+1)+" "+singleData.tag)
	}
	result = append(result, "_fishgenErr Exception")
	return strings.Join(result, ",")
}

func generateSingleResult(data []FieldInfo) string {
	var result []string
	for index, singleData := range data {
		result = append(result, fmt.Sprintf("_fishgen%d", index+1)+" "+singleData.tag)
	}
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
	if isPublicFunction(fun.name) == false {
		return "", nil
	}
	if isPublicStructModel(receiverName) == false {
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
		"(" + generateSingleResultWithException(fun.results) + "){\n" +
		"defer Catch(func(exception Exception) {\n" +
		"_fishgenErr = exception\n" +
		"})\n" +
		generateSingleReturn(fun) +
		"}\n", nil
}

func isPublicStructController(name string) bool {
	if strings.HasSuffix(name, "Controller") == true {
		return true
	}
	return false
}

func isPublicStructModel(name string) bool {
	if strings.HasSuffix(name, "Model") == true {
		return true
	}
	return false
}

func isPublicStructTest(name string) bool {
	if strings.HasSuffix(name, "Test") == true {
		return true
	}
	return false
}

func isPublicStruct(name string) bool {
	if isPublicStructController(name) ||
		isPublicStructModel(name) ||
		isPublicStructTest(name) {
		return true
	}
	return false
}

func firstUpper(name string) string {
	return strings.ToUpper(name[0:1]) + name[1:]
}

func isPublicFunction(name string) bool {
	if name[0] < 'A' || name[0] > 'Z' {
		return false
	}
	return true

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

func generateSingleFileInit(data []ParserInfo) (string, error) {
	//生成init
	result := "func init(){\n"
	initIndex := 0
	for _, singleParseInfo := range data {
		for _, singleStruct := range singleParseInfo.declType {
			var funName string
			if isPublicStructController(singleStruct) {
				funName = "InitController"
			} else if isPublicStructModel(singleStruct) {
				funName = "InitModel"
			} else if isPublicStructTest(singleStruct) {
				funName = "InitTest"
			} else {
				continue
			}
			result += fmt.Sprintf("v%v := %v(&%v{})\n%v(&v%v)\n", initIndex, firstUpper(singleStruct), singleStruct, funName, initIndex)
			initIndex++
		}
	}
	result += "}\n"
	return result, nil
}

func generateSingleFileInterface(data []ParserInfo) (string, error) {
	//提取所有接口
	structInterface := map[string][]string{}
	for _, singleParserInfo := range data {
		for _, singleFun := range singleParserInfo.functions {
			if len(singleFun.receiver) == 0 {
				continue
			}
			receiverName := singleFun.receiver[0].tag
			receiverName = strings.Trim(receiverName, "*")
			if isPublicFunction(singleFun.name) == false {
				continue
			}
			if isPublicStruct(receiverName) == false {
				continue
			}
			funSignature := singleFun.name + "(" + generateSingleField(singleFun.params) + ")(" + generateSingleResult(singleFun.results) + ")"
			structInterface[receiverName] = append(structInterface[receiverName], funSignature)
			if isPublicStructModel(receiverName) {
				funSignature := singleFun.name + "_WithError(" + generateSingleField(singleFun.params) + ")(" + generateSingleResultWithException(singleFun.results) + ")"
				structInterface[receiverName] = append(structInterface[receiverName], funSignature)
			}
		}
	}

	//生成接口
	result := []string{}
	for _, singleParseInfo := range data {
		for _, singleStruct := range singleParseInfo.declType {
			if isPublicStruct(singleStruct) == false {
				continue
			}
			singleResult := "type " + firstUpper(singleStruct) + " interface{\n"
			for _, singleFun := range structInterface[singleStruct] {
				singleResult += singleFun + "\n"
			}
			singleResult += "}\n"
			result = append(result, singleResult)
		}
	}
	return strings.Join(result, "\n"), nil
}

func generateSingleFileContent(data []ParserInfo) (string, error) {
	packageInfo := "package " + data[0].packname + "\n"

	interfaceInfo, err := generateSingleFileInterface(data)
	if err != nil {
		return "", err
	}

	contentInfo, err := generateSingleFileFunction(data)
	if err != nil {
		return "", err
	}

	initInfo, err := generateSingleFileInit(data)
	if err != nil {
		return "", err
	}

	content := interfaceInfo + contentInfo + initInfo
	importInfo, err := generateSingleFileImport(data, packageInfo+content)
	if err != nil {
		return "", err
	}
	return packageInfo + importInfo + content, nil
}
