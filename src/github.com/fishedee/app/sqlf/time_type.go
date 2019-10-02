package sqlf

import (
	gosql "database/sql"
	"errors"
	"reflect"
	"strings"
	"time"
)

func initTimeSqlTypeOperation() {
	zeroTime := time.Time{}
	zeroTimeString := zeroTime.Local().Format("2006-01-02T15:04:05")

	a := time.Time{}
	stringType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.(time.Time)
			builder.WriteByte('?')
			if data.IsZero() && driver == "mysql" {
				in = append(in, zeroTimeString)
			} else {
				in = append(in, v)
			}

			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("time.Time dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("time.Time dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("time.Time dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringType, &sqlTypeOperation)
}

func initTimePtrSqlTypeOperation() {
	zeroTime := time.Time{}
	zeroTimeString := zeroTime.Local().Format("2006-01-02T15:04:05")

	var a *time.Time
	timePtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.(*time.Time)
			builder.WriteByte('?')
			if data.IsZero() && driver == "mysql" {
				in = append(in, zeroTimeString)
			} else {
				in = append(in, *data)
			}
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
			return errors.New("*time.Time dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*time.Time dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(timePtrType, &sqlTypeOperation)
}

func initTimeSliceSqlTypeOperation() {
	zeroTime := time.Time{}
	zeroTimeString := zeroTime.Local().Format("2006-01-02T15:04:05")

	a := []time.Time{}
	stringSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([]time.Time)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				if single.IsZero() && driver == "mysql" {
					in = append(in, zeroTimeString)
				} else {
					in = append(in, single)
				}
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("[]time.Time dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("[]time.Time dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]time.Time dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringSliceType, &sqlTypeOperation)
}

func initTimeSlicePtrSqlTypeOperation() {
	zeroTime := time.Time{}
	zeroTimeString := zeroTime.Local().Format("2006-01-02T15:04:05")

	var a *[]time.Time
	stringSlicePtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := *(v.(*[]time.Time))
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				if single.IsZero() && driver == "mysql" {
					in = append(in, zeroTimeString)
				} else {
					in = append(in, single)
				}
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			data := v.(*[]time.Time)
			result := []time.Time{}
			var temp time.Time
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
			return errors.New("*[]time.Time dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*[]time.Time dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(stringSlicePtrType, &sqlTypeOperation)
}

func init() {
	initTimeSqlTypeOperation()
	initTimePtrSqlTypeOperation()
	initTimeSliceSqlTypeOperation()
	initTimeSlicePtrSqlTypeOperation()
}
