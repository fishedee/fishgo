package main

import (
	"strings"
)
func generateSingleTestFileContent(data []ParserInfo) (string, error) {
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
