package test

import (
	. "github.com/fishedee/language"
	. "github.com/fishedee/web"
)

type ClientLoginAoModel interface {
	IsLogin() (_fishgen1 bool)
	IsLogin_WithError() (_fishgen1 bool, _fishgenErr Exception)
	Logout()
	Logout_WithError() (_fishgenErr Exception)
	Login(name string, password string) (_fishgen1 bool)
	Login_WithError(name string, password string) (_fishgen1 bool, _fishgenErr Exception)
}

type ClientAoTest interface {
	TestBasic()
}

type ConfigAoModel interface {
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
	TestBasic()
	TestStruct()
}

type CounterAoModel interface {
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
	TestBasic()
}

type InnerTest interface {
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
	v0 := ClientLoginAoModel(&clientLoginAoModel{})
	InitModel(&v0)
	v1 := ClientAoTest(&clientAoTest{})
	InitTest(&v1)
	v2 := ConfigAoModel(&configAoModel{})
	InitModel(&v2)
	v3 := ConfigAoTest(&configAoTest{})
	InitTest(&v3)
	v4 := CounterAoModel(&counterAoModel{})
	InitModel(&v4)
	v5 := CounterAoTest(&counterAoTest{})
	InitTest(&v5)
	v6 := InnerTest(&innerTest{})
	InitTest(&v6)
}
