package main

import (
	. "github.com/fishedee/language"
	"go/types"
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
	isSliceReturn := ""
	var returnElemType types.Type
	returnSliceType, isSlice := thirdArgFuncReturn[0].(*types.Slice)
	if isSlice {
		isSliceReturn = "..."
		returnElemType = returnSliceType.Elem()
	} else {
		isSliceReturn = ""
		returnElemType = thirdArgFuncReturn[0]
	}
	hasQueryGroupGenerate[signature] = true
	importPackage := map[string]bool{}
	setImportPackage(line, firstArgSliceElem, importPackage)
	setImportPackage(line, returnElemType, importPackage)
	setImportPackage(line, groupFieldType, importPackage)
	argumentDefine := getFunctionArgumentCode(line, args, []bool{false, true, false})
	funcBody := excuteTemplate(queryGroupFuncTmpl, map[string]string{
		"signature":          signature,
		"isFunctorGroup":     "true",
		"firstArgElemType":   getTypeDeclareCode(line, firstArgSliceElem),
		"thirdArgType":       getTypeDeclareCode(line, thirdArgFunc),
		"thirdArgReturnType": getTypeDeclareCode(line, returnElemType),
		"columnType":         getTypeDeclareCode(line, groupFieldType),
		"columnExtract":      groupFieldExtract,
		"isSliceReturn":      isSliceReturn,
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
	 {{if eq .isFunctorGroup "true"}}` +
		"func queryGroup_{{ .signature }}(data interface{},groupType string,groupFunctor interface{})interface{}{\n" +
		`{{else}}` +
		"func queryColumnMap_{{ .signature }}(data interface{},column string)interface{}{\n" +
		`{{end}}dataIn := data.([]{{ .firstArgElemType }})
		bufferData := make([]{{ .firstArgElemType }},len(dataIn),len(dataIn))
		mapData := make(map[{{ .columnType }}]int,len(dataIn))
		{{if eq .isFunctorGroup "true"}}groupFunctorIn := groupFunctor.({{ .thirdArgType }})
		result := make([]{{ .thirdArgReturnType}},0,len(dataIn))
		{{else}}result := make(map[{{ .columnType }}][]{{ .firstArgElemType}},len(dataIn))
		{{end}}
		length := len(dataIn)
		nextData := make([]int, length, length)
		for i := 0; i != length; i++ {
			single := dataIn[i]{{ .columnExtract }}
			lastIndex,isExist := mapData[single]
			if isExist == true {
				nextData[lastIndex] = i
			}
			nextData[i] = -1
			mapData[single] = i
		}
		k := 0
		for i := 0; i != length; i++ {
			j := i
			if nextData[j] == 0 {
				continue
			}
			kbegin := k
			for nextData[j] != -1 {
				nextJ := nextData[j]
				bufferData[k] = dataIn[j]
				nextData[j] = 0
				j = nextJ
				k++
			}
			bufferData[k] = dataIn[j]
			k++
			nextData[j] = 0
			{{if eq .isFunctorGroup "true"}}
			single := groupFunctorIn(bufferData[kbegin:k])
			result = append(result,single{{ .isSliceReturn }})
			{{else}}
			result[bufferData[kbegin]{{ .columnExtract }}] = bufferData[kbegin:k]
			{{end}}
		}
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
