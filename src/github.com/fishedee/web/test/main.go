package main

import (
	"github.com/astaxie/beego"
)

//go:generate fishgen ^./models/.*(ao|db)\.go$
func main() {
	beego.Run()
}
