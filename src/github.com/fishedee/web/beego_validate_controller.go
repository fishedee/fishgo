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
		Log.Error("Buiness Error Code:[%d] Message:[%s]\nStackTrace:[%s]",exception.GetCode(),exception.GetMessage(),exception.GetStackTrace())
		result = []reflect.Value{reflect.ValueOf(exception)}
	})
	result = method.Call(arguments)
	if len(result) == 0 {
		result = []reflect.Value{reflect.Zero( reflect.TypeOf(this) )}
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
	appControllerType := reflect.TypeOf(appController)
	appControllerValue := reflect.ValueOf(appController)
	appMethodInfos,ok := routerControllerMethod[appControllerType]
	if !ok{
		panic("appMethodInfo not found router "+appControllerType.String())
	}
	appMethodInfo,ok := appMethodInfos[lastUrlSegment]
	if !ok{
		panic("appMethodInfo not found router "+appControllerType.String()+","+lastUrlSegment)
	}
	appControllerResult := this.runMethod(appMethodInfo.method.Func,[]reflect.Value{appControllerValue})
	
	//处理返回值
	if len(appControllerResult) != 1 {
		panic("url controller should has return value "+url)
	}
	appControllerValueResult := []reflect.Value{appControllerResult[0],reflect.ValueOf(appMethodInfo.viewName)}
	appControllerValue.MethodByName("AutoRender").Call(appControllerValueResult)
}

func (this *BeegoValidateController)AutoRender(result interface{},view string){

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
			//FIXME 缺少time的解析
			//language.Throw(1,"不合法的参数"+singleRequireStructValueKind.String())
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

type methodInfo struct{
	name string
	viewName string
	method reflect.Method
}
var routerControllerMethod map[reflect.Type]map[string]methodInfo

func init(){
	routerControllerMethod = map[reflect.Type]map[string]methodInfo{}
}

func firstLowerName(name string)(string){
	return strings.ToLower(name[0:1]) + name[1:]
}

func firstUpperName(name string)(string){
	return strings.ToUpper(name[0:1]) + name[1:]
}

func isPublic(name string)bool{
	fisrtStr := name[0:1]
	if fisrtStr >= "A" && fisrtStr <= "Z"{
		return true
	}else{
		return false
	}
}

func getMethodInfo(name string)(*methodInfo){
	data := strings.Split(name,"_")
	if len(data) != 2{
		return nil
	}
	return &methodInfo{
		name:firstLowerName(data[0]),
		viewName:firstLowerName(data[1]),
	}
}

func InitBeegoVaildateControllerRoute(namespace string,target beego.ControllerInterface){
	controllerType := reflect.TypeOf(target)
	routerControllerMethod[controllerType] = map[string]methodInfo{}

	for i := 0 ; i != controllerType.NumMethod() ; i++{
		singleMethod := controllerType.Method(i)
		singleMethodName := singleMethod.Name
		if isPublic( singleMethodName ) == false{
			continue
		}
		singleMethodInfo := getMethodInfo( singleMethodName )
		if singleMethodInfo == nil{
			continue
		}
		singleMethodInfo.method = singleMethod

		beego.Router(
			namespace+"/"+singleMethodInfo.name,
			target,
			"*:AutoRouteMethod",
		);
		routerControllerMethod[controllerType][singleMethodInfo.name] = *singleMethodInfo
	}
}