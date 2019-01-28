package macro

import (
	"bytes"
	"fmt"
	. "github.com/fishedee/assert"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"testing"
)

func nodeString(fset *token.FileSet, n ast.Node) string {
	var buf bytes.Buffer
	format.Node(&buf, fset, n)
	return buf.String()
}

func TestFuncCallInspect(t *testing.T) {
	macro := NewMacro()
	err := macro.ImportRecursive("github.com/fishedee/app")
	AssertEqual(t, err, nil)

	err = macro.Walk(func(pkg MacroPackage) {
		pkg.OnFuncCall(func(expr *ast.CallExpr, caller *types.Func, args []types.TypeAndValue) {
			if caller.Pkg().Path() != "github.com/fishedee/language" {
				return
			}
			fmt.Printf("%v:%v\n", pkg.FileSet().Position(expr.Pos()), nodeString(pkg.FileSet(), expr))
		})
		pkg.Inspect()
	})
	AssertEqual(t, err, nil)
}

func TestFormatSource(t *testing.T) {
	s := `
	package main

	import ("fmt")

	func main(){fmt.Println("123")}
	`
	fmt.Println(NewMacro().FormatSource(s))
}
