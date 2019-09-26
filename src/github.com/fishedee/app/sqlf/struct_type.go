package sqlf

import (
	gosql "database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func initStructTypeOperation(t reflect.Type) sqlTypeOperation {
	return sqlTypeOperation{
		toArgs:     initSqlToArgs(t),
		fromResult: initSqlFromResult(t),
		column:     initSqlColumn(t),
		setValue:   initSqlSetValue(t),
	}
}

func getTypeKind(t reflect.Type) int {
	timeType := reflect.TypeOf(time.Time{})
	if t.Kind() == reflect.Struct && t != timeType {
		//struct类型
		return 1
	} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct && t.Elem() != timeType {
		//*struct类型
		return 2
	} else if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Struct && t.Elem() != timeType {
		//[]struct类型
		return 3
	} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice && t.Elem().Elem().Kind() == reflect.Struct && t.Elem().Elem() != timeType {
		//*[]struct类型
		return 4
	} else {
		return 5
	}
}

type sqlStructPublicField struct {
	name  string
	index []int
}

func getStructPublicField(t reflect.Type) []sqlStructPublicField {
	result := []sqlStructPublicField{}
	numField := t.NumField()
	for i := 0; i != numField; i++ {
		field := t.Field(i)
		fieldName := field.Name
		if fieldName[0] >= 'A' && fieldName[0] <= 'Z' {
			fieldName = strings.ToLower(fieldName[0:1]) + fieldName[1:]
			result = append(result, sqlStructPublicField{
				name:  fieldName,
				index: field.Index,
			})
		}
	}
	return result
}

func initSqlToArgs(t reflect.Type) sqlToArgsType {
	tKind := getTypeKind(t)
	if tKind == 5 {
		tName := t.String()
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New(fmt.Sprintf("%v dos not support toArgs", tName))
		}
	}

	structToArgs := func(t reflect.Type) func(v reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
		fields := getStructPublicField(t)
		return func(v reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			//生成占位符
			builder.WriteString(getSqlComma(len(fields)))

			//生成字段
			for _, field := range fields {
				in = append(in, v.FieldByIndex(field.index).Interface())
			}
			return in, nil
		}
	}

	if tKind == 1 {
		structHandler := structToArgs(t)
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			value := reflect.ValueOf(v)
			return structHandler(value, in, builder)
		}
	} else if tKind == 2 {
		structHandler := structToArgs(t.Elem())
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			value := reflect.ValueOf(v).Elem()
			return structHandler(value, in, builder)
		}
	} else {
		structSliceToArgs := func(t reflect.Type) func(value reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			structHandler := structToArgs(t.Elem())

			return func(value reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
				var err error
				length := value.Len()
				for i := 0; i != length; i++ {
					if i != 0 {
						builder.WriteString(",(")
					} else {
						builder.WriteByte('(')
					}
					in, err = structHandler(value.Index(i), in, builder)
					if err != nil {
						return nil, err
					}
					builder.WriteByte(')')
				}
				return in, nil
			}
		}

		if tKind == 3 {
			structSliceHandler := structSliceToArgs(t)
			return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
				value := reflect.ValueOf(v)
				return structSliceHandler(value, in, builder)
			}
		} else {
			structSliceHandler := structSliceToArgs(t.Elem())
			return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
				value := reflect.ValueOf(v).Elem()
				return structSliceHandler(value, in, builder)
			}
		}
	}
}

func initSqlColumn(t reflect.Type) sqlColumnType {
	tKind := getTypeKind(t)
	if tKind == 5 {
		tName := t.String()
		return func(builder *strings.Builder) error {
			return errors.New(fmt.Sprintf("%v dos not support column", tName))
		}
	}

	structColumn := func(t reflect.Type) string {
		fields := getStructPublicField(t)

		builder := strings.Builder{}
		for i, field := range fields {
			if i != 0 {
				builder.WriteByte(',')
			}
			builder.WriteByte('`')
			builder.WriteString(field.name)
			builder.WriteByte('`')
		}
		return builder.String()
	}
	var result = ""

	if tKind == 1 {
		result = structColumn(t)
	} else if tKind == 2 {
		result = structColumn(t.Elem())
	} else if tKind == 3 {
		result = structColumn(t.Elem())
	} else {
		result = structColumn(t.Elem().Elem())
	}

	return func(builder *strings.Builder) error {
		builder.WriteString(result)
		return nil
	}
}

func initSqlFromResult(t reflect.Type) sqlFromResultType {
	tKind := getTypeKind(t)
	if tKind != 4 {
		tName := t.String()
		return func(v interface{}, rows *gosql.Rows) error {
			return errors.New(fmt.Sprintf("%v dos not support fromResult", tName))
		}
	}
	sliceType := t.Elem()
	structType := sliceType.Elem()
	structInfo := getStructPublicField(structType)
	fieldInfoMap := map[string]sqlStructPublicField{}
	for _, single := range structInfo {
		fieldInfoMap[single.name] = single
	}
	return func(v interface{}, rows *gosql.Rows) error {
		//配置列
		columns, err := rows.Columns()
		if err != nil {
			return err
		}
		temp := reflect.New(structType).Elem()
		tempScan := make([]interface{}, len(columns), len(columns))
		for i, column := range columns {
			fieldInfo, isExist := fieldInfoMap[column]
			if isExist == false {
				return errors.New(fmt.Sprintf("%v dos not have column %v", structType.String(), column))
			}

			tempScan[i] = temp.FieldByIndex(fieldInfo.index).Addr().Interface()
		}

		//写入数组
		result := reflect.MakeSlice(sliceType, 0, 16)
		for rows.Next() {
			err := rows.Scan(tempScan...)
			if err != nil {
				return err
			}
			result = reflect.Append(result, temp)
		}
		reflect.ValueOf(v).Elem().Set(result)
		return nil
	}
}

func initSqlSetValue(t reflect.Type) sqlSetValueType {
	tKind := getTypeKind(t)
	if tKind != 1 && tKind != 2 {
		tName := t.String()
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			return nil, errors.New(fmt.Sprintf("%v dos not support setValue", tName))
		}
	}
	structSetValue := func(t reflect.Type) func(v reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
		fields := getStructPublicField(t)
		return func(v reflect.Value, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			for i, field := range fields {
				if i != 0 {
					builder.WriteByte(',')
				}
				builder.WriteByte('`')
				builder.WriteString(field.name)
				builder.WriteString("` = ? ")
				in = append(in, v.FieldByIndex(field.index).Interface())
			}
			return in, nil
		}
	}

	if tKind == 1 {
		structSetValueHandler := structSetValue(t)
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			value := reflect.ValueOf(v)
			return structSetValueHandler(value, in, builder)
		}
	} else {
		structSetValueHandler := structSetValue(t.Elem())
		return func(v interface{}, in []interface{}, builder *strings.Builder) ([]interface{}, error) {
			value := reflect.ValueOf(v).Elem()
			return structSetValueHandler(value, in, builder)
		}
	}
}

func getSqlComma(num int) string {
	if num < len(commaCache) {
		return commaCache[num]
	} else {
		return strings.Repeat("?,", num-1) + "?"
	}
}

func init() {
	commaCache := make([]string, 128, 128)
	commaCache[1] = "?"
	for i := 2; i != len(commaCache); i++ {
		commaCache[i] = commaCache[i-1] + ",?"
	}
}
