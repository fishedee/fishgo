package command

import (
	. "github.com/fishedee/language"
	"github.com/fishedee/web/fishcmd/modules"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func Router(argv []string) (string, error) {
	//读取配置
	err := modules.InitConfig()
	if err != nil {
		return "", err
	}

	controllersDir := modules.GetAppDir() + "/controllers"

	dir, err := modules.ReadDir(controllersDir)
	if err != nil {
		return "", err
	}

	importPackage, routerContent, err := extractContents(dir)
	if err != nil {
		return "", err
	}

	err = Generator(importPackage, routerContent)
	if err != nil {
		return "", err
	}

	return "", nil
}

func extractContents(dirInfo map[string][]os.FileInfo) ([]string, []routerContent, error) {

	//获取项目名称
	appName := modules.GetAppName()
	controllersDir := modules.GetAppDir() + "/controllers"

	//遍历类名出来
	routerInitData := map[string][]string{}
	for dirName, files := range dirInfo {

		for _, singleFileInfo := range files {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, dirName+"/"+singleFileInfo.Name(), nil, parser.ParseComments)
			if err != nil {
				modules.Log.Error("ParseFile fail! error: %v", err.Error())
				return nil, nil, err
			}

			for _, s := range f.Scope.Objects {

				typeName := s.Decl.(*ast.TypeSpec).Name.String()

				if strings.HasSuffix(typeName, "Controller") == false {
					continue
				}
				packageName := dirName[len(controllersDir):] //包名
				routerInitData[packageName] = append(routerInitData[packageName], typeName)
			}
		}

	}

	//导入包名
	importPackage := []string{}
	//路由内容
	routerContents := []routerContent{}
	for controllerDir, conterollerNames := range routerInitData {

		//控制器根目录就直接用匿名导入
		if controllerDir == "" {
			importPackage = append(importPackage, `. "`+appName+`/controllers"`)
		} else {
			importPackage = append(importPackage, `"`+appName+"/controllers"+controllerDir+`"`)
		}

		for _, conterollerName := range conterollerNames {

			re := regexp.MustCompile("(.*?)Controller")
			functionName := re.FindStringSubmatch(conterollerName)[1]

			if controllerDir == "" {
				dirName := strings.ToLower("/" + functionName)
				routerContents = append(routerContents, routerContent{dirName, "&" + conterollerName + "{}"})
			} else {
				dirName := strings.ToLower(controllerDir + "/" + functionName)
				routerContents = append(routerContents, routerContent{dirName, "&" + controllerDir[1:] + "." + conterollerName + "{}"})
			}

		}
	}
	return importPackage, routerContents, nil
}

func Generator(importPackage []string, routerContents []routerContent) error {
	//把参数数组排序，防止每次生成内容顺序不一样
	importPackage = ArraySort(importPackage).([]string)
	routerContents = ArrayColumnSort(routerContents, "dirName").([]routerContent)

	//头部
	routerContent :=
		`package routers
		import (
		`

	for _, imports := range importPackage {
		routerContent += imports + "\n"
	}

	routerContent += `. "github.com/fishedee/web"` + " \n )"

	//路由内容
	routerContent += " \n func init() { \n"
	for _, singleRouterContent := range routerContents {
		routerContent += "\n" + `InitRoute("` + singleRouterContent.dirName + `",` + singleRouterContent.ConterollerName + `)`
	}
	routerContent += "}"

	routerFileDir := modules.GetAppDir() + "/routers/router.go"

	//格式化数据
	routerData, err := format.Source([]byte(routerContent))
	if err != nil {
		modules.Log.Error("format.Source fail! error: %v", err.Error())
		return err
	}

	if modules.IsExistFile(routerFileDir) {
		//读取文件
		fileData, err := ioutil.ReadFile(routerFileDir)
		if err != nil {
			modules.Log.Error("ReadFile fail! error: %v", err.Error())
			return err
		}
		//数据比较，如果数据一样的话，就不进行写入文件
		if reflect.DeepEqual(fileData, routerData) {
			return nil
		}
	}

	modules.Log.Debug("%v", "自动编写路由文件")
	//打开文件
	f, err := os.OpenFile(routerFileDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		modules.Log.Error("format.Source fail! error: %v", err.Error())
		return err
	}
	defer f.Close()

	//写入文件(字符串)
	_, err = io.WriteString(f, string(routerData))
	if err != nil {
		modules.Log.Error("WriteString fail! error: %v", err.Error())
		return err
	}
	return nil
}

type routerContent struct {
	dirName         string
	ConterollerName string
}
