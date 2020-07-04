package sqlf

import (
	gosql "database/sql"
	"errors"
	"reflect"
	"strings"
)

func initByteSliceSqlTypeOperation() {
	var a []byte
	byteSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, v)
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("[]byte dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("[]byte dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]byte dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(byteSliceType, &sqlTypeOperation)
}

func initByteSlicePtrSqlTypeOperation() {
	var a *[]byte
	byteSlicePtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.(*[]byte)
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
			return errors.New("*[]byte dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*[]byte dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(byteSlicePtrType, &sqlTypeOperation)
}

func initByteSliceSliceSqlTypeOperation() {
	var a [][]byte
	byteSliceSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([][]byte)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("[][]byte dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("[][]byte dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[][]byte dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(byteSliceSliceType, &sqlTypeOperation)
}

func initByteSliceSlicePtrSqlTypeOperation() {
	var a *[][]byte
	byteSliceSlicePtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := *(v.(*[][]byte))
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			data := v.(*[][]byte)
			result := [][]byte{}
			var temp []byte
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
			return errors.New("*[][]byte dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*[][]byte dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(byteSliceSlicePtrType, &sqlTypeOperation)
}

func init() {
	initByteSliceSqlTypeOperation()
	initByteSlicePtrSqlTypeOperation()
	initByteSliceSliceSqlTypeOperation()
	initByteSliceSlicePtrSqlTypeOperation()
}
