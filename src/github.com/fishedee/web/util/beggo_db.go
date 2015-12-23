package util

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	. "github.com/fishedee/util"
	"strconv"
	"strings"
)

type DatabaseManagerConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Passowrd string
	Database string
	Debug    bool
}

type DatabaseManager struct {
	*xorm.Engine
	config DatabaseManagerConfig
}

var newDatabaseManagerMemory *MemoryFunc
var newDatabaseManagerFromConfigMemory *MemoryFunc

func init(){
	var err error
	newDatabaseManagerMemory,err = NewMemoryFunc(newDatabaseManager,MemoryFuncCacheNormal)
	if err != nil{
		panic(err)
	}
	newDatabaseManagerFromConfigMemory,err = NewMemoryFunc(newDatabaseManagerFromConfig,MemoryFuncCacheNormal)
	if err != nil{
		panic(err)
	}
}

func newDatabaseManager(config DatabaseManagerConfig) (*DatabaseManager, error) {
	if config.Driver == ""{
		return nil,nil
	}
	dblink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		config.User,
		config.Passowrd,
		config.Host,
		config.Port,
		config.Database,
	)
	tempDb, err := xorm.NewEngine(config.Driver, dblink)
	if err != nil {
		return nil, err
	}

	tempDb.SetTableMapper(&tableMapper{})
	tempDb.SetColumnMapper(&columnMapper{})
	if config.Debug {
		tempDb.ShowSQL = true
	}
	return &DatabaseManager{
		Engine: tempDb,
		config: config,
	}, nil
}

func NewDatabaseManager(config DatabaseManagerConfig)(*DatabaseManager,error){
	result,err := newDatabaseManagerMemory.Call(config)
	return result.(*DatabaseManager),err
}

func newDatabaseManagerFromConfig(configName string) (*DatabaseManager, error) {
	dbdirver := beego.AppConfig.String(configName + "dirver")
	dbhost := beego.AppConfig.String(configName + "host")
	dbport := beego.AppConfig.String(configName + "port")
	dbuser := beego.AppConfig.String(configName + "user")
	dbpassword := beego.AppConfig.String(configName + "password")
	dbdatabase := beego.AppConfig.String(configName + "database")
	dbdebug := beego.AppConfig.String(configName + "debug")

	config := DatabaseManagerConfig{}
	config.Driver = dbdirver
	config.Host = dbhost
	config.Port, _ = strconv.Atoi(dbport)
	config.User = dbuser
	config.Passowrd = dbpassword
	config.Database = dbdatabase
	config.Debug, _ = strconv.ParseBool(dbdebug)

	return NewDatabaseManager(config)
}

func NewDatabaseManagerFromConfig(configName string)(*DatabaseManager,error){
	result,err := newDatabaseManagerFromConfigMemory.Call(configName)
	return result.(*DatabaseManager),err
}

type tableMapper struct {
}

func (this *tableMapper) Obj2Table(name string) string {
	result := []rune{}
	result = append(result, 't')
	for _, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			result = append(result, '_')
			chr -= ('A' - 'a')
		}
		result = append(result, chr)
	}
	return string(result)
}

func (this *tableMapper) Table2Obj(in string) string {
	fmt.Println("Obj2Table2 " + in)
	return in
}

type columnMapper struct {
}

func (this *columnMapper) Obj2Table(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func (this *columnMapper) Table2Obj(in string) string {
	fmt.Println("Obj2Table4 " + in)
	return in
}
