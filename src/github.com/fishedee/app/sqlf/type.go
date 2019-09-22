package sqlf

import (
	gosql "database/sql"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

type sqlToArgsType = func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error)

type sqlFromResultType = func(v interface{}, rows *gosql.Rows) error

type sqlColumnType = func(builder *strings.Builder) error

type sqlSetValueType = func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error)

type sqlTypeOperation struct {
	toArgs     sqlToArgsType
	fromResult sqlFromResultType
	column     sqlColumnType
	setValue   sqlSetValueType
}

var (
	sqlTypeOperationMap = sync.Map{}
)

func extractResult(data interface{}, rows *gosql.Rows) error {
	operation := getSqlOperationFromInterface(data)
	return operation.fromResult(data, rows)
}

func genSql(query string, args []interface{}) (string, []interface{}, error) {
	//获得operation
	operation := make([]sqlTypeOperation, len(args), len(args))
	for i, arg := range args {
		operation[i] = getSqlOperationFromInterface(arg)
	}

	//拼凑sql
	realArgs := make([]interface{}, 0, len(args))
	argsIndex := 0
	sqlBuilder := strings.Builder{}
	tempArgSql := ""
	var err error
	for {
		index := strings.IndexByte(query, '?')
		if index == -1 {
			sqlBuilder.WriteString(query)
			break
		}
		sqlBuilder.WriteString(query[0 : index-1])
		query = query[index:]
		if argsIndex >= len(args) {
			return "", nil, errors.New(fmt.Sprintf("invalid ? index %v,%v", argsIndex, len(args)))
		}
		columnSql := ".column "
		setValueSql := ".setValue "
		if len(query) >= len(columnSql) && query[0:len(columnSql)] == columnSql {
			//提取column的name
			query = query[len(columnSql)+1:]
			err = operation[argsIndex].column(&sqlBuilder)
			if err != nil {
				return "", nil, err
			}
		} else if len(query) >= len(setValueSql) && query[0:len(setValueSql)] == setValueSql {
			//提取set sql
			query = query[len(setValueSql)+1:]
			realArgs, err = operation[argsIndex].setValue(args[argsIndex], realArgs, &sqlBuilder)
			if err != nil {
				return "", nil, err
			}
		} else {
			//普通的提取方式
			query = query[1:]
			realArgs, err = operation[argsIndex].toArgs(args[argsIndex], realArgs, &sqlBuilder)
			if err != nil {
				return "", nil, err
			}
		}
	}
	return sqlBuilder.String(), realArgs, nil
}

func getSqlOperationFromInterface(i interface{}) sqlTypeOperation {
	return getSqlOperation(reflect.TypeOf(i))
}

func getSqlOperation(t reflect.Type) (sqlTypeOperation, error) {
	result, isExist := sqlTypeOperationMap.Load(t)

	if isExist == true {
		return *(result.(*sqlTypeOperation))
	}

	newResult = initSqlOperation(t)
	sqlTypeOperationMap.Store(t, &newResult)

	return newResult
}

func initSqlOperation(t reflect.Type) (sqlTypeOperation, error) {
	return sqlTypeOperation{
		toArgs:     initSqlToArgs(t),
		fromResult: initSqlFromResult(t),
		column:     initSqlColumn(t),
		setValue:   initSqlSetValue(t),
	}
}

func getTypeKind(t reflect.Type) int {
	timeType := reflect.TypeOf(time.Time{})
	if t.Kind() == reflect.Struct && t != timeType {
		//struct类型
		return 1
	} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct && t.Elem() != timeType {
		//*struct类型
		return 2
	} else if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Struct && t.Elem() != timeType {
		//[]struct类型
		return 3
	} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice && t.Elem().Elem().Kind() == reflect.Struct && t.Elem().Elem() != timeType {
		//*[]struct类型
		return 4
	} else {
		return 5
	}
}

type sqlStructPublicField struct {
	name  string
	index []int
}

func getStructPublicField(t reflect.Type) []sqlStructPublicField {
	result := []sqlStructPublicField{}
	numField := t.NumField()
	for i := 0; i != numField; i++ {
		field := t.Field(i)
		fieldName := field.Name
		if fieldName[0] >= 'A' && fieldName[0] <= 'Z' {
			result = append(result, sqlStructPublicField{
				name:  field.Name,
				index: field.Index,
			})
		}
	}
	sort.Sort(result, func(i int, j int) {
		return result[i].name < result[j].name
	})
	return result
}

func getSqlComma(num int) string {
	if num == 1 {
		return "?"
	} else {
		return strings.Repeat("?,", num-1) + "?"
	}
}

func initSqlToArgs(t reflect.Type) sqlToArgsType {
	tKind := getTypeKind(t)
	if tKind == 5 {
		tName := t.String()
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New(fmt.Sprintf("%v dos not support toArgs", tName))
		}
	}

	structToArgs := func(t reflect.Type) func(v reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
		fields := getStructPublicField(t)
		return func(v reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			//生成占位符
			builder.WriteString(getSqlComma(len(fields)))

			//生成字段
			for _, field := range fields {
				in = append(in, v.FieldByIndex(field.index).Interface())
			}
			return in, nil
		}
	}

	if tKind == 1 {
		structHandler := structToArgs(t)
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			value := reflect.ValueOf(v)
			return structHandler(value, in, builder)
		}
	} else if tKind == 2 {
		structHandler := structToArgs(t.Elem())
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			value := reflect.ValueOf(v).Elem()
			return structHandler(value, in, builder)
		}
	} else {
		structSliceToArgs := func(t reflect.Type) func(value reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			structHandler := structToArgs(t.Elem())

			return func(value reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
				length := value.Len()
				for i := 0; i != length; i++ {
					if i != 0 {
						builder.WriteString(",(")
					} else {
						builder.WriteString("(")
					}
					in, err = structHandler(value.Index(i), in, builder)
					if err != nil {
						return nil, err
					}
					builder.WriteString(")")
				}
				return in, nil
			}
		}

		if tKind == 3 {
			structSliceHandler := structSliceToArgs(t)
			return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
				value := reflect.ValueOf(v)
				return structSliceHandler(value, in, builder)
			}
		} else {
			structSliceHandler := structSliceToArgs(t.Elem())
			return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
				value := reflect.ValueOf(v).Elem()
				return structSliceHandler(value, in, builder)
			}
		}
	}
}

func initSqlColumn(t reflect.Type) sqlColumnType {
	tKind := getTypeKind(t)
	if tKind == 5 {
		tName := t.String()
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New(fmt.Sprintf("%v dos not support toArgs", tName))
		}
	}

	structColumn := func(t reflect.Type) string {
		fields := getStructPublicField(t)

		buffer := strings.Buffer{}
		for i, field := range filed {
			if i != 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString("`")
			buffer.WriteString(field.name)
			buffer.WriteString("`")
		}
		return buffer.String()
	}
	var result = ""

	if tKind == 1 {
		result = structColumn(t)
	} else if tKind == 2 {
		result = structColumn(t.Elem())
	} else if tKind == 3 {
		result = structColumn(t.Elem())
	} else {
		result = structColumn(t.Elem().Elem())
	}

	return func(builder *strings.Builder) error {
		builder.WriteString(result)
		return nil
	}
}
