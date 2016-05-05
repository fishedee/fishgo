package web

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type DatabaseCommon interface {
	Close() error
	Sql(querystring string, args ...interface{}) DatabaseSession
	NoAutoTime() DatabaseSession
	NoAutoCondition(no ...bool) DatabaseSession
	Cascade(trueOrFalse ...bool) DatabaseSession
	Where(querystring string, args ...interface{}) DatabaseSession
	Id(id interface{}) DatabaseSession
	Distinct(columns ...string) DatabaseSession
	Select(str string) DatabaseSession
	Cols(columns ...string) DatabaseSession
	AllCols() DatabaseSession
	MustCols(columns ...string) DatabaseSession
	UseBool(columns ...string) DatabaseSession
	Omit(columns ...string) DatabaseSession
	Nullable(columns ...string) DatabaseSession
	In(column string, args ...interface{}) DatabaseSession
	Incr(column string, arg ...interface{}) DatabaseSession
	Decr(column string, arg ...interface{}) DatabaseSession
	SetExpr(column string, expression string) DatabaseSession
	Table(tableNameOrBean interface{}) DatabaseSession
	Alias(alias string) DatabaseSession
	Limit(limit int, start ...int) DatabaseSession
	Desc(colNames ...string) DatabaseSession
	Asc(colNames ...string) DatabaseSession
	OrderBy(order string) DatabaseSession
	Join(join_operator string, tablename interface{}, condition string, args ...interface{}) DatabaseSession
	GroupBy(keys string) DatabaseSession
	Having(conditions string) DatabaseSession
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Query(sql string, paramStr ...interface{}) (resultsSlice []map[string][]byte, err error)
	Insert(beans ...interface{}) (int64, error)
	InsertOne(bean interface{}) (int64, error)
	Update(bean interface{}, condiBeans ...interface{}) (int64, error)
	Delete(bean interface{}) (int64, error)
	Get(bean interface{}) (bool, error)
	Find(beans interface{}, condiBeans ...interface{}) error
	Count(bean interface{}) (int64, error)
}

type DatabaseSession interface {
	DatabaseCommon
	And(querystring string, args ...interface{}) DatabaseSession
	Or(querystring string, args ...interface{}) DatabaseSession
	ForUpdate() DatabaseSession
	Begin() error
	Commit() error
}

type Database interface {
	DatabaseCommon
	NewSession() DatabaseSession
	UpdateBatch(rowsSlicePtr interface{}, indexColName string) (int64, error)
}

type DatabaseConfig struct {
	Driver        string
	Host          string
	Port          int
	User          string
	Passowrd      string
	Database      string
	Debug         bool
	MaxConnection int
}

type databaseImplement struct {
	*xorm.Engine
	config DatabaseConfig
}

type databaseSessionImplement struct {
	*xorm.Session
}

func NewDatabase(config DatabaseConfig) (Database, error) {
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
	return &databaseImplement{
		Engine: tempDb,
		config: config,
	}, nil
}

func NewDatabaseFromConfig(configName string) (Database, error) {
	dbdirver := globalBasic.Config.GetString(configName + "dirver")
	dbhost := globalBasic.Config.GetString(configName + "host")
	dbport := globalBasic.Config.GetString(configName + "port")
	dbuser := globalBasic.Config.GetString(configName + "user")
	dbpassword := globalBasic.Config.GetString(configName + "password")
	dbdatabase := globalBasic.Config.GetString(configName + "database")
	dbmaxconnection := globalBasic.Config.GetString(configName + "maxconnection")
	dbdebug := globalBasic.Config.GetString(configName + "debug")

	config := DatabaseConfig{}
	config.Driver = dbdirver
	config.Host = dbhost
	config.Port, _ = strconv.Atoi(dbport)
	config.User = dbuser
	config.Passowrd = dbpassword
	config.Database = dbdatabase
	config.Debug, _ = strconv.ParseBool(dbdebug)
	config.MaxConnection, _ = strconv.Atoi(dbmaxconnection)

	return NewDatabase(config)
}

type zeroable interface {
	IsZero() bool
}

func (this *databaseImplement) rValue(bean interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(bean))
}

