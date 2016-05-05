package number

import (
	. "github.com/fishedee/web"
)

type NumberAoModel struct {
	Model
}

func (this *NumberAoModel) Add(left int, right int) int {
	return left + right
}
