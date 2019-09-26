package sqlf

import (
	gosql "database/sql"
	"errors"
	"reflect"
	"strings"
	"time"
)

func initTimeSqlTypeOperation() {
	a := time.Time{}
	stringType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, v)
			return in, nil
		},
		fromResult: func(v interface{}, rows *gosql.Rows) error {
			return errors.New("time.Time dos not support setValue")
		},
		column: func(builder *strings.Builder) error {
			return errors.New("time.Time dos not support column")
		},
		setValue: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("time.Time dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringType, &sqlTypeOperation)
}

func initTimeSliceSqlTypeOperation() {
	a := []time.Time{}
	stringSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([]time.Time)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(v interface{}, rows *gosql.Rows) error {
			return errors.New("[]time.Time dos not support setValue")
		},
		column: func(builder *strings.Builder) error {
			return errors.New("[]time.Time dos not support column")
		},
		setValue: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]time.Time dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringSliceType, &sqlTypeOperation)
}

func init() {
	initTimeSqlTypeOperation()
	initTimeSliceSqlTypeOperation()
}