func (this *databaseImplement) isZero(k interface{}) bool {
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

func (this *databaseImplement) value2Interface(fieldValue reflect.Value) (interface{}, error) {
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
	case reflect.Float32, reflect.Float64:
		return fieldValue.Float(), nil
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

func (this *databaseImplement) autoMapType(v reflect.Value) *core.Table {
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

func (this *databaseImplement) NewSession() DatabaseSession {
	return this.NewSession()
}

func (this *databaseImplement) Sql(querystring string, args ...interface{}) DatabaseSession {
	return this.Sql(querystring, args...)
}

func (this *databaseImplement) NoAutoTime() DatabaseSession {
	return this.NoAutoTime()
}

func (this *databaseImplement) NoAutoCondition(no ...bool) DatabaseSession {
	return this.NoAutoCondition(no...)
}

func (this *databaseImplement) Cascade(trueOrFalse ...bool) DatabaseSession {
	return this.Cascade(trueOrFalse...)
}

func (this *databaseImplement) Where(querystring string, args ...interface{}) DatabaseSession {
	return this.Where(querystring, args...)
}

func (this *databaseImplement) Id(id interface{}) DatabaseSession {
	return this.Id(id)
}

func (this *databaseImplement) Distinct(columns ...string) DatabaseSession {
	return this.Distinct(columns...)
}

func (this *databaseImplement) Select(str string) DatabaseSession {
	return this.Select(str)
}

func (this *databaseImplement) Cols(columns ...string) DatabaseSession {
	return this.Cols(columns...)
}

func (this *databaseImplement) AllCols() DatabaseSession {
	return this.AllCols()
}

func (this *databaseImplement) MustCols(columns ...string) DatabaseSession {
	return this.MustCols(columns...)
}

func (this *databaseImplement) UseBool(columns ...string) DatabaseSession {
	return this.UseBool(columns...)
}

func (this *databaseImplement) Omit(columns ...string) DatabaseSession {
	return this.Omit(columns...)
}

func (this *databaseImplement) Nullable(columns ...string) DatabaseSession {
	return this.Nullable(columns...)
}

func (this *databaseImplement) In(column string, args ...interface{}) DatabaseSession {
	return this.In(column, args...)
}

func (this *databaseImplement) Incr(column string, args ...interface{}) DatabaseSession {
	return this.Incr(column, args...)
}

func (this *databaseImplement) Decr(column string, args ...interface{}) DatabaseSession {
	return this.Decr(column, args...)
}

func (this *databaseImplement) SetExpr(column string, expression string) DatabaseSession {
	return this.SetExpr(column, expression)
}

func (this *databaseImplement) Table(tableNameOrBean interface{}) DatabaseSession {
	return this.Table(tableNameOrBean)
}

func (this *databaseImplement) Alias(alias string) DatabaseSession {
	return this.Alias(alias)
}

func (this *databaseImplement) Limit(limit int, start ...int) DatabaseSession {
	return this.Limit(limit, start...)
}

func (this *databaseImplement) Desc(colNames ...string) DatabaseSession {
	return this.Desc(colNames...)
}

func (this *databaseImplement) Asc(colNames ...string) DatabaseSession {
	return this.Asc(colNames...)
}

func (this *databaseImplement) OrderBy(order string) DatabaseSession {
	return this.OrderBy(order)
}

func (this *databaseImplement) Join(join_operator string, tablename interface{}, condition string, args ...interface{}) DatabaseSession {
	return this.Join(join_operator, tablename, condition, args...)
}

func (this *databaseImplement) GroupBy(keys string) DatabaseSession {
	return this.GroupBy(keys)
}

func (this *databaseImplement) Having(conditions string) DatabaseSession {
	return this.Having(conditions)
}

func (this *databaseImplement) UpdateBatch(rowsSlicePtr interface{}, indexColName string) (int64, error) {
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
	updateCols := make([]bool, 0)
	var indexCol *core.Column

	//提取字段
	for i := 0; i < size; i++ {
		v := sliceValue.Index(i)
		vv := reflect.Indirect(v)

		//处理需要的update的列
		if i == 0 {
			for _, col := range table.Columns() {
				if col.Name == indexColName {
					indexCol = col
				} else {
					cols = append(cols, col)
					updateCols = append(updateCols, false)
				}
			}
			if indexCol == nil {
				return 0, errors.New("counld not found index col " + indexColName)
			}
		}

		//处理需要的update的值
		var singleRow = make([]interface{}, 0)
		for colIndex, col := range cols {
			ptrFieldValue, err := col.ValueOfV(&vv)
			if err != nil {
				return 0, err
			}
			fieldValue := *ptrFieldValue
			var arg interface{}
			if this.isZero(fieldValue.Interface()) {
				arg = nil
			} else {
				var err error
				arg, err = this.value2Interface(fieldValue)
				if err != nil {
					return 0, err
				}
				updateCols[colIndex] = true
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
	var isFirstUpdateCol = true
	for colIndex, col := range cols {
		if updateCols[colIndex] == false {
			continue
		}
		if isFirstUpdateCol == false {
			sql += " , "
		}
		sql += this.Engine.QuoteStr() + col.Name + this.Engine.QuoteStr()
		sql += " = CASE "
		sql += this.Engine.QuoteStr() + indexCol.Name + this.Engine.QuoteStr()
		for rowIndex, row := range rows {
			if row[colIndex] == nil {
				continue
			}
			sql += " WHEN ? THEN ? "
			sqlArgs = append(sqlArgs, indexRow[rowIndex])
			sqlArgs = append(sqlArgs, row[colIndex])
		}
		sql += " END "
		isFirstUpdateCol = false
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
