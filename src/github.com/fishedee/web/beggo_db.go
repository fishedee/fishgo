package web

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
var DB2 *gorm.DB
var DB3 *gorm.DB
var DB4 *gorm.DB
var DB5 *gorm.DB

func initSingleDb(prefix string)(*gorm.DB){
	dbdirver := beego.AppConfig.String("fishdbdirver"+prefix)
	dbhost := beego.AppConfig.String("fishdbhost"+prefix)
	dbport := beego.AppConfig.String("fishdbport"+prefix)
	dbuser := beego.AppConfig.String("fishdbuser"+prefix)
	dbpassword := beego.AppConfig.String("fishdbpassword"+prefix)
	dbdatabase := beego.AppConfig.String("fishdbdatabase"+prefix)

	if dbdirver == ""{
		return nil
	}
	
	dblink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",dbuser,dbpassword,dbhost,dbport,dbdatabase)
	tempDb, err := gorm.Open(dbdirver,dblink)
	if err != nil {
		panic("open mysql error!")
	}

	err = tempDb.DB().Ping()
	if err != nil {
		panic("open ping error!")
	}
	tempDb.DB().SetMaxIdleConns(10000)
	tempDb.DB().SetMaxOpenConns(10000)
	return &tempDb
}

func init() {
	DB = initSingleDb("")
	DB2 = initSingleDb("2")
	DB3 = initSingleDb("3")
	DB4 = initSingleDb("4")
	DB5 = initSingleDb("4")
}
