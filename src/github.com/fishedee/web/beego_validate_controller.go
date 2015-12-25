package web

import (
	"github.com/astaxie/beego"
	"github.com/fishedee/language"
	"net/url"
	"strings"
	"reflect"
	"strconv"
	"time"
	"sync"
)

type beegoValidateControllerInfo struct{
	getBasicFunc int
}

var beegoValidateControllerInfoMap struct{
	mutex sync.RWMutex
	data map[reflect.Type]beegoValidateControllerInfo
}

var beegoValidateControllerType reflect.Type

var routerControllerMethod map[reflect.Type]map[string]methodInfo

func init(){
	beegoValidateControllerInfoMap.data = map[reflect.Type]beegoValidateControllerInfo{}
	beegoValidateControllerType = reflect.TypeOf(BeegoValidateController{})
	routerControllerMethod = map[reflect.Type]map[string]methodInfo{}
}

type BeegoValidateController struct {
	beego.Controller
	*BeegoValidateBasic
}

func (this *BeegoValidateController)Prepare(){
	this.BeegoValidateBasic = NewBeegoValidateBasic(this.Ctx)
	PrepareBeegoValidateModel(this.AppController)
}

func (this *BeegoValidateController)Finish(){
	FinishBeegoValidateModel(this.AppController)
}

func (this *BeegoValidateController)GetBasic()(*BeegoValidateBasic){
	return this.BeegoValidateBasic
}

func (this *BeegoValidateController)runMethod(method reflect.Value,arguments []reflect.Value)(result []reflect.Value){
	defer language.Catch(func(exception language.Exception){
		this.Log.Error("Buiness Error Code:[%d] Message:[%s]\nStackTrace:[%s]",exception.GetCode(),exception.GetMessage(),exception.GetStackTrace())
		if this.Monitor != nil{
			this.Monitor.AscErrorCount()
		}
		result = []reflect.Value{reflect.ValueOf(exception)}
	})
	result = method.Call(arguments)
	if len(result) == 0 {
		result = []reflect.Value{reflect.Zero( reflect.TypeOf(this) )}
	}
	return
}

