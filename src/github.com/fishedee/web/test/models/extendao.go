package test

import (
	. "github.com/fishedee/web"
)

type BaseAoModel struct {
	Model
	ConfigAo ConfigAoModel
}

type ExtendAoModel struct {
	BaseAoModel
}
