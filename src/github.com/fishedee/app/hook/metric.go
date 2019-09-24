package hook

import (
	"context"
	"fmt"
	. "github.com/fishedee/app/metric"
	. "github.com/fishedee/app/proxy"
	. "github.com/fishedee/encoding"
	"reflect"
	"sync"
	"time"
)

func NewModuleMetricHook(metric Metric) ProxyHook {
	contextType := reflect.TypeOf([]context.Context{}).Elem()
	return func(ctx ProxyContext, origin reflect.Value) reflect.Value {
		originFirstArgType := origin.Type().In(0)
		if originFirstArgType.ConvertibleTo(contextType) == false {
			panic(fmt.Sprintf("%v first argument can not convert to context.Context", origin.Type()))
		}
		ctxName := ctx.InterfaceName + "." + ctx.MethodName
		ctxNameEncode, err := EncodeUrl(ctxName)
		if err != nil {
			panic(err)
		}
		moduleTag := "module=" + ctxNameEncode
		rwMutex := &sync.RWMutex{}
		pathMapper := map[string]MetricTimer{}
		getTimer := func(name string) MetricTimer {
			rwMutex.RLock()
			result, isExist := pathMapper[name]
			rwMutex.RUnlock()

			if isExist {
				return result
			}

			result = metric.GetTimer(name)

			rwMutex.Lock()
			pathMapper[name] = result
			rwMutex.Unlock()

			return result

		}
		return reflect.MakeFunc(origin.Type(), func(args []reflect.Value) []reflect.Value {
			//获取path参数
			context := args[0].Convert(contextType).Interface().(context.Context)
			path := context.Value("path").(string)
			pathEncode, err := EncodeUrl(path)
			if err != nil {
				panic(err)
			}
			pathTag := "path=" + pathEncode

			//获取metric
			metricName := "module.request?" + pathTag + "&" + moduleTag
			metricTimer := getTimer(metricName)

			//调用接口
			begin := time.Now()
			result := origin.Call(args)
			end := time.Now()
			metricTimer.Update(end.Sub(begin))
			return result
		})
	}
}
