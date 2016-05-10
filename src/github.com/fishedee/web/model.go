package web

type Model struct {
	*Basic
}

func InitModel(model interface{}) {
	err := addIocTarget(model)
	if err != nil {
		panic(err)
	}
}
