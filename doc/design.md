# 1 概述

fishgo设计原则与历史

# 2 目标

简要的公共库

# 3 原则

* 模块化，每个模块只做一个事情，能随意组合模块实现需要的结果
* 无侵入性，避免一个模块的使用需要侵入性引入其他框架和模块，每个模块都应该是能独立使用的
* 无偏好性，例如，有些人喜欢exception，有些人喜欢error，接口中已经支持这两种做法，而不是只偏向一种。
* 灵活，提供简单而且扩展性大的接口，恰当地引入中间件机制，装饰器机制，依赖注入等等
* 一致，名称一致，start对应stop，connect对应disconnect，run对应close，不要乱配对
* 兼容，已经发布的接口不能改变语义，提供新语义要么新建模块，要么新建函数
* 稳定，必须有完整和强健的单元测试

# 4 约定

## 4.1 点号引入

一般模块都是点号引入的，注意模块的名字要用全名，不用New，而用NewXXXX

## 4.2 Daemon生命周期

Daemon在使用时需要考虑到四个场景

* 单个daemon运行，一个daemon的关闭触发整个进程的关闭
* 多个daemon运行，所有daemon的关闭才能触发整个进程的关闭
* 异常退出的情况下，一个daemon的退出会触发其他daemon的关闭
* 主动退出的情况下，逐个让每个daemon优雅关闭

```
func WhenExit(){
	daemon1.Close()
}
func main(){
	err := daemon1.Run()
	if err != nil{
		panic(err)
	}
}
```

单个daemon的运行，异常退出时panic触发整个进程的退出，正常退出时Close触发Run退出

```
func WhenExit(){
	workgroup.Close()
}

func main(){
	workgroup := NewWorkGroup()
	workgroup.Add(queue)
	workgroup.Add(timer)
	workgroup.Add(server)
	err := workgroup.Run()
	if err != nil{
		panic(err)
	}
}
```

多个daemon的运行，异常退出时panic触发整个进程的退出，正常退出时Close触发Run退出

所以，所有的daemon都需要有以下的两个接口

* Run()error，正常运行时阻塞进程，正常退出时返回nil，异常发生时马上退出，返回非nil。要注意，Run的返回意味着所有的收尾工作都已经做完了。
* Close()，立即关闭进程，这个操作应该是马上返回的，不需要等待收尾工作的完成。

## 4.3 error与exception

除了language包，其他模块的接口都不允许在实现里面直接panic，也就是接口必须能支持返回error来检查出错的。另外，对于那些经常需要检查error返回值的业务而言，可以提供额外的Must接口，以提高使用的体验。这样做可以同时兼顾返回error的灵活性，和使用exception的便利性。

# 5 历史

## 5.1 20180225

重新设计web的模块，将其重构到app的模块，目的是更加模块化的实现，当然，代价是写代码远远没有原来方便了，普通的web的项目依然建议用成熟的web模块。app模块暂时只作为实验性质，不要使用。

