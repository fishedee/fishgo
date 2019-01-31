package main

import (
	"go/types"
	"html/template"
)

func QuerySortGen(request queryGenRequest) *queryGenResponse {
	args := request.args
	line := request.pkg.FileSet().Position(request.expr.Pos()).String()

	//解析第一个参数
	firstArgSlice := getSliceType(line, args[0].Type)
	firstArgSliceNamed := getNamedType(line, firstArgSlice.Elem())
	firstArgSliceStruct := getStructType(line, firstArgSliceNamed.Underlying())

	//解析第二个参数
	secondArgSortType := getContantStringValue(line, args[1].Value)
	sortFieldNames, sortFieldIsAscs := analyseSortType(secondArgSortType)
	sortFieldTypes := make([]types.Type, len(sortFieldNames), len(sortFieldNames))
	for i, sortFieldName := range sortFieldNames {
		sortFieldTypes[i] = getFieldType(line, firstArgSliceStruct, sortFieldName)
	}

	//生成函数
	signature := getFunctionSignature(line, args, []bool{false, true})
	if hasQuerySortGenerate[signature] == true {
		return nil
	}
	hasQuerySortGenerate[signature] = true
	importPackage := map[string]bool{}
	setImportPackage(line, firstArgSliceNamed, importPackage)
	argumentDefine := getFunctionArgumentCode(line, args, []bool{false, true})
	funcBody := excuteTemplate(querySortFuncTmpl, map[string]string{
		"signature":        signature,
		"firstArgElemType": getTypeDeclareCode(line, firstArgSliceNamed),
		"sortCode":         getCombineLessCompareCode(line, "newData[i]", "newData[j]", sortFieldNames, sortFieldIsAscs, sortFieldTypes),
	})
	initBody := excuteTemplate(querySortInitTmpl, map[string]string{
		"signature":      signature,
		"argumentDefine": argumentDefine,
	})
	return &queryGenResponse{
		importPackage: importPackage,
		funcName:      "querySort_" + signature,
		funcBody:      funcBody,
		initBody:      initBody,
	}
}

var (
	querySortFuncTmpl    *template.Template
	querySortInitTmpl    *template.Template
	hasQuerySortGenerate map[string]bool
)

func init() {
	var err error
	querySortFuncTmpl, err = template.New("name").Parse(`
	func querySort_{{ .signature }}(data interface{},sortType string)interface{}{
		dataIn := data.([]{{ .firstArgElemType }})
		newData := make([]{{ .firstArgElemType }},len(dataIn),len(dataIn))
		copy(newData,dataIn)

		language.QuerySortInternal(len(newData),func(i int, j int)int{
			{{ .sortCode }}
			return 0
		},func(i int,j int){
			newData[j] , newData[i] = newData[i],newData[j]
		})
		return newData
	}
	`)
	if err != nil {
		panic(err)
	}
	querySortInitTmpl, err = template.New("name").Parse(`
		language.QuerySortMacroRegister({{.argumentDefine}},querySort_{{.signature}})
	`)
	if err != nil {
		panic(err)
	}
	registerQueryGen("github.com/fishedee/language.QuerySort", QuerySortGen)
	hasQuerySortGenerate = map[string]bool{}
}
