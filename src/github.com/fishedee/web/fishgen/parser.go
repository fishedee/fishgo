package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

type FieldInfo struct {
	name string
	tag  string
}
type FunctionInfo struct {
	name     string
	receiver []FieldInfo
	params   []FieldInfo
	results  []FieldInfo
}

type ImportInfo struct {
	comment string
	name    string
	path    string
}

type ParserInfo struct {
	file      string
	dir       string
	packname  string
	declType  []string
	useType   []string
	functions []FunctionInfo
	imports   []ImportInfo
}

var parserInfoCache map[string]ParserInfo

func getFieldType(fieldType ast.Expr, useType map[string]bool) string {
	ident, ok := fieldType.(*ast.Ident)
	if ok {
		useType[ident.Name] = true
		return ident.Name
	}
	starExpr, ok := fieldType.(*ast.StarExpr)
	if ok {
		return "*" + getFieldType(starExpr.X, useType)
	}
	selectorType, ok := fieldType.(*ast.SelectorExpr)
	if ok {
		return getFieldType(selectorType.X, useType) + "." + selectorType.Sel.Name
	}
	arrayType, ok := fieldType.(*ast.ArrayType)
	if ok {
		return "[]" + getFieldType(arrayType.Elt, useType)
	}
	mapType, ok := fieldType.(*ast.MapType)
	if ok {
		return "map[" + getFieldType(mapType.Key, useType) + "]" + getFieldType(mapType.Value, useType)
	}
	ellipse, ok := fieldType.(*ast.Ellipsis)
	if ok {
		return "..." + getFieldType(ellipse.Elt, useType)
	}
	funcType, ok := fieldType.(*ast.FuncType)
	if ok {
		data := ""
		if funcType.Func != token.NoPos {
			data += "func"
		}
		data += "(" + getFieldListTypeString(funcType.Params, useType) + ")"
		data += "(" + getFieldListTypeString(funcType.Results, useType) + ")"
		return data
	}
	interfaceType, ok := fieldType.(*ast.InterfaceType)
	if ok {
		data := "interface{"
		for _, singleMethod := range interfaceType.Methods.List {
			fieldListInner := []*ast.Field{singleMethod}
			fieldList := &ast.FieldList{List: fieldListInner}

			data += "\n" + getFieldListTypeString(fieldList, useType)
		}
		data += "}"
		return data
	}
	chanType, ok := fieldType.(*ast.ChanType)
	if ok {
		data := "chan "
		data += "\n" + getFieldType(chanType.Value, useType)
		return data
	}
	panic(fmt.Sprintf("%#v unknown fieldType ", fieldType))
}

func getFieldListTypeString(fieldList *ast.FieldList, useType map[string]bool) string {
	var result []string
	data := getFieldListType(fieldList, useType)
	for _, singleData := range data {
		result = append(result, singleData.name+" "+singleData.tag)
	}
	return strings.Join(result, ",")
}

func getFieldListType(fieldList *ast.FieldList, useType map[string]bool) []FieldInfo {
	if fieldList == nil {
		return nil
	}
	var result []FieldInfo
	for _, singleField := range fieldList.List {
		typeName := getFieldType(singleField.Type, useType)
		if singleField.Names != nil {
			for _, singleName := range singleField.Names {
				result = append(
					result,
					FieldInfo{
						name: singleName.Name,
						tag:  typeName,
					},
				)
			}
		} else {
			result = append(
				result,
				FieldInfo{
					name: "",
					tag:  typeName,
				},
			)
		}

	}
	return result
}

func getFunction(funcDecl *ast.FuncDecl, useType map[string]bool) FunctionInfo {
	return FunctionInfo{
		name:     funcDecl.Name.Name,
		receiver: getFieldListType(funcDecl.Recv, useType),
		params:   getFieldListType(funcDecl.Type.Params, useType),
		results:  getFieldListType(funcDecl.Type.Results, useType),
	}
}

func getImport(imporDecl *ast.ImportSpec) ImportInfo {
	result := ImportInfo{}
	if imporDecl.Name != nil {
		result.name = imporDecl.Name.Name
	}
	if imporDecl.Comment != nil && imporDecl.Comment.List != nil && len(imporDecl.Comment.List) != 0 {
		result.comment = imporDecl.Comment.List[0].Text
	}
	result.path = strings.Trim(imporDecl.Path.Value, "\"")
	return result
}

func mapToArray(data map[string]bool) []string {
	result := []string{}
	for singleData, _ := range data {
		result = append(result, singleData)
	}
	return result
}

func parserSingleFile(filename string, source interface{}) (ParserInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, source, parser.ParseComments)
	if err != nil {
		return ParserInfo{}, errors.New("file parse error " + err.Error())
	}

	result := ParserInfo{}
	result.dir = path.Dir(filename)
	result.file = path.Base(filename)
	useType := map[string]bool{}
	for _, singleDecl := range f.Decls {
		singleFuncDecl, ok := singleDecl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		singleFuncInfo := getFunction(singleFuncDecl, useType)
		result.functions = append(result.functions, singleFuncInfo)
	}
	for _, singleImport := range f.Imports {
		singleImportInfo := getImport(singleImport)
		result.imports = append(result.imports, singleImportInfo)
	}
	declType := map[string]bool{}
	for name, singleObject := range f.Scope.Objects {
		if singleObject.Kind != ast.Typ {
			continue
		}
		declType[name] = true
	}
	result.useType = mapToArray(useType)
	result.declType = mapToArray(declType)
	result.packname = f.Name.Name
	return result, nil
}

func ParserSingleSource(source string) (ParserInfo, error) {
	data, err := parserSingleFile("", source)
	if err != nil {
		return ParserInfo{}, err
	}
	return data, nil
}

func ParserSingleFile(path string) (ParserInfo, error) {
	data, ok := parserInfoCache[path]
	if ok {
		return data, nil
	}

	data, err := parserSingleFile(path, nil)
	if err != nil {
		return ParserInfo{}, err
	}
	parserInfoCache[path] = data

	return data, nil
}

func Parser(data map[string][]os.FileInfo) (map[string][]ParserInfo, error) {
	result := map[string][]ParserInfo{}
	for singleKey, singleDir := range data {
		singleResult := []ParserInfo{}
		for _, singleFile := range singleDir {
			singleFileResult, err := ParserSingleFile(singleKey + "/" + singleFile.Name())
			if err != nil {
				return nil, err
			}
			singleResult = append(singleResult, singleFileResult)
		}
		result[singleKey] = singleResult
	}
	return result, nil
}

func init() {
	parserInfoCache = map[string]ParserInfo{}
}
