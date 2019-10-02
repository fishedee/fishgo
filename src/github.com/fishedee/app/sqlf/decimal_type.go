package sqlf

import (
	gosql "database/sql"
	"errors"
	. "github.com/fishedee/language"
	"reflect"
	"strings"
)

func initDecimalSqlTypeOperation() {
	a := Decimal("")
	decimalType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, v)
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("Decimal dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("Decimal dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("Decimal dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(decimalType, &sqlTypeOperation)
}

func initDecimalPtrSqlTypeOperation() {
	var a *Decimal
	decimalPtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.(*Decimal)
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
			return errors.New("*Decimal dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*Decimal dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(decimalPtrType, &sqlTypeOperation)
}

func initDecimalSliceSqlTypeOperation() {
	a := []Decimal{}
	decimalSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([]Decimal)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("[]Decimal dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("[]Decimal dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]Decimal dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(decimalSliceType, &sqlTypeOperation)
}

func initDecimalSlicePtrSqlTypeOperation() {
	var a *[]Decimal
	stringSlicePtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := *(v.(*[]Decimal))
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			data := v.(*[]Decimal)
			result := []Decimal{}
			var temp Decimal
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
			return errors.New("*[]Decimal dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*[]Decimal dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringSlicePtrType, &sqlTypeOperation)
}

func init() {
	initDecimalSqlTypeOperation()
	initDecimalPtrSqlTypeOperation()
	initDecimalSliceSqlTypeOperation()
	initDecimalSlicePtrSqlTypeOperation()
}