func (this *BeegoValidateController)AutoRouteMethod(){
	defer language.CatchCrash(func(exception language.Exception){
		this.Log.Critical("Buiness Crash Code:[%d] Message:[%s]\nStackTrace:[%s]",exception.GetCode(),exception.GetMessage(),exception.GetStackTrace())
		if this.Monitor != nil{
			this.Monitor.AscCriticalCount()
		}
		this.Ctx.Output.SetStatus(500)
		this.Ctx.WriteString("server internal error!")
	})
	//查找路由
	appController := this.AppController
	url := this.Ctx.Input.Url()
	urlArray := strings.Split(url,"/")
	if len(urlArray) == 0{
		panic("unknown url segement"+url)
	}
	lastUrlSegment := urlArray[ 2 ]

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

func (this *BeegoValidateController)Check(requireStruct interface{}){	
	//获取require字段
	requireStructType := reflect.TypeOf(requireStruct).Elem()
	requireStructValue := reflect.ValueOf(requireStruct).Elem()
	for i := 0 ; i != requireStructType.NumField() ; i++{
		singleRequireStruct := requireStructType.Field(i)
		singleRequireStructName := singleRequireStruct.Name
		if isPublic(singleRequireStructName) == false {
			continue
		}

		filedName := singleRequireStruct.Tag.Get("validate")
		if filedName == ""{
			filedName = firstLowerName(singleRequireStructName)
		}

		result := this.Ctx.Input.Query(filedName)
		if result == ""{
			continue
		}
		result,err := url.QueryUnescape(result)
		if err != nil{
			language.Throw(1,"参数"+singleRequireStructName+"解析失败，其值为：["+result+"]")
		}

		singleRequireStructValue := requireStructValue.Field(i)
		singleRequireStructValueType := singleRequireStructValue.Type()
		if singleRequireStructValueType == reflect.TypeOf(""){
			singleRequireStructValue.SetString(result)
		}else if singleRequireStructValueType == reflect.TypeOf(1){
			var resultInt int
			resultInt,err := strconv.Atoi(result)
			if err != nil{
				language.Throw(1,"参数"+singleRequireStructName+"不是合法的整数，其值为：["+result+"]")
			}
			singleRequireStructValue.SetInt( int64(resultInt) )
		}else if singleRequireStructValueType == reflect.TypeOf(time.Time{}){
			var resultTime time.Time
			resultTime,err := time.ParseInLocation("2006-01-02 15:04:05",result,time.Now().Local().Location())
			if err != nil{
				language.Throw(1,"参数"+singleRequireStructName+"不是合法的时间，其值为：["+result+"]")
			}
			singleRequireStructValue.Set( reflect.ValueOf(resultTime) )
		}else{
			language.Throw(1,"不合法的参数类型： "+singleRequireStructValueType.String())
		}
	}
}

func (this *BeegoValidateController)CheckGet(requireStruct interface{}){
	if this.Ctx.Input.Method() != "GET"{
		language.Throw(1,"请求Method不是Get方法")
	}
	this.Check(requireStruct)
}

func (this *BeegoValidateController)CheckPost(requireStruct interface{}){
	if this.Ctx.Input.Method() != "POST"{
		language.Throw(1,"请求Method不是POST方法")
	}
	this.Check(requireStruct)
}

func (this *BeegoValidateController)Write(data []byte){
	writer := this.Ctx.ResponseWriter
	writer.Write(data)
}

func (this *BeegoValidateController)WriteMimeHeader(mime string,title string){
	writer := this.Ctx.ResponseWriter
	writerHeader := writer.Header()
	if mime == "json"{
		writerHeader.Set("Content-Type","application/x-javascript; charset=utf-8");
	}else if mime == "javascript"{
		writerHeader.Set("Content-Type","application/x-javascript; charset=utf-8"); 
	}else if mime == "plain"{
		writerHeader.Set("Content-Type","text/plain; charset=utf-8")
	}else if mime == "xlsx"{
		writerHeader.Set("Content-Type","application/vnd.openxmlformats-officedocument; charset=UTF-8"); 
		writerHeader.Set("Pragma","public"); 
		writerHeader.Set("Expires","0"); 
		writerHeader.Set("Cache-Control","must-revalidate, post-check=0, pre-check=0"); 
		writerHeader.Set("Content-Type","application/force-download"); 
		writerHeader.Set("Content-Type","application/octet-stream"); 
		writerHeader.Set("Content-Type","application/download"); 
		writerHeader.Set("Content-Disposition","attachment;filename="+title+".xlsx"); 
		writerHeader.Set("Content-Transfer-Encoding","binary");
	}else if mime == "csv"{
		writerHeader.Set("Content-Type","application/vnd.ms-excel; charset=UTF-8"); 
		writerHeader.Set("Pragma","public"); 
		writerHeader.Set("Expires","0"); 
		writerHeader.Set("Cache-Control","must-revalidate, post-check=0, pre-check=0"); 
		writerHeader.Set("Content-Type","application/force-download"); 
		writerHeader.Set("Content-Type","application/octet-stream"); 
		writerHeader.Set("Content-Type","application/download"); 
		writerHeader.Set("Content-Disposition","attachment;filename="+title+".csv"); 
		writerHeader.Set("Content-Transfer-Encoding","binary");
	}else{
		panic("invalid mime ["+mime+"]")
	}
}

type methodInfo struct{
	name string
	viewName string
	method reflect.Method
}

func getControllerInfoInner(target reflect.Type)(beegoValidateControllerInfo){
	result := beegoValidateControllerInfo{
		-1,
	}
	targetType := reflect.TypeOf(target)
	for i := 0 ; i != targetType.NumMethod() ; i++{
		singleMethod := target.Method(i)
		if singleMethod.Name == "GetBasic"{
			result.getBasicFunc = i
		}
	}
	return result
}

func getControllerInfo(target reflect.Type)(beegoValidateControllerInfo){
	beegoValidateControllerInfoMap.mutex.RLock()
	singleMethodInfo,ok := beegoValidateControllerInfoMap.data[target]
	beegoValidateControllerInfoMap.mutex.RUnlock()

	if ok{
		return singleMethodInfo
	}

	result := getControllerInfoInner(target)

	beegoValidateControllerInfoMap.mutex.Lock()
	beegoValidateControllerInfoMap.data[target] = result
	beegoValidateControllerInfoMap.mutex.Unlock()
	return result
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
		beego.Router(
			namespace+"/"+singleMethodInfo.name+"/*.*",
			target,
			"*:AutoRouteMethod",
		);
		routerControllerMethod[controllerType][singleMethodInfo.name] = *singleMethodInfo
	}
}