package macro

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/loader"
	"io/ioutil"
	"os"
)

type MacroFuncCallInspector func(expr *ast.CallExpr, caller *types.Func, args []types.TypeAndValue)

type MacroPackage struct {
	fset     *token.FileSet
	pkg      *loader.PackageInfo
	funcCall MacroFuncCallInspector
}

func (this *MacroPackage) Package() *types.Package {
	return this.pkg.Pkg
}

func (this *MacroPackage) FileSet() *token.FileSet {
	return this.fset
}

func (this *MacroPackage) TypeInfo() types.Info {
	return this.pkg.Info
}

func (this *MacroPackage) OnFuncCall(funcCall MacroFuncCallInspector) {
	this.funcCall = funcCall
}

func (this *MacroPackage) fireFuncCall(n ast.Node) {
	if this.funcCall == nil {
		return
	}

	expr, ok := n.(*ast.CallExpr)
	if ok == false {
		return
	}

	//获取caller信息
	exprIdent, ok := expr.Fun.(*ast.Ident)
	if ok == false {
		return
	}

	info := this.pkg.Info
	funcObj, ok := info.Uses[exprIdent].(*types.Func)
	if ok == false {
		return
	}

	//获取argument信息
	typeAndValues := []types.TypeAndValue{}
	for _, arg := range expr.Args {
		t1, isExist := info.Types[arg]
		if isExist == false {
			panic(fmt.Sprintf("unknown argument type:%v", expr.Args))
		}
		typeAndValues = append(typeAndValues, t1)
	}

	//触发
	this.funcCall(expr, funcObj, typeAndValues)
}

func (this *MacroPackage) Inspect() {
	for _, file := range this.pkg.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			//检查函数
			this.fireFuncCall(n)
			return true
		})
	}
}

type MacroWalker func(inspector MacroPackage)

type Macro struct {
	packages map[string]bool
}

func NewMacro() *Macro {
	return &Macro{
		packages: map[string]bool{},
	}
}

func (this *Macro) Import(pkg string) error {
	this.packages[pkg] = true
	return nil
}

func (this *Macro) getAllDir(baseDir string, pkgName string) ([]string, error) {
	files, err := ioutil.ReadDir(baseDir + "/" + pkgName)
	if err != nil {
		return nil, err
	}
	result := []string{}
	result = append(result, pkgName)
	for _, file := range files {
		if file.IsDir() {
			subPackageName := pkgName + "/" + file.Name()
			subPackages, err := this.getAllDir(baseDir, subPackageName)
			if err != nil {
				return nil, err
			}
			result = append(result, subPackages...)
		}
	}
	return result, nil
}

func (this *Macro) ImportRecursive(pkg string) error {
	gopath, _ := os.LookupEnv("GOPATH")
	allPackage, err := this.getAllDir(gopath+"/src", pkg)
	if err != nil {
		return err
	}
	for _, packageSingle := range allPackage {
		this.packages[packageSingle] = true
	}
	return nil
}

func (this *Macro) Walk(walker MacroWalker) error {
	var conf loader.Config
	if len(this.packages) == 0 {
		return errors.New("none package have to load")
	}
	for singlePackage, _ := range this.packages {
		conf.Import(singlePackage)
	}
	lprog, err := conf.Load()
	if err != nil {
		return err
	}
	for _, singlePackage := range lprog.Imported {
		if len(singlePackage.Errors) != 0 {
			return singlePackage.Errors[0]
		}
		inspector := MacroPackage{
			fset: lprog.Fset,
			pkg:  singlePackage,
		}
		walker(inspector)
	}
	return nil
}

func (this *Macro) FormatSource(data string) (string, error) {
	result, err := format.Source([]byte(data))
	if err != nil {
		return "", errors.New(err.Error() + "," + data)
	}
	return string(result), nil
}
