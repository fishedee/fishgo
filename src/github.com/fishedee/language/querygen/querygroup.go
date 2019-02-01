package main

import (
	. "github.com/fishedee/language"
	"html/template"
	"strings"
)

func QueryGroupGen(request queryGenRequest) *queryGenResponse {
	args := request.args
	line := request.pkg.FileSet().Position(request.expr.Pos()).String()

	//解析第一个参数
	firstArgSlice := getSliceType(line, args[0].Type)
	firstArgSliceElem := firstArgSlice.Elem()

	//解析第二个参数
	secondArgGroupType := getContantStringValue(line, args[1].Value)
	groupType := strings.Trim(secondArgGroupType, " ")
	groupFieldExtract, groupFieldType := getExtendFieldType(line, firstArgSliceElem, groupType)

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
		Throw(1, "%v:groupFunctor argument should be equal with first argument %v!=%v", line, thirdArgFuncArgument[0], firstArgSliceElem)
	}

	//生成函数
	signature := getFunctionSignature(line, args, []bool{false, true, false})
	if hasQueryGroupGenerate[signature] == true {
		return nil
	}
	hasQueryGroupGenerate[signature] = true
	importPackage := map[string]bool{}
	setImportPackage(line, firstArgSliceElem, importPackage)
	setImportPackage(line, thirdArgFuncReturn[0], importPackage)
	setImportPackage(line, groupFieldType, importPackage)
	argumentDefine := getFunctionArgumentCode(line, args, []bool{false, true, false})
	funcBody := excuteTemplate(queryGroupFuncTmpl, map[string]string{
		"signature":          signature,
		"firstArgElemType":   getTypeDeclareCode(line, firstArgSliceElem),
		"thirdArgType":       getTypeDeclareCode(line, thirdArgFunc),
		"thirdArgReturnType": getTypeDeclareCode(line, thirdArgFuncReturn[0]),
		"columnType":         getTypeDeclareCode(line, groupFieldType),
		"columnExtract":      groupFieldExtract,
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
		bufferData := make([]{{ .firstArgElemType }},len(dataIn),len(dataIn))
		mapData := make(map[{{ .columnType }}]int,len(dataIn))
		result := make([]{{ .thirdArgReturnType}},0,len(dataIn))

		language.QueryGroupInternal(len(dataIn),
			func(i int) (int, bool) {
				lastIndex,isExist := mapData[dataIn[i]{{ .columnExtract }}]
				return lastIndex,isExist
			}, func(i int, index int) {
				mapData[dataIn[i]{{ .columnExtract }}] = index
			}, func(k int, i int) {
				bufferData[k] = dataIn[i]
			}, func(i int, j int) {
				single := groupFunctorIn(bufferData[i:j])
				result = append(result,single)
			},
		)
		return result
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
