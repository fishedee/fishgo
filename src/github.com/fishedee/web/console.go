package web

import (
	"reflect"
)

func InitModel(handler interface{}) {
	target := reflect.ValueOf(handler)
	basic := initEmptyBasic(nil)
	injectIoc(target, basic)
}

func InitController(handler interface{}) {
	InitModel(handler)
}
