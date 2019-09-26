package sqlf

import (
	gosql "database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
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
	commaCache          = []string{}
	stringPtrType       = reflect.TypeOf((*string)(nil))
)

func extractResult(data interface{}, rows *gosql.Rows) error {
	operation := getSqlOperationFromInterface(data)
	return operation.fromResult(data, rows)
}

func notWordChar(data uint8) bool {
	if data >= '0' && data <= '9' {
		return false
	} else if data >= 'A' && data <= 'Z' {
		return false
	} else if data >= 'a' && data <= 'z' {
		return false
	} else {
		return true
	}
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
	sqlBuilder.Grow(len(query) * 2)
	var err error
	for {
		index := strings.IndexByte(query, '?')
		if index == -1 {
			sqlBuilder.WriteString(query)
			break
		}
		sqlBuilder.WriteString(query[0:index])
		query = query[index:]
		if argsIndex >= len(args) {
			return "", nil, errors.New(fmt.Sprintf("invalid ? index %v,%v", argsIndex, len(args)))
		}
		columnSql := "?.column"
		setValueSql := "?.setValue"
		if len(query) > len(columnSql) && query[0:len(columnSql)] == columnSql && notWordChar(query[len(columnSql)]) == true {
			//提取column的name
			query = query[len(columnSql):]
			err = operation[argsIndex].column(&sqlBuilder)
			if err != nil {
				return "", nil, err
			}
		} else if len(query) > len(setValueSql) && query[0:len(setValueSql)] == setValueSql && notWordChar(query[len(setValueSql)]) == true {
			//提取set sql
			query = query[len(setValueSql):]
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
		argsIndex++
	}
	return sqlBuilder.String(), realArgs, nil
}

func getSqlOperationFromInterface(i interface{}) sqlTypeOperation {
	return getSqlOperation(reflect.TypeOf(i))
}

func getSqlOperation(t reflect.Type) sqlTypeOperation {
	result, isExist := sqlTypeOperationMap.Load(t)

	if isExist == true {
		return *(result.(*sqlTypeOperation))
	}

	newResult := initStructTypeOperation(t)
	sqlTypeOperationMap.Store(t, &newResult)
	return newResult
}
