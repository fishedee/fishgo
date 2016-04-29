package number

import (
	. "github.com/fishedee/language"
)

func (this *NumberAoModel) Add_WithError(left int, right int) (_fishgen1 int, _fishgenErr Exception) {
	defer Catch(func(exception Exception) {
		_fishgenErr = exception
	})
	_fishgen1 = this.Add(left, right)
	return
}
