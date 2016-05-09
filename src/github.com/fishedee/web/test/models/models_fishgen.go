package test

import (
	. "github.com/fishedee/language"
	. "github.com/fishedee/web"
)

type ClientLoginAoModel interface {
	ModelInterface
	IsLogin() (_fishgen1 bool)
	IsLogin_WithError() (_fishgen1 bool, _fishgenErr Exception)
	Logout()
	Logout_WithError() (_fishgenErr Exception)
	Login(name string, password string) (_fishgen1 bool)
	Login_WithError(name string, password string) (_fishgen1 bool, _fishgenErr Exception)
}

type ClientAoTest interface {
	TestInterface
}

type ConfigAoModel interface {
	ModelInterface
	Set(key string, value string)
	Set_WithError(key string, value string) (_fishgenErr Exception)
	Get(key string) (_fishgen1 string)
	Get_WithError(key string) (_fishgen1 string, _fishgenErr Exception)
	SetStruct(key string, value ConfigData)
	SetStruct_WithError(key string, value ConfigData) (_fishgenErr Exception)
	GetStruct(key string) (_fishgen1 ConfigData)
	GetStruct_WithError(key string) (_fishgen1 ConfigData, _fishgenErr Exception)
}

type ConfigAoTest interface {
	TestInterface
	TestBasic()
	TestStruct()
}

type CounterAoModel interface {
	ModelInterface
	Incr()
	Incr_WithError() (_fishgenErr Exception)
	IncrAtomic()
	IncrAtomic_WithError() (_fishgenErr Exception)
	Reset()
	Reset_WithError() (_fishgenErr Exception)
	Get() (_fishgen1 int)
	Get_WithError() (_fishgen1 int, _fishgenErr Exception)
}

type CounterAoTest interface {
	TestInterface
	TestBasic()
}

type BaseAoModel interface {
	ModelInterface
}

type ExtendAoModel interface {
	ModelInterface
}

type ExtendAoTest interface {
	TestInterface
	TestBasic()
}

type InnerTest interface {
	TestInterface
	TestBasic()
}

func (this *clientLoginAoModel) IsLogin_WithError() (_fishgen1 bool, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.IsLogin()
	return
}

func (this *clientLoginAoModel) Logout_WithError() (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Logout()
	return
}

func (this *clientLoginAoModel) Login_WithError(name string, password string) (_fishgen1 bool, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Login(name, password)
	return
}

func (this *configAoModel) Set_WithError(key string, value string) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Set(key, value)
	return
}

func (this *configAoModel) Get_WithError(key string) (_fishgen1 string, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get(key)
	return
}

func (this *configAoModel) SetStruct_WithError(key string, value ConfigData) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.SetStruct(key, value)
	return
}

func (this *configAoModel) GetStruct_WithError(key string) (_fishgen1 ConfigData, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.GetStruct(key)
	return
}

func (this *counterAoModel) Incr_WithError() (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Incr()
	return
}

func (this *counterAoModel) IncrAtomic_WithError() (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.IncrAtomic()
	return
}

func (this *counterAoModel) Reset_WithError() (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Reset()
	return
}

func (this *counterAoModel) Get_WithError() (_fishgen1 int, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get()
	return
}

func init() {
	InitModel(ClientLoginAoModel(&clientLoginAoModel{}))
	InitTest(ClientAoTest(&clientAoTest{}))
	InitModel(ConfigAoModel(&configAoModel{}))
	InitTest(ConfigAoTest(&configAoTest{}))
	InitModel(CounterAoModel(&counterAoModel{}))
	InitTest(CounterAoTest(&counterAoTest{}))
	InitModel(BaseAoModel(&baseAoModel{}))
	InitModel(ExtendAoModel(&extendAoModel{}))
	InitTest(ExtendAoTest(&extendAoTest{}))
	InitTest(InnerTest(&innerTest{}))
}
