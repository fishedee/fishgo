package main

import (
	. "github.com/fishedee/language"
	"go/types"
	"html/template"
)

func QueryGroupGen(request queryGenRequest) *queryGenResponse {
	args := request.args
	line := request.pkg.FileSet().Position(request.expr.Pos()).String()

	//解析第一个参数
	firstArgSlice := getSliceType(line, args[0].Type)
	firstArgSliceNamed := getNamedType(line, firstArgSlice.Elem())
	firstArgSliceStruct := getStructType(line, firstArgSliceNamed.Underlying())

	//解析第二个参数
	secondArgGroupType := getContantStringValue(line, args[1].Value)
	sortFieldNames, sortFieldIsAscs := analyseSortType(secondArgGroupType)
	sortFieldTypes := make([]types.Type, len(sortFieldNames), len(sortFieldNames))
	for i, sortFieldName := range sortFieldNames {
		sortFieldTypes[i] = getFieldType(line, firstArgSliceStruct, sortFieldName)
	}

	//解析第三个参数
	thirdArgFunc := getFunctionType(line, args[2].Type)
	thirdArgFuncArgument := getArgumentType(line, thirdArgFunc)
	thirdArgFuncReturn := getReturnType(line, thirdArgFunc)
	if len(thirdArgFuncArgument) != 1 {
		Throw(1, "%v:should be one argument", line)
	}
	if len(thirdArgFuncReturn) != 1 {
		Throw(1, "%v:should be one return", line)
	}
	if thirdArgFuncArgument[0].String() != firstArgSlice.String() {
		Throw(1, "%v:groupFunctor argument should be equal with first argument %v!=%v", thirdArgFuncArgument[0], firstArgSliceNamed)
	}

	//生成函数
	signature := getFunctionSignature(line, args, []bool{false, true, false})
	if hasQueryGroupGenerate[signature] == true {
		return nil
	}
	hasQueryGroupGenerate[signature] = true
	importPackage := map[string]bool{}
	setImportPackage(line, firstArgSliceNamed, importPackage)
	setImportPackage(line, thirdArgFuncReturn[0], importPackage)
	argumentDefine := getFunctionArgumentCode(line, args, []bool{false, true, false})
	funcBody := excuteTemplate(queryGroupFuncTmpl, map[string]string{
		"signature":          signature,
		"firstArgElemType":   getTypeDeclareCode(line, firstArgSliceNamed),
		"thirdArgType":       getTypeDeclareCode(line, thirdArgFunc),
		"thirdArgReturnType": getTypeDeclareCode(line, thirdArgFuncReturn[0]),
		"sortCode":           getCombineLessCompareCode(line, "newData[i]", "newData[j]", sortFieldNames, sortFieldIsAscs, sortFieldTypes),
	})
	initBody := excuteTemplate(queryGroupInitTmpl, map[string]string{
		"signature":      signature,
		"argumentDefine": argumentDefine,
	})
	return &queryGenResponse{
		importPackage: importPackage,
		funcName:      "queryGroup_" + signature,
		funcBody:      funcBody,
		initBody:      initBody,
	}
}

var (
	queryGroupFuncTmpl    *template.Template
	queryGroupInitTmpl    *template.Template
	hasQueryGroupGenerate map[string]bool
)

func init() {
	var err error
	queryGroupFuncTmpl, err = template.New("name").Parse(`
	func queryGroup_{{ .signature }}(data interface{},groupType string,groupFunctor interface{})interface{}{
		dataIn := data.([]{{ .firstArgElemType }})
		groupFunctorIn := groupFunctor.({{ .thirdArgType }})
		newData := make([]{{ .firstArgElemType }},len(dataIn),len(dataIn))
		copy(newData,dataIn)
		newData2 := make([]{{ .thirdArgReturnType}},0,len(dataIn))

		language.QueryGroupInternal(len(newData),func(i int, j int)int{
			{{ .sortCode }}
			return 0
		},func(i int,j int){
			newData[j] , newData[i] = newData[i],newData[j]
		},func(i int,j int){
			single := groupFunctorIn(newData[i:j])
			newData2 = append(newData2,single)
		})
		return newData2
	}
	`)
	if err != nil {
		panic(err)
	}
	queryGroupInitTmpl, err = template.New("name").Parse(`
		language.QueryGroupMacroRegister({{.argumentDefine}},queryGroup_{{.signature}})
	`)
	if err != nil {
		panic(err)
	}
	registerQueryGen("github.com/fishedee/language.QueryGroup", QueryGroupGen)
	hasQueryGroupGenerate = map[string]bool{}
}
