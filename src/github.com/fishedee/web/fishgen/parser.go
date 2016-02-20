package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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
	functions []FunctionInfo
	imports   []ImportInfo
}

func getFieldType(fieldType ast.Expr) string {
	ident, ok := fieldType.(*ast.Ident)
	if ok {
		return ident.Name
	}
	starExpr, ok := fieldType.(*ast.StarExpr)
	if ok {
		return "*" + getFieldType(starExpr.X)
	}
	selectorType, ok := fieldType.(*ast.SelectorExpr)
	if ok {
		return selectorType.Sel.Name + "." + getFieldType(selectorType.X)
	}
	arrayType, ok := fieldType.(*ast.ArrayType)
	if ok {
		return "[]" + getFieldType(arrayType.Elt)
	}
	mapType, ok := fieldType.(*ast.MapType)
	if ok {
		return "map[" + getFieldType(mapType.Key) + "]" + getFieldType(mapType.Value)
	}
	panic(fmt.Sprintf("%#v unknown fieldType ", fieldType))
}

func getFieldListType(fieldList *ast.FieldList) []FieldInfo {
	if fieldList == nil {
		return nil
	}
	var result []FieldInfo
	for _, singleField := range fieldList.List {
		typeName := getFieldType(singleField.Type)
		for _, singleName := range singleField.Names {
			result = append(
				result,
				FieldInfo{
					name: singleName.Name,
					tag:  typeName,
				},
			)
		}
	}
	return result
}

func getFunction(funcDecl *ast.FuncDecl) FunctionInfo {
	return FunctionInfo{
		name:     funcDecl.Name.Name,
		receiver: getFieldListType(funcDecl.Recv),
		params:   getFieldListType(funcDecl.Type.Params),
		results:  getFieldListType(funcDecl.Type.Results),
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
	result.path = imporDecl.Path.Value
	return result
}

func parserSingleFile(filepath string, filename string) (ParserInfo, error) {
	path := filepath + "/" + filename
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return ParserInfo{}, errors.New("file parse error " + err.Error())
	}

	result := ParserInfo{}
	result.file = filename
	for _, singleDecl := range f.Decls {
		singleFuncDecl, ok := singleDecl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		singleFuncInfo := getFunction(singleFuncDecl)
		result.functions = append(result.functions, singleFuncInfo)
	}
	for _, singleImport := range f.Imports {
		singleImportInfo := getImport(singleImport)
		result.imports = append(result.imports, singleImportInfo)
	}
	return result, nil
}

func Parser(data map[string][]os.FileInfo) (map[string][]ParserInfo, error) {
	result := map[string][]ParserInfo{}
	for singleKey, singleDir := range data {
		singleResult := []ParserInfo{}
		for _, singleFile := range singleDir {
			singleFileResult, err := parserSingleFile(singleKey, singleFile.Name())
			if err != nil {
				return nil, err
			}
			singleResult = append(singleResult, singleFileResult)
		}
		result[singleKey] = singleResult
	}
	return result, nil
}
