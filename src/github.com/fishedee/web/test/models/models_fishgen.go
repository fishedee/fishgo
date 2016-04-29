package test

import (
	. "github.com/fishedee/language"
)

func (this *ClientLoginAoModel) IsLogin_WithError() (_fishgen1 bool, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.IsLogin()
	return
}

func (this *ClientLoginAoModel) Logout_WithError() (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Logout()
	return
}

func (this *ClientLoginAoModel) Login_WithError(name string, password string) (_fishgen1 bool, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Login(name, password)
	return
}

func (this *ConfigAoModel) Set_WithError(key string, value string) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Set(key, value)
	return
}

func (this *ConfigAoModel) Get_WithError(key string) (_fishgen1 string, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get(key)
	return
}

func (this *ConfigAoModel) SetStruct_WithError(key string, value ConfigData) (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.SetStruct(key, value)
	return
}

func (this *ConfigAoModel) GetStruct_WithError(key string) (_fishgen1 ConfigData, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.GetStruct(key)
	return
}

func (this *CounterAoModel) Incr_WithError() (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Incr()
	return
}

func (this *CounterAoModel) IncrAtomic_WithError() (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.IncrAtomic()
	return
}

func (this *CounterAoModel) Reset_WithError() (_fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	this.Reset()
	return
}

func (this *CounterAoModel) Get_WithError() (_fishgen1 int, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Get()
	return
}
