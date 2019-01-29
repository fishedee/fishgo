#!/bin/sh
cd `pwd`/bin 
go install github.com/beego/bee
go install github.com/fishedee/web/fishgen
go install github.com/fishedee/web/fishcmd
go install github.com/fishedee/app/mock
go install github.com/fishedee/language/querygen
