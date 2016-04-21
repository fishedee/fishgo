package number

type NumberAoModel struct {
	BaseModel
}

func (this *NumberAoModel) Add(left int, right int) int {
	return left + right
}
