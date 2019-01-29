package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	. "github.com/fishedee/language"
	"go/constant"
	"go/types"
	"html/template"
)

func getFunctionSignature(line string, arguments []types.TypeAndValue, isConstant []bool) string {
	var buffer bytes.Buffer
	for i, argument := range arguments {
		single := ""
		if isConstant[i] == true {
			single = getContantStringValue(line, argument.Value)
		} else {
			single = argument.Type.String()
		}
		buffer.WriteString("_" + single)
	}
	hash := sha1.New()
	hash.Write(buffer.Bytes())
	etag := hash.Sum(nil)
	etagString := hex.EncodeToString(etag)
	return etagString
}

func getContantStringValue(line string, value constant.Value) string {
	if value == nil {
		Throw(1, "%v:should be constant!%v", line, value)
	}
	return constant.StringVal(value)
}

func getNamedType(line string, t types.Type) *types.Named {
	t1, isNamed := t.(*types.Named)
	if isNamed == false {
		Throw(1, "%v:should be named type!%v", t)
	}
	return t1
}

func getSliceType(line string, t types.Type) *types.Slice {
	t1, isSlice := t.(*types.Slice)
	if isSlice == false {
		Throw(1, "%v:should be slice type!%v", t)
	}
	return t1
}

func getStructType(line string, t types.Type) *types.Struct {
	t1, isStruct := t.(*types.Struct)
	if isStruct == false {
		Throw(1, "%v:should be struct type!%v", line, t)
	}
	return t1
}

func getFieldType(line string, tStruct *types.Struct, column string) types.Type {
	for i := 0; i != tStruct.NumFields(); i++ {
		field := tStruct.Field(i)
		if field.Name() == column {
			return field.Type()
		}
	}
	Throw(1, "%v:%v has not found column %v", line, tStruct, column)
	return nil
}

func getTypeDeclareCode(line string, t types.Type) string {
	if tBasic, ok := t.(*types.Basic); ok {
		switch tBasic.Kind() {
		case types.Bool:
			return "bool"
		case types.Int:
			return "int"
		case types.String:
			return "string"
		default:
			Throw(1, "%v:unknown basic type %v", t.String())
			return ""
		}
	} else if tSlice, ok := t.(*types.Slice); ok {
		elemType := tSlice.Elem()
		return "[]" + getTypeDeclareCode(line, elemType)
	} else if tMap, ok := t.(*types.Map); ok {
		keyType := getTypeDeclareCode(line, tMap.Key())
		elemType := getTypeDeclareCode(line, tMap.Elem())
		return "map[" + keyType + "]" + elemType
	} else if tNamed, ok := t.(*types.Named); ok {
		obj := tNamed.Obj()
		if obj.Pkg().Path() == globalGeneratePackagePath {
			return obj.Name()
		} else {
			return obj.Pkg().Name() + "." + obj.Name()
		}
	} else {
		Throw(1, "%v:unknown type to declare: %v", t.String())
		return ""
	}

}

func getTypeDefineCodeInner(line string, t types.Type, isTop bool) string {
	if tBasic, ok := t.(*types.Basic); ok {
		switch tBasic.Kind() {
		case types.Bool:
			return "false"
		case types.Int:
			return "0"
		case types.String:
			return "\"\""
		default:
			Throw(1, "%v:unknown basic type %v", t.String())
			return ""
		}
	} else if tSlice, ok := t.(*types.Slice); ok {
		elemType := tSlice.Elem()
		return "[]" + getTypeDefineCodeInner(line, elemType, false) + "{}"
	} else if tMap, ok := t.(*types.Map); ok {
		keyType := getTypeDefineCodeInner(line, tMap.Key(), false)
		elemType := getTypeDefineCodeInner(line, tMap.Elem(), false)
		return "map[" + keyType + "]" + elemType + "{}"
	} else if tNamed, ok := t.(*types.Named); ok {
		obj := tNamed.Obj()
		underType := tNamed.Underlying()
		if _, isStruct := underType.(*types.Struct); isStruct {
			declareName := ""
			if obj.Pkg().Path() == globalGeneratePackagePath {
				declareName = obj.Name()
			} else {
				declareName = obj.Pkg().Name() + "." + obj.Name()
			}
			if isTop == true {
				declareName = declareName + "{}"
			}
			return declareName
		} else {
			underTypeDefine := getTypeDefineCode(line, underType)
			if obj.Pkg().Path() == globalGeneratePackagePath {
				return obj.Name() + "(" + underTypeDefine + ")"
			} else {
				return obj.Pkg().Name() + "." + obj.Name() + "(" + underTypeDefine + ")"
			}
		}
	} else {
		Throw(1, "%v:unknown type to define %v", line, t.String())
		return ""
	}
}

func setImportPackage(line string, t types.Type, importPkg map[string]bool) {
	if _, ok := t.(*types.Basic); ok {
		return
	} else if tSlice, ok := t.(*types.Slice); ok {
		elemType := tSlice.Elem()
		setImportPackage(line, elemType, importPkg)
	} else if tMap, ok := t.(*types.Map); ok {
		keyType := tMap.Key()
		elemType := tMap.Elem()
		setImportPackage(line, keyType, importPkg)
		setImportPackage(line, elemType, importPkg)
	} else if tNamed, ok := t.(*types.Named); ok {
		obj := tNamed.Obj()
		pkg := obj.Pkg()
		importPkg[pkg.Path()] = true
	} else {
		Throw(1, "%v:unknown type to define %v", line, t.String())
	}
}

func getTypeDefineCode(line string, t types.Type) string {
	return getTypeDefineCodeInner(line, t, true)
}

func getCompareCode(line string, name1 string, name1Type types.Type, name2 string, name2Type types.Type) string {
	//FIXME
	return ""
}

func getFunctionArgumentCode(line string, arguments []types.TypeAndValue, isConstant []bool) string {
	argvs := []string{}
	for i, argument := range arguments {
		if isConstant[i] == true {
			argvs = append(argvs, "\""+getContantStringValue(line, argument.Value)+"\"")
		} else {
			argvs = append(argvs, getTypeDefineCode(line, argument.Type))
		}
	}
	return Implode(argvs, ",")
}

func excuteTemplate(tmpl *template.Template, data map[string]string) string {
	newData := make(map[string]template.HTML, len(data))
	for key, value := range data {
		newData[key] = template.HTML(value)
	}
	var buffer bytes.Buffer
	err := tmpl.Execute(&buffer, newData)
	if err != nil {
		Throw(1, "execute fail %v", err)
	}
	return buffer.String()
}

var (
	globalGeneratePackagePath = ""
)
