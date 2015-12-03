package web

import (
	"github.com/astaxie/beego"
	"github.com/fishedee/language"
	"strings"
	"reflect"
	"strconv"
)

type BeegoValidateController struct {
	beego.Controller
}

func (this *BeegoValidateController)Prepare(){
	this.Controller.Prepare();
}

func (this *BeegoValidateController)Finish(){
	this.Controller.Finish();
}

func (this *BeegoValidateController)runMethod(method reflect.Value,arguments []reflect.Value)(result []reflect.Value){
	defer language.Catch(func(exception language.Exception){
		Log.Error(exception.GetStackTrace())
		result = []reflect.Value{reflect.ValueOf(exception)}
	})
	result = method.Call(arguments)
	if len(result) == 0 {
		result = []reflect.Value{reflect.ValueOf(nil)}
	}
	return
}

func (this *BeegoValidateController)AutoRouteMethod(){
	//查找路由
	appController := this.AppController
	url := this.Ctx.Input.Url()
	urlArray := strings.Split(url,"/")
	if len(urlArray) == 0{
		panic("unknown url segement"+url)
	}
	lastUrlSegment := urlArray[ len(urlArray) - 1 ]

	//执行路由
	appControllerValue := reflect.ValueOf(appController)
	methodName := firstUpperName(lastUrlSegment)
	appControllerResult := this.runMethod(appControllerValue.MethodByName(methodName),[]reflect.Value{})
	
	//处理返回值
	if len(appControllerResult) != 1 {
		panic("url controller should has return value "+url)
	}
	appControllerValueResult := []reflect.Value{appControllerResult[0]}
	appControllerValue.MethodByName("AutoRender").Call(appControllerValueResult)
}

func (this *BeegoValidateController)AutoRender(result interface{}){

}

func (this *BeegoValidateController)check(requireStruct interface{}){	
	//获取require字段
	requireStructType := reflect.TypeOf(requireStruct).Elem()
	requireStructValue := reflect.ValueOf(requireStruct).Elem()
	for i := 0 ; i != requireStructType.NumField() ; i++{
		singleRequireStruct := requireStructType.Field(i)
		singleRequireStructName := firstLowerName(singleRequireStruct.Name)

		result := this.Ctx.Input.Query(singleRequireStructName)

		singleRequireStructValue := requireStructValue.Field(i)
		singleRequireStructValueKind := singleRequireStructValue.Kind()
		if singleRequireStructValueKind == reflect.String{
			singleRequireStructValue.SetString(result)
		}else if singleRequireStructValueKind == reflect.Int{
			var resultInt int
			var err error
			if result == ""{
				resultInt = 0
			}else{
				resultInt,err = strconv.Atoi(result)
				if err != nil{
					language.Throw(1,"参数"+singleRequireStructName+"不是合法的整数，其值为：["+result+"]")
				}
			}
			singleRequireStructValue.SetInt( int64(resultInt) )
		}else{
			language.Throw(1,"不合法的参数"+singleRequireStructValueKind.String())
		}
	}
}

func (this *BeegoValidateController)CheckGet(requireStruct interface{}){
	if this.Ctx.Input.Method() != "GET"{
		language.Throw(1,"请求Method不是Get方法")
	}
	this.check(requireStruct)
}

func (this *BeegoValidateController)CheckPost(requireStruct interface{}){
	if this.Ctx.Input.Method() != "POST"{
		language.Throw(1,"请求Method不是POST方法")
	}
	this.check(requireStruct)
}

var vaildateControllerMethod map[string]bool

func init(){
	vaildateControllerMethod = map[string]bool{}
	vaildateControllerType := reflect.TypeOf(&BeegoValidateController{}).Elem()
	for i := 0 ; i != vaildateControllerType.NumMethod() ; i++{
		singleMethod := vaildateControllerType.Method(i).Name
		vaildateControllerMethod[ singleMethod ] = true
	}
}

func firstLowerName(name string)(string){
	return strings.ToLower(name[0:1]) + name[1:]
}

func firstUpperName(name string)(string){
	return strings.ToUpper(name[0:1]) + name[1:]
}

func InitBeegoVaildateControllerRoute(namespace string,target beego.ControllerInterface){
	controllerType := reflect.TypeOf(target)
	for i := 0 ; i != controllerType.NumMethod() ; i++{
		singleMethod := firstLowerName(controllerType.Method(i).Name)
		if _,ok := vaildateControllerMethod[singleMethod];ok{
			continue
		}
		beego.Router(
			namespace+"/"+singleMethod,
			target,
			"*:AutoRouteMethod",
		);
	}
}