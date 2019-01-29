package main

import (
	"bytes"
	"flag"
	"fmt"
	. "github.com/fishedee/app/macro"
	. "github.com/fishedee/language"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"os"
)

var (
	recursive      = flag.Bool("r", false, "generate package including sub package")
	queryGenMapper = map[string]queryGenHandler{}
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of querygen:\n")
	fmt.Fprintf(os.Stderr, "\tlarge improve performance language/query.go function \n")
	fmt.Fprintf(os.Stderr, "\tquerygen [flags] [packageName]\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttps://github.com/fishedee/fishgo/tree/master/src/github.com/fishedee/language/querygen\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func nodeString(fset *token.FileSet, n ast.Node) string {
	var buf bytes.Buffer
	format.Node(&buf, fset, n)
	return buf.String()
}

type queryGenRequest struct {
	pkg    MacroPackage
	expr   *ast.CallExpr
	caller *types.Func
	args   []types.TypeAndValue
}

type queryGenResponse struct {
	importPackage map[string]bool
	funcName      string
	funcBody      string
	initBody      string
}

type queryGenHandler func(request queryGenRequest) *queryGenResponse

func handleQueryGen(name string, request queryGenRequest) (*queryGenResponse, bool) {
	handler, isExist := queryGenMapper[name]
	if isExist == false {
		return nil, false
	}
	return handler(request), true
}

func registerQueryGen(name string, handler queryGenHandler) {
	queryGenMapper[name] = handler
}

func run() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		panic("need package name")
	}

	macro := NewMacro()
	if *recursive {
		err := macro.ImportRecursive(args[0])
		if err != nil {
			panic(err)
		}
	} else {
		err := macro.Import(args[0])
		if err != nil {
			panic(err)
		}
	}

	err := macro.Walk(func(pkg MacroPackage) {

		pkg.OnFuncCall(func(expr *ast.CallExpr, caller *types.Func, args []types.TypeAndValue) {
			callerFullName := caller.FullName()
			request := queryGenRequest{
				pkg:    pkg,
				expr:   expr,
				caller: caller,
				args:   args,
			}
			response, hasHandle := handleQueryGen(callerFullName, request)
			if hasHandle == true {
				//should generate
				fmt.Println(response)
			}
		})
		pkg.Inspect()
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("querygen fail: ")
	defer CatchCrash(func(e Exception) {
		log.Fatal(e.GetMessage())
	})
	run()
}
