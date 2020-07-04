package sqlf

import (
	gosql "database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

func initJsonRawMessageSqlTypeOperation() {
	var a json.RawMessage
	jsonRawMessageType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			builder.WriteByte('?')
			in = append(in, v)
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("json.RawMessage dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("json.RawMessage dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("json.RawMessage dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(jsonRawMessageType, &sqlTypeOperation)
}

func initJsonRawMessagePtrSqlTypeOperation() {
	var a *json.RawMessage
	jsonRawMessagePtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.(*json.RawMessage)
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
			return errors.New("*json.RawMessage dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*json.RawMessage dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(jsonRawMessagePtrType, &sqlTypeOperation)
}

func initJsonRawMessageSliceSqlTypeOperation() {
	var a []json.RawMessage
	jsonRawMessageSliceType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := v.([]json.RawMessage)
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			return errors.New("[]json.RawMessage dos not support setValue")
		},
		column: func(driver string, isInsert bool, builder *strings.Builder) error {
			return errors.New("[]json.RawMessage dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("[]json.RawMessage dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(jsonRawMessageSliceType, &sqlTypeOperation)
}

func initJsonRawMessageSlicePtrSqlTypeOperation() {
	var a *[]json.RawMessage
	jsonRawMessageSlicePtrType := reflect.TypeOf(a)
	sqlTypeOperation := sqlTypeOperation{
		toArgs: func(driver string, isInsert bool, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			data := *(v.(*[]json.RawMessage))
			builder.WriteString(getSqlComma(len(data)))
			for _, single := range data {
				in = append(in, single)
			}
			return in, nil
		},
		fromResult: func(driver string, v interface{}, rows *gosql.Rows) error {
			data := v.(*[]json.RawMessage)
			result := []json.RawMessage{}
			var temp json.RawMessage
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
			return errors.New("*[]json.RawMessage dos not support column")
		},
		setValue: func(driver string, v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New("*[]json.RawMessage dos not support setValue")
		},
	}
	sqlTypeOperationMap.Store(jsonRawMessageSlicePtrType, &sqlTypeOperation)
}

func init() {
	initJsonRawMessageSqlTypeOperation()
	initJsonRawMessagePtrSqlTypeOperation()
	initJsonRawMessageSliceSqlTypeOperation()
	initJsonRawMessageSlicePtrSqlTypeOperation()
}
