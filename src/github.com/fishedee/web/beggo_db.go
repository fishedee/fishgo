package web

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var DB *xorm.Engine
var DB2 *xorm.Engine
var DB3 *xorm.Engine
var DB4 *xorm.Engine
var DB5 *xorm.Engine

type tableMapper struct{

}

func (this *tableMapper)Obj2Table(name string)string{
	result := []rune{}
	result = append(result,'t')
	for _,chr := range name{
		if isUpper := 'A' <= chr && chr <= 'Z' ;isUpper{
			result = append(result,'_')
			chr -= ('A'-'a')
		}
		result = append(result,chr)
	}
	return string(result)
}

func (this *tableMapper)Table2Obj(in string)string{
	fmt.Println("Obj2Table2 "+in)
	return in
}

type columnMapper struct{

}

func (this *columnMapper)Obj2Table(name string)string{
	return strings.ToLower(name[0:1])+name[1:]
}

func (this *columnMapper)Table2Obj(in string)string{
	fmt.Println("Obj2Table4 "+in)
	return in
}

func initSingleDb(prefix string)(*xorm.Engine){
	dbdirver := beego.AppConfig.String("fishdbdirver"+prefix)
	dbhost := beego.AppConfig.String("fishdbhost"+prefix)
	dbport := beego.AppConfig.String("fishdbport"+prefix)
	dbuser := beego.AppConfig.String("fishdbuser"+prefix)
	dbpassword := beego.AppConfig.String("fishdbpassword"+prefix)
	dbdatabase := beego.AppConfig.String("fishdbdatabase"+prefix)
	dbdebug := beego.AppConfig.String("fishdbdebug"+prefix)
	dbdebugBool,_ := strconv.ParseBool(dbdebug)

	if dbdirver == ""{
		return nil
	}
	
	dblink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",dbuser,dbpassword,dbhost,dbport,dbdatabase)
	tempDb,err := xorm.NewEngine(dbdirver,dblink)
	if err != nil {
		panic("open mysql error! "+err.Error()+","+dblink)
	}

	tempDb.SetTableMapper(&tableMapper{})
	tempDb.SetColumnMapper(&columnMapper{})
	if dbdebugBool{
		tempDb.ShowSQL = true
	}
	return tempDb
}

func init() {
	DB = initSingleDb("")
	DB2 = initSingleDb("2")
	DB3 = initSingleDb("3")
	DB4 = initSingleDb("4")
	DB5 = initSingleDb("5")
}
