package number

import (
	. "github.com/fishedee/web"
)

type numberAoModel struct {
	Model
}

func (this *numberAoModel) Add(left int, right int) int {
	return left + right
}
