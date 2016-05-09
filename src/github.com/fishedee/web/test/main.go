package main

import (
	"github.com/astaxie/beego"
)

//go:generate fishgen -force ^.*(ao|_testing|controller)\.go$
func main() {
	beego.Run()
}
