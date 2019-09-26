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
		toArgs: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, v)
			return in, nil
		},
		fromResult: func(v interface{}, rows *gosql.Rows) error {
			return errors.New("Decimal dos not support setValue")
		},
		column: func(builder *strings.Builder) error {
			return errors.New("Decimal dos not support column")
		},
		setValue: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("Decimal dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(decimalType, &sqlTypeOperation)
}

func initDecimalSliceSqlTypeOperation() {
	a := []Decimal{}
	decimalSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([]Decimal)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(v interface{}, rows *gosql.Rows) error {
			return errors.New("[]Decimal dos not support setValue")
		},
		column: func(builder *strings.Builder) error {
			return errors.New("[]Decimal dos not support column")
		},
		setValue: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]Decimal dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(decimalSliceType, &sqlTypeOperation)
}

func init() {
	initDecimalSqlTypeOperation()
	initDecimalSliceSqlTypeOperation()
}
