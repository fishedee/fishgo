package language

import (
	"reflect"
	"fmt"
)

func Table(column interface{},data interface{})([][]string){
	result := [][]string{}

	columnKeys,columnValues := MapToSlice(column)
	columnKeysReal := columnKeys.([]string)
	columnValuesReal := columnValues.([]string)
	result = append(result,columnValuesReal)

	dataValue := reflect.ValueOf(data)
	for i := 0 ; i != dataValue.Len() ; i++{
		singleDataValue := dataValue.Index(i)
		singleDataStringValue := StructToMap(singleDataValue.Interface())
		singleDataStringValueData := reflect.ValueOf(singleDataStringValue)
		singleResult := []string{}
		for _,singleColumn := range columnKeysReal{
			singleResultString := ""
			singleValue := singleDataStringValueData.MapIndex(reflect.ValueOf(singleColumn))
			if singleValue.IsValid() == false{
				singleResultString = ""
			}else{
				singleResultString = fmt.Sprintf("%v",singleValue)
			}
			singleResult = append(
				singleResult,
				singleResultString,
			)
		}
		result = append(result,singleResult)
	}
	return result
}