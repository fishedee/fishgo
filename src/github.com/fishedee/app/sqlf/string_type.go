package sqlf

import (
	gosql "database/sql"
	"errors"
	"reflect"
	"strings"
)

func initStringSqlTypeOperation() {
	a := ""
	stringType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, v)
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("string dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("string dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("string dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringType, &sqlTypeOperation)
}

func initStringPtrSqlTypeOperation() {
	var a *string
	stringPtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.(*string)
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
			return errors.New("*string dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*string dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringPtrType, &sqlTypeOperation)
}

func initStringSliceSqlTypeOperation() {
	a := []string{""}
	stringSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([]string)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("[]string dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("[]string dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]string dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringSliceType, &sqlTypeOperation)
}

func initStringSlicePtrSqlTypeOperation() {
	var a *[]string
	stringSlicePtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := *(v.(*[]string))
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			data := v.(*[]string)
			result := []string{}
			var temp string
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
			return errors.New("*[]string dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*[]string dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringSlicePtrType, &sqlTypeOperation)
}

func init() {
	initStringSqlTypeOperation()
	initStringPtrSqlTypeOperation()
	initStringSliceSqlTypeOperation()
	initStringSlicePtrSqlTypeOperation()
}
