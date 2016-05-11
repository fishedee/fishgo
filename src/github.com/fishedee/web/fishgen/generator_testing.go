package main

import (
	"strings"
)

func generateSingleTestFileContent(data []ParserInfo) (string, error) {
	hasModel := false
	for _, singleParseInfo := range data {
		for _, singleStruct := range singleParseInfo.declType {
			if isPublicStructModel(singleStruct) {
				hasModel = true
				break
			}
		}
		if hasModel == true {
			break
		}
	}
	if hasModel == false {
		return "", nil
	}
	packageName := data[0].packname
	packageInfo := "package " + packageName + "\n"
	packageImport := `
		import (
			"testing"
			. "github.com/fishedee/web"
		)`
	packageFunc := "\ntype testFishGenStruct struct{}\n" +
		"func Test" + strings.ToUpper(packageName[0:1]) + packageName[1:] + "(t *testing.T){\n" +
		"RunTest(t,&testFishGenStruct{})\n" +
		"}\n"

	return packageInfo + packageImport + packageFunc, nil
}
