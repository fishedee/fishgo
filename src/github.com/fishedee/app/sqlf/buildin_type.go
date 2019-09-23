package sqlf

import (
	gosql "database/sql"
	"errors"
	. "github.com/fishedee/language"
	"reflect"
	"strings"
	"time"
)

func initIntSqlTypeOperation() {
	var a int = 10
	intType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, v)
			return in, nil
		},
		fromResult: func(v interface{}, rows *gosql.Rows) error {
			return errors.New("int dos not support setValue")
		},
		column: func(builder *strings.Builder) error {
			return errors.New("int dos not support column")
		},
		setValue: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("int dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(intType, &sqlTypeOperation)
}

func initIntSliceSqlTypeOperation() {
	a := []int{10}
	intSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([]int)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(v interface{}, rows *gosql.Rows) error {
			return errors.New("[]int dos not support setValue")
		},
		column: func(builder *strings.Builder) error {
			return errors.New("[]int dos not support column")
		},
		setValue: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]int dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(intSliceType, &sqlTypeOperation)
}

func initStringSqlTypeOperation() {
	a := ""
	stringType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, v)
			return in, nil
		},
		fromResult: func(v interface{}, rows *gosql.Rows) error {
			return errors.New("string dos not support setValue")
		},
		column: func(builder *strings.Builder) error {
			return errors.New("string dos not support column")
		},
		setValue: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("string dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringType, &sqlTypeOperation)
}

func initStringSliceSqlTypeOperation() {
	a := []string{""}
	stringSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([]string)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(v interface{}, rows *gosql.Rows) error {
			return errors.New("[]string dos not support setValue")
		},
		column: func(builder *strings.Builder) error {
			return errors.New("[]string dos not support column")
		},
		setValue: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]string dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringSliceType, &sqlTypeOperation)
}

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

func initDecimalSqlTypeOperation() {
	a := Decimal("")
	decimalType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, string(v.(Decimal)))
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
				in = append(in, string(single))
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
	initIntSqlTypeOperation()
	initIntSliceSqlTypeOperation()
	initStringSqlTypeOperation()
	initStringSliceSqlTypeOperation()
	initTimeSqlTypeOperation()
	initTimeSliceSqlTypeOperation()
	initDecimalSqlTypeOperation()
	initDecimalSliceSqlTypeOperation()
}
