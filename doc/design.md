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

* Run()error，正常运行时阻塞进程，正常退出时返回nil，异常发生时马上退出，返回非nil。要注意，Run的返回意味着必要的收尾工作都已经做完了，返回后就能保证重新start了。
* Close()，关闭进程。要注意的是，这个函数的返回更加严格，需要所有的收尾的工作都已经做完，由run启动的所有go进程都必须全部收尾完成。这样做的原因是为了保证多个daemon的close时的依赖顺序能按照需要的执行。
* Close时不关闭非Daemon接口，因为在Daemon关闭了以后，外部的timer或者http请求可能还会调用非Daemon接口，例如queue中的Produce。

## 4.3 error与exception

除了language包，其他模块的接口都不允许在实现里面直接panic，也就是接口必须能支持返回error来检查出错的。另外，对于那些经常需要检查error返回值的业务而言，可以提供额外的Must接口，以提高使用的体验。这样做可以同时兼顾返回error的灵活性，和使用exception的便利性。

# 5 使用

## 5.1 xorm

```
func (this *ContactDb) getWhere(where Contact) (string, []interface{}) {
	whereSql := []string{}
	argSql := []interface{}{}
	if where.Name != "" {
		whereSql = append(whereSql, "name like ?")
		argSql = append(argSql, "%"+where.Name+"%")
	}
	if where.Remark != "" {
		whereSql = append(whereSql, "remark like ?")
		argSql = append(argSql, "%"+where.Remark+"%")
	}
	if where.GroupId != 0 {
		whereSql = append(whereSql, "groupId = ?")
		argSql = append(argSql, strconv.Itoa(where.GroupId))
	}
	if where.Phone != "" {
		whereSql = append(whereSql, "(phone like ? or phone2 like ? or phone3 like ?)")
		argSql = append(argSql, "%"+where.Phone+"%", "%"+where.Phone+"%", "%"+where.Phone+"%")
	}
	return Implode(whereSql, " and "), argSql
}

func (this *ContactDb) Search(where Contact, limit CommonPage) Contacts {
	db := this.db.NewSession()

	whereSql, whereArg := this.getWhere(where)

	if whereSql != "" {
		db.Where(whereSql, whereArg...)
	}
	result := Contacts{}
	db.OrderBy("createTime desc").Limit(limit.PageSize, limit.PageIndex).MustFind(&result.Data)

	if whereSql != "" {
		db.Where(whereSql, whereArg...)
	}
	result.Count = int(db.MustCount(&Contact{}))

	return result
}

func (this *ContactDb) Get(contactId int) Contact {
	var contacts []Contact
	this.db.Where("contactId = ?", contactId).MustFind(&contacts)
	if len(contacts) == 0 {
		Throw(1, "该"+strconv.Itoa(contactId)+"联系人不存在")
	}
	return contacts[0]
}

func (this *ContactDb) Del(contactId int) {
	this.db.Where("contactId = ?", contactId).MustDelete(&Contact{})
}

func (this *ContactDb) Add(contact Contact) {
	this.db.MustInsert(contact)
}

func (this *ContactDb) Mod(contactId int, contact Contact) {
	this.db.Where("contactId = ?", contactId).AllCols().MustUpdate(&contact)
}

func (this *ContactDb) ResetGroupId(groupId int) {
	this.db.Where("groupId = ?", groupId).Cols("groupId").MustUpdate(&Contact{
		GroupId: 0,
	})
}
```

db的标准写法，如上，注意search需要两次的where,mod需要allcols。要注意在，性能要求较高的场合不能使用pageIndex+pageSize，而是nextPage+pageSize。

## 5.2 queue

* 没有queue的topic在发送时会报错
* 有queue但没有consumer的topic发送时不会报错，但会导致queue大量堆积

# 6 历史

## 6.1 20180225

重新设计web的模块，将其重构到app的模块，目的是更加模块化的实现，当然，代价是写代码远远没有原来方便了，普通的web的项目依然建议用成熟的web模块。app模块暂时只作为实验性质，不要使用。

