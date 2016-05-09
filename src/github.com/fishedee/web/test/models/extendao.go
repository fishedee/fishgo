package test

import (
	. "github.com/fishedee/web"
)

type baseAoModel struct {
	Model
	configAo ConfigAoModel
}

type extendAoModel struct {
	baseAoModel
}
