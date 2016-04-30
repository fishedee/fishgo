package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.com/astaxie/beego"
	. "github.com/fishedee/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
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

type zeroable interface {
	IsZero() bool
}

func (this *DatabaseManager) rValue(bean interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(bean))
}

func (this *DatabaseManager) isZero(k interface{}) bool {
	switch k.(type) {
	case int:
		return k.(int) == 0
	case int8:
		return k.(int8) == 0
	case int16:
		return k.(int16) == 0
	case int32:
		return k.(int32) == 0
	case int64:
		return k.(int64) == 0
	case uint:
		return k.(uint) == 0
	case uint8:
		return k.(uint8) == 0
	case uint16:
		return k.(uint16) == 0
	case uint32:
		return k.(uint32) == 0
	case uint64:
		return k.(uint64) == 0
	case float32:
		return k.(float32) == 0
	case float64:
		return k.(float64) == 0
	case bool:
		return k.(bool) == false
	case string:
		return k.(string) == ""
	case zeroable:
		return k.(zeroable).IsZero()
	}
	return false
}

func (this *DatabaseManager) value2Interface(fieldValue reflect.Value) (interface{}, error) {
	fieldType := fieldValue.Type()
	fieldTypeKind := fieldType.Kind()
	switch fieldTypeKind {
	case reflect.Bool:
		return fieldValue.Bool(), nil
	case reflect.String:
		return fieldValue.String(), nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return fieldValue.Int(), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return fieldValue.Uint(), nil
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			t := fieldValue.Interface().(time.Time)
			tf := t.Format("2006-01-02 15:04:05")
			return tf, nil
		} else {
			return nil, fmt.Errorf("Unsupported type %v", fieldType)
		}
	default:
		return nil, fmt.Errorf("Unsupported type %v", fieldType)
	}
}

type tableName interface {
	TableName() string
}

func (this *DatabaseManager) autoMapType(v reflect.Value) *core.Table {
	t := v.Type()
	table := core.NewEmptyTable()
	if tb, ok := v.Interface().(tableName); ok {
		table.Name = tb.TableName()
	} else {
		if v.CanAddr() {
			if tb, ok = v.Addr().Interface().(tableName); ok {
				table.Name = tb.TableName()
			}
		}
		if table.Name == "" {
			table.Name = this.TableMapper.Obj2Table(t.Name())
		}
	}
	table.Type = t
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		ormTagStr := tag.Get("xorm")
		if ormTagStr == "-" || ormTagStr == "<-" {
			continue
		}
		col := &core.Column{FieldName: t.Field(i).Name, Nullable: true, IsPrimaryKey: false,
			IsAutoIncrement: false, MapType: core.TWOSIDES, Indexes: make(map[string]bool)}
		col.Name = this.ColumnMapper.Obj2Table(t.Field(i).Name)
		table.AddColumn(col)
	}
	return table
}

func (this *DatabaseManager) UpdateBatch(rowsSlicePtr interface{}, indexColName string) (int64, error) {
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	if sliceValue.Kind() != reflect.Slice {
		return 0, errors.New("needs a pointer to a slice")
	}
	if sliceValue.Len() == 0 {
		return 0, errors.New("update rows is empty")
	}

	bean := sliceValue.Index(0).Interface()
	elementValue := this.rValue(bean)
	table := this.autoMapType(elementValue)
	size := sliceValue.Len()

	var rows = make([][]interface{}, 0)
	var indexRow = make([]interface{}, 0)
	cols := make([]*core.Column, 0)
	var indexCol *core.Column

	//提取字段
	for i := 0; i < size; i++ {
		v := sliceValue.Index(i)
		vv := reflect.Indirect(v)

		//处理需要的update的列
		if i == 0 {
			for _, col := range table.Columns() {
				ptrFieldValue, err := col.ValueOfV(&vv)
				if err != nil {
					return 0, err
				}
				fieldValue := *ptrFieldValue
				if this.isZero(fieldValue.Interface()) {
					continue
				}
				if col.Name == indexColName {
					indexCol = col
				} else {
					cols = append(cols, col)
				}
			}
			if indexCol == nil {
				return 0, errors.New("counld not found index col " + indexColName)
			}
		}

		//处理需要的update的值
		var singleRow = make([]interface{}, 0)
		for _, col := range cols {
			ptrFieldValue, err := col.ValueOfV(&vv)
			if err != nil {
				return 0, err
			}
			fieldValue := *ptrFieldValue
			arg, err := this.value2Interface(fieldValue)
			if err != nil {
				return 0, err
			}
			singleRow = append(singleRow, arg)
		}
		rows = append(rows, singleRow)
		ptrFieldValue, err := indexCol.ValueOfV(&vv)
		if err != nil {
			return 0, err
		}
		fieldValue := *ptrFieldValue
		arg, err := this.value2Interface(fieldValue)
		if err != nil {
			return 0, err
		}
		indexRow = append(indexRow, arg)
	}
	if len(cols) == 0 {
		return 0, errors.New("update cols is empty! " + fmt.Sprintf("%v", rowsSlicePtr))
	}

	//拼接sql
	var sqlArgs = make([]interface{}, 0)
	var sql = "UPDATE " + table.Name + " SET "
	for colIndex, col := range cols {
		if colIndex != 0 {
			sql += " , "
		}
		sql += this.Engine.QuoteStr() + col.Name + this.Engine.QuoteStr()
		sql += " = CASE "
		sql += this.Engine.QuoteStr() + indexCol.Name + this.Engine.QuoteStr()
		for rowIndex, row := range rows {
			sql += " WHEN ? THEN ? "
			sqlArgs = append(sqlArgs, indexRow[rowIndex])
			sqlArgs = append(sqlArgs, row[colIndex])
		}
		sql += " END "
	}
	sql += " WHERE " + this.Engine.QuoteStr() + indexCol.Name + this.Engine.QuoteStr() + " IN ( "
	for rowIndex, row := range indexRow {
		if rowIndex != 0 {
			sql += " , "
		}
		sql += " ? "
		sqlArgs = append(sqlArgs, row)
	}
	sql += " ) "

	//执行sql
	res, err := this.Exec(sql, sqlArgs...)
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
