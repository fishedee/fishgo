package language

import (
	"reflect"
	"strconv"
	"strings"
)

type EnumStruct struct {
	names map[string]string
}

func InitEnumStruct( this interface{} ){
	enumInfo := reflect.TypeOf(this).Elem()
	enumValue := reflect.ValueOf(this)
	result := enumValue.Elem().FieldByName("EnumStruct").Addr().Interface().(*EnumStruct);
	result.names = map[string]string{}

	for i := 0 ; i != enumInfo.NumField() ; i++{
		singleField := enumInfo.Field(i)

		singleFieldName := singleField.Name
		singleFieldTag := singleField.Tag.Get("enum")
		singleFieldTagArray := strings.Split(singleFieldTag,",")
		if len(singleFieldTagArray) != 2{
			continue
		}

		singleFieldTagValue,err := strconv.Atoi( singleFieldTagArray[0] )
		if err != nil{
			panic(singleFieldName+": "+singleFieldTag+" is not a integer")
		}
		singleFieldTagSeeName := singleFieldTagArray[1]

		result.names[singleFieldTagArray[0]] = singleFieldTagSeeName
		enumValue.Elem().Field(i).SetInt(int64(singleFieldTagValue))
	}
}


func (this *EnumStruct) Names() map[string]string {
	return this.names;
}

