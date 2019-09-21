package sqlf

import (
	gosql "database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type sqlTypeOperation struct {
	toArgs     func(v interface{}, in []interface{}) (string, []interface{}, error)
	fromResult func(v interface{}, rows *gosql.Rows) error
	column     func() (string, error)
	setValue   func(v interface{}, in []interface{}) (string, []interface{}, error)
}

var (
	sqlTypeOperationMap = map[reflect.Type]sqlTypeOperation{}
	sqlTypeRwLock       = sync.RWMutex{}
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
			tempArgSql, err = operation[argsIndex].column()
			if err != nil {
				return "", nil, err
			}
			sqlBuilder.WriteString(tempArgSql)
		} else if len(query) >= len(setValueSql) && query[0:len(setValueSql)] == setValueSql {
			//提取set sql
			query = query[len(setValueSql)+1:]
			tempArgSql, realArgs, err = operation[argsIndex].setValue(args[argsIndex], realArgs)
			if err != nil {
				return "", nil, err
			}
			sqlBuilder.WriteString(tempArgSql)
		} else {
			//普通的提取方式
			query = query[1:]
			tempArgSql, realArgs, err = operation[argsIndex].toArgs(args[argsIndex], realArgs)
			if err != nil {
				return "", nil, err
			}
			sqlBuilder.WriteString(tempArgSql)
		}
	}
	return sqlBuilder.String(), realArgs, nil
}

func getSqlOperationFromInterface(i interface{}) sqlTypeOperation {
	return getSqlOperation(reflect.TypeOf(i))
}

func getSqlOperation(t reflect.Type) sqlTypeOperation {
	sqlTypeRwLock.RLock()
	result, isExist := sqlTypeOperationMap[t]
	sqlTypeRwLock.RUnlock()

	if isExist == true {
		return result
	}

	result = initSqlOperation(t)

	sqlTypeRwLock.Lock()
	sqlTypeOperationMap[t] = result
	sqlTypeRwLock.Unlock()

	return result
}

func initSqlOperation(t reflect.Type) sqlTypeOperation {
	return sqlTypeOperation{}
}
