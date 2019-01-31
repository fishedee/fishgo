package main

import (
	"html/template"
	"strings"
)

func QueryColumnMapGen(request queryGenRequest) *queryGenResponse {
	args := request.args
	line := request.pkg.FileSet().Position(request.expr.Pos()).String()

	//解析第二个参数
	secondArgValue := getContantStringValue(line, args[1].Value)
	column := strings.Trim(secondArgValue, " ")

	//解析第一个参数
	firstArgSlice := getSliceType(line, args[0].Type)
	firstArgSliceNamed := getNamedType(line, firstArgSlice.Elem())
	firstArgSliceStruct := getStructType(line, firstArgSliceNamed.Underlying())
	columnArgType := getFieldType(line, firstArgSliceStruct, column)

	//生成函数
	signature := getFunctionSignature(line, args, []bool{false, true})
	if hasQueryColumnMapGenerate[signature] == true {
		return nil
	}
	hasQueryColumnMapGenerate[signature] = true
	importPackage := map[string]bool{}
	setImportPackage(line, firstArgSliceNamed, importPackage)
	setImportPackage(line, columnArgType, importPackage)
	argumentDefine := getFunctionArgumentCode(line, args, []bool{false, true})
	funcBody := excuteTemplate(queryColumnMapFuncTmpl, map[string]string{
		"signature":              signature,
		"firstArgElemType":       getTypeDeclareCode(line, firstArgSliceNamed),
		"firstArgElemColumnType": getTypeDeclareCode(line, columnArgType),
		"column":                 column,
	})
	initBody := excuteTemplate(queryColumnMapInitTmpl, map[string]string{
		"signature":      signature,
		"argumentDefine": argumentDefine,
	})
	return &queryGenResponse{
		importPackage: importPackage,
		funcName:      "queryColumnMap_" + signature,
		funcBody:      funcBody,
		initBody:      initBody,
	}
}

var (
	queryColumnMapFuncTmpl    *template.Template
	queryColumnMapInitTmpl    *template.Template
	hasQueryColumnMapGenerate map[string]bool
)

func init() {
	var err error
	queryColumnMapFuncTmpl, err = template.New("name").Parse(`
	func queryColumnMap_{{ .signature }}(data interface{},column string)interface{}{
		dataIn := data.([]{{ .firstArgElemType }})
		result := make(map[{{ .firstArgElemColumnType }}]{{ .firstArgElemType }},len(dataIn))

		for _,single := range dataIn{
			result[single.{{ .column }}] = single
		}
		return result
	}
	`)
	if err != nil {
		panic(err)
	}
	queryColumnMapInitTmpl, err = template.New("name").Parse(`
		language.QueryColumnMapMacroRegister({{.argumentDefine}},queryColumnMap_{{.signature}})
	`)
	if err != nil {
		panic(err)
	}
	registerQueryGen("github.com/fishedee/language.QueryColumnMap", QueryColumnMapGen)
	hasQueryColumnMapGenerate = map[string]bool{}
}
