package sqlf

import (
	gosql "database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const (
	NormalColumn      = "?.column"
	NormalValue       = "?"
	InsertColumn      = "?.insertColumn"
	InsertValue       = "?.insertValue"
	UpdateColumnValue = "?.updateColumnValue"
)

type sqlToArgsType = func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error)

type sqlFromResultType = func(driver string, v interface{}, rows *gosql.Rows) error

type sqlColumnType = func(driver string, isInsert bool, builder *strings.Builder) error

type sqlSetValueType = func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error)

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

func extractResult(driver string, data interface{}, rows *gosql.Rows) error {
	operation := getSqlOperationFromInterface(data)
	return operation.fromResult(driver, data, rows)
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

func checkStartWith(query string, match string) bool {
	if len(query) < len(match) {
		return false
	}
	if query[0:len(match)] != match {
		return false
	}
	if len(query) == len(match) ||
		notWordChar(query[len(match)]) == true {
		return true
	} else {
		return false
	}
}

func genSql(driver string, query string, args []interface{}) (string, []interface{}, error) {
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
		if checkStartWith(query, InsertColumn) {
			//提取insert的column
			query = query[len(InsertColumn):]
			err = operation[argsIndex].column(driver, true, &sqlBuilder)
			if err != nil {
				return "", nil, err
			}
		} else if checkStartWith(query, NormalColumn) {
			//提取normal的column
			query = query[len(NormalColumn):]
			err = operation[argsIndex].column(driver, false, &sqlBuilder)
			if err != nil {
				return "", nil, err
			}
		} else if checkStartWith(query, UpdateColumnValue) {
			//提取update的column与value
			query = query[len(UpdateColumnValue):]
			realArgs, err = operation[argsIndex].setValue(driver, args[argsIndex], realArgs, &sqlBuilder)
			if err != nil {
				return "", nil, err
			}
		} else if checkStartWith(query, InsertValue) {
			//提取insert的value
			query = query[len(InsertValue):]
			realArgs, err = operation[argsIndex].toArgs(driver, true, args[argsIndex], realArgs, &sqlBuilder)
			if err != nil {
				return "", nil, err
			}
		} else {
			//普通的提取方式
			query = query[1:]
			realArgs, err = operation[argsIndex].toArgs(driver, false, args[argsIndex], realArgs, &sqlBuilder)
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
