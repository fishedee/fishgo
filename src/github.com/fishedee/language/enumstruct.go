package language

import (
	"reflect"
	"strconv"
	"strings"
)

type EnumData struct {
	Id   int
	Name string
}

type EnumStruct struct {
	names map[string]string
}

func InitEnumStruct(this interface{}) {
	enumInfo := reflect.TypeOf(this).Elem()
	enumValue := reflect.ValueOf(this)
	result := enumValue.Elem().FieldByName("EnumStruct").Addr().Interface().(*EnumStruct)
	result.names = map[string]string{}

	for i := 0; i != enumInfo.NumField(); i++ {
		singleField := enumInfo.Field(i)

		singleFieldName := singleField.Name
		singleFieldTag := singleField.Tag.Get("enum")
		singleFieldTagArray := strings.Split(singleFieldTag, ",")
		if len(singleFieldTagArray) != 2 {
			continue
		}

		singleFieldTagValue, err := strconv.Atoi(singleFieldTagArray[0])
		if err != nil {
			panic(singleFieldName + ": " + singleFieldTag + " is not a integer")
		}
		singleFieldTagSeeName := singleFieldTagArray[1]

		result.names[singleFieldTagArray[0]] = singleFieldTagSeeName
		enumValue.Elem().Field(i).SetInt(int64(singleFieldTagValue))
	}
}

func (this *EnumStruct) Names() map[string]string {
	return this.names
}

func (this *EnumStruct) Entrys() map[int]string {
	result := map[int]string{}
	for key, value := range this.names {
		singleKey, _ := strconv.Atoi(key)
		result[singleKey] = value
	}
	return result
}

func (this *EnumStruct) Datas() []EnumData {
	result := []EnumData{}
	for key, value := range this.names {
		singleKey, _ := strconv.Atoi(key)
		result = append(result, EnumData{
			Id:   singleKey,
			Name: value,
		})
	}
	return result
}

func (this *EnumStruct) Keys() []int {
	result := []int{}
	for key, _ := range this.names {
		singleResult, _ := strconv.Atoi(key)
		result = append(
			result,
			singleResult,
		)
	}
	return result
}

func (this *EnumStruct) Values() []string {
	result := []string{}
	for _, value := range this.names {
		result = append(
			result,
			value,
		)
	}
	return result
}

type EnumDataString struct {
	Id   string
	Name string
}

type EnumStructString struct {
	names map[string]string
}

func InitEnumStructString(this interface{}) {
	enumInfo := reflect.TypeOf(this).Elem()
	enumValue := reflect.ValueOf(this)
	result := enumValue.Elem().FieldByName("EnumStructString").Addr().Interface().(*EnumStructString)
	result.names = map[string]string{}

	for i := 0; i != enumInfo.NumField(); i++ {
		singleField := enumInfo.Field(i)

		singleFieldTag := singleField.Tag.Get("enum")
		singleFieldTagArray := strings.Split(singleFieldTag, ",")
		if len(singleFieldTagArray) != 2 {
			continue
		}

		singleFieldTagValue := singleFieldTagArray[0]
		singleFieldTagSeeName := singleFieldTagArray[1]

		result.names[singleFieldTagValue] = singleFieldTagSeeName
		enumValue.Elem().Field(i).SetString(singleFieldTagValue)
	}
}

func (this *EnumStructString) Names() map[string]string {
	return this.names
}

func (this *EnumStructString) Entrys() map[string]string {
	return this.names
}

func (this *EnumStructString) Datas() []EnumDataString {
	result := []EnumDataString{}
	for key, value := range this.names {
		result = append(result, EnumDataString{
			Id:   key,
			Name: value,
		})
	}
	return result
}

func (this *EnumStructString) Keys() []string {
	result := []string{}
	for key, _ := range this.names {
		result = append(
			result,
			key,
		)
	}
	return result
}

func (this *EnumStructString) Values() []string {
	result := []string{}
	for _, value := range this.names {
		result = append(
			result,
			value,
		)
	}
	return result
}
