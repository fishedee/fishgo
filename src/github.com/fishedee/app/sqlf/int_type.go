package sqlf

import (
	gosql "database/sql"
	"errors"
	"reflect"
	"strings"
)

func initIntSqlTypeOperation() {
	var a int = 10
	intType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, v)
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("int dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("int dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("int dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(intType, &sqlTypeOperation)
}

func initIntPtrSqlTypeOperation() {
	var a *int
	intPtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.(*int)
			builder.WriteByte('?')
			in = append(in, *data)
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			if rows.Next() {
				err := rows.Scan(v)
				if err != nil {
					return err
				}
				return nil
			} else {
				return errors.New("has no result")
			}
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("*int dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*int dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(intPtrType, &sqlTypeOperation)
}

func initIntSliceSqlTypeOperation() {
	a := []int{10}
	intSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([]int)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("[]int dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("[]int dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]int dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(intSliceType, &sqlTypeOperation)
}

func initIntSlicePtrSqlTypeOperation() {
	var a *[]int
	intSlicePtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := *(v.(*[]int))
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			data := v.(*[]int)
			result := []int{}
			var temp int
			for rows.Next() {
				err := rows.Scan(&temp)
				if err != nil {
					return err
				}
				result = append(result, temp)
			}
			*data = result
			return nil
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("*[]int dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*[]int dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(intSlicePtrType, &sqlTypeOperation)
}

func init() {
	initIntSqlTypeOperation()
	initIntPtrSqlTypeOperation()
	initIntSliceSqlTypeOperation()
	initIntSlicePtrSqlTypeOperation()
}
