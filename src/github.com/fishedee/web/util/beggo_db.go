package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	. "github.com/fishedee/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DatabaseManagerConfig struct {
	Driver        string
	Host          string
	Port          int
	User          string
	Passowrd      string
	Database      string
	Debug         bool
	MaxConnection int
}

type DatabaseManager struct {
	*xorm.Engine
	config DatabaseManagerConfig
}

func (this *DatabaseManager) UpdateBatch(data interface{}, indexCol string) (int64, error) {
	tableName := ""
	updateBatchSqlMap := map[string][]string{}

	//判断输入参数类型
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Slice {
		panic("update batch should be a slice")
	}
	dataElemType := dataType.Elem()
	if dataElemType.Kind() != reflect.Struct {
		panic("update btach element should be a struct")
	}
	sf, ok := dataElemType.FieldByName(indexCol)
	if !ok {
		panic("dataElemFieldType has not filed " + indexCol)
	}
	indexType := ""
	if sf.Type.Kind() == reflect.Int {
		indexType = "int"
	} else if sf.Type.Kind() == reflect.String {
		indexType = "string"
	} else {
		panic("非法类型数据!")
	}

	//遍历结构体--取出各字段的值
	dataValue := reflect.ValueOf(data)
	fmt.Printf("%+v", dataValue)
	dataLen := dataValue.Len()
	for i := 0; i != dataLen; i++ {
		singleDataValue := dataValue.Index(i)
		fmt.Printf("%+v", singleDataValue.Type())
		tableInfo := this.TableInfo(singleDataValue.Interface())
		if tableName == "" {
			tableName = tableInfo.Name
		}
		fmt.Printf("tableInfo:%+v, tableName:%+v", tableInfo, tableName)
		singleDataValueLen := singleDataValue.NumField()
		for j := 0; j != singleDataValueLen; j++ {
			colName := dataElemType.Field(j).Name
			singleDataField := singleDataValue.Field(j)
			singleDataFieldType := singleDataField.Kind()
			colVal := ""
			if singleDataFieldType == reflect.String {
				colVal = singleDataField.String()
			} else if singleDataFieldType == reflect.Int {
				colVal = strconv.Itoa(int(singleDataField.Int()))
			} else {
				panic("结构体中含有非法类型数据！")
			}
			fmt.Printf("%+v, %+v", colName, colVal)
			if strings.ToLower(colName) != strings.ToLower(tableInfo.AutoIncrement) {
				updateBatchSqlMap[colName] = append(updateBatchSqlMap[colName], colVal)
			}
		}
	}
	fmt.Printf("%+v", updateBatchSqlMap)

	//拼接sql语句
	updateBatchSql := " update " + tableName + " set "
	sum := 0
	for k, v := range updateBatchSqlMap {
		if k != indexCol {
			otherColType := ""
			sf, ok = dataElemType.FieldByName(k)
			if !ok {
				panic("dataElemFieldType has not filed " + k)
			}
			if sf.Type.Kind() == reflect.Int {
				otherColType = "int"
			} else if sf.Type.Kind() == reflect.String {
				otherColType = "string"
			} else {
				panic("非法类型数据!")
			}
			sum++
			updateBatchSql += " " + k + " = case " + indexCol
			for n := 0; n < len(v); n++ {
				//根据数据类型判断是否添加单引号
				whenSql := ""
				if indexType == "int" {
					whenSql = " when " + updateBatchSqlMap[indexCol][n] + " "
				} else if indexType == "string" {
					whenSql = " when '" + updateBatchSqlMap[indexCol][n] + "' "
				} else {
					panic("非法类型数据!")
				}
				thenSql := ""
				if otherColType == "int" {
					thenSql = " then " + v[n] + " "
				} else if otherColType == "string" {
					thenSql = " then '" + v[n] + "' "
				} else {
					panic("非法类型数据!")
				}
				updateBatchSql += whenSql + thenSql
			}
			updateBatchSql += " end "
			if sum < len(updateBatchSqlMap)-1 {
				updateBatchSql += ", "
			}
		}
	}
	updateBatchSql += " where " + indexCol + " in ("
	for singleKey, singleValue := range updateBatchSqlMap[indexCol] {
		if indexType == "int" {
			updateBatchSql += singleValue
		} else if indexType == "string" {
			updateBatchSql += " '" + singleValue + "' "
		}
		if singleKey < len(updateBatchSqlMap[indexCol])-1 {
			updateBatchSql += ","
		}
	}
	updateBatchSql += ")"
	fmt.Printf("%+v", updateBatchSql)

	//执行sql
	res, err := this.Exec(updateBatchSql)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

var newDatabaseManagerMemory *MemoryFunc
var newDatabaseManagerFromConfigMemory *MemoryFunc

func init() {
	var err error
	newDatabaseManagerMemory, err = NewMemoryFunc(newDatabaseManager, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
	newDatabaseManagerFromConfigMemory, err = NewMemoryFunc(newDatabaseManagerFromConfig, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
}

func newDatabaseManager(config DatabaseManagerConfig) (*DatabaseManager, error) {
	if config.Driver == "" {
		return nil, nil
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
		tempDb.ShowSQL(true)
	}
	if config.MaxConnection > 0 {
		tempDb.SetMaxOpenConns(config.MaxConnection)
	}
	tempDb.Ping()
	return &DatabaseManager{
		Engine: tempDb,
		config: config,
	}, nil
}

func NewDatabaseManager(config DatabaseManagerConfig) (*DatabaseManager, error) {
	result, err := newDatabaseManagerMemory.Call(config)
	return result.(*DatabaseManager), err
}

func newDatabaseManagerFromConfig(configName string) (*DatabaseManager, error) {
	dbdirver := beego.AppConfig.String(configName + "dirver")
	dbhost := beego.AppConfig.String(configName + "host")
	dbport := beego.AppConfig.String(configName + "port")
	dbuser := beego.AppConfig.String(configName + "user")
	dbpassword := beego.AppConfig.String(configName + "password")
	dbdatabase := beego.AppConfig.String(configName + "database")
	dbmaxconnection := beego.AppConfig.String(configName + "maxconnection")
	dbdebug := beego.AppConfig.String(configName + "debug")

	config := DatabaseManagerConfig{}
	config.Driver = dbdirver
	config.Host = dbhost
	config.Port, _ = strconv.Atoi(dbport)
	config.User = dbuser
	config.Passowrd = dbpassword
	config.Database = dbdatabase
	config.Debug, _ = strconv.ParseBool(dbdebug)
	config.MaxConnection, _ = strconv.Atoi(dbmaxconnection)

	return NewDatabaseManager(config)
}

func NewDatabaseManagerFromConfig(configName string) (*DatabaseManager, error) {
	result, err := newDatabaseManagerFromConfigMemory.Call(configName)
	return result.(*DatabaseManager), err
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
