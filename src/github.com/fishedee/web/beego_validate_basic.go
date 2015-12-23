package web

import (
	. "github.com/fishedee/web/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/context"
)

type BeegoValidateBasic struct {
	ctx *context.Context
	Session *SessionManager
	DB *DatabaseManager
	DB2 *DatabaseManager
	DB3 *DatabaseManager
	DB4 *DatabaseManager
	DB5 *DatabaseManager
	logger *logs.BeeLogger
	Log *LogManager
}

var globalBasic BeegoValidateBasic

func init(){
	var err error
	globalBasic.Session,err = NewSessionManagerFromConfig("fishsession")
	if err != nil{
		panic(err)
	}
	globalBasic.DB,err = NewDatabaseManagerFromConfig("fishdb")
	if err != nil{
		panic(err)
	}
	globalBasic.DB2,err = NewDatabaseManagerFromConfig("fishdb2")
	if err != nil{
		panic(err)
	}
	globalBasic.DB3,err = NewDatabaseManagerFromConfig("fishdb3")
	if err != nil{
		panic(err)
	}
	globalBasic.DB4,err = NewDatabaseManagerFromConfig("fishdb4")
	if err != nil{
		panic(err)
	}
	globalBasic.DB5,err = NewDatabaseManagerFromConfig("fishdb5")
	if err != nil{
		panic(err)
	}
	globalBasic.logger,err = NewLogManagerFromConfig("fishlog")
	if err != nil{
		panic(err)
	}
} 
func NewBeegoValidateBasic(ctx *context.Context)(*BeegoValidateBasic){
	result := globalBasic
	result.ctx = ctx
	result.Log = NewLogManagerWithCtx(ctx,result.logger)
	return &result
}