package number

import (
	. "github.com/fishedee/language"
	. "github.com/fishedee/web"
)

type NumberAoModel interface {
	Add(left int, right int) (_fishgen1 int)
	Add_WithError(left int, right int) (_fishgen1 int, _fishgenErr Exception)
}

type NumberAoTest interface {
	TestBasic()
}

func (this *numberAoModel) Add_WithError(left int, right int) (_fishgen1 int, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Add(left, right)
	return
}

func init() {
	InitModel(NumberAoModel(&numberAoModel{}))
	InitTest(NumberAoTest(&numberAoTest{}))
}
