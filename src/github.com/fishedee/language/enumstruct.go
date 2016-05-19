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
	datas []EnumData
}

func InitEnumStruct(this interface{}) {
	enumInfo := reflect.TypeOf(this).Elem()
	enumValue := reflect.ValueOf(this)
	result := enumValue.Elem().FieldByName("EnumStruct").Addr().Interface().(*EnumStruct)
	result.names = map[string]string{}
	result.datas = []EnumData{}

	for i := 0; i != enumInfo.NumField(); i++ {
		singleField := enumInfo.Field(i)
		if singleField.PkgPath != "" || singleField.Anonymous {
			continue
		}

		singleFieldName := singleField.Name
		singleFieldTag := singleField.Tag.Get("enum")
		singleFieldTagArray := strings.Split(singleFieldTag, ",")
		if len(singleFieldTagArray) != 2 {
			panic("invalid enum " + enumInfo.String() + ":" + singleFieldName)
		}

		singleFieldTagValue, err := strconv.Atoi(singleFieldTagArray[0])
		if err != nil {
			panic("invalid enum " + enumInfo.String() + ":" + singleFieldName)
		}
		singleFieldTagSeeName := singleFieldTagArray[1]
		if singleFieldTagSeeName == "" {
			panic("invalid enum " + enumInfo.String() + ":" + singleFieldName)
		}

		result.names[singleFieldTagArray[0]] = singleFieldTagSeeName
		result.datas = append(result.datas, EnumData{
			Id:   singleFieldTagValue,
			Name: singleFieldTagSeeName,
		})
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
	return this.datas
}

func (this *EnumStruct) Keys() []int {
	result := []int{}
	for _, singleEnum := range this.datas {
		result = append(
			result,
			singleEnum.Id,
		)
	}
	return result
}

func (this *EnumStruct) Values() []string {
	result := []string{}
	for _, singleEnum := range this.datas {
		result = append(
			result,
			singleEnum.Name,
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
	datas []EnumDataString
}

func InitEnumStructString(this interface{}) {
	enumInfo := reflect.TypeOf(this).Elem()
	enumValue := reflect.ValueOf(this)
	result := enumValue.Elem().FieldByName("EnumStructString").Addr().Interface().(*EnumStructString)
	result.names = map[string]string{}
	result.datas = []EnumDataString{}

	for i := 0; i != enumInfo.NumField(); i++ {
		singleField := enumInfo.Field(i)
		if singleField.PkgPath != "" || singleField.Anonymous {
			continue
		}
		singleFieldName := singleField.Name

		singleFieldTag := singleField.Tag.Get("enum")
		singleFieldTagArray := strings.Split(singleFieldTag, ",")
		if len(singleFieldTagArray) != 2 {
			panic("invalid enum " + enumInfo.String() + ":" + singleFieldName)
		}

		singleFieldTagValue := singleFieldTagArray[0]
		singleFieldTagSeeName := singleFieldTagArray[1]
		if singleFieldTagValue == "" || singleFieldTagSeeName == "" {
			panic("invalid enum " + enumInfo.String() + ":" + singleFieldName)
		}

		result.names[singleFieldTagValue] = singleFieldTagSeeName
		result.datas = append(result.datas, EnumDataString{
			Id:   singleFieldTagValue,
			Name: singleFieldTagSeeName,
		})

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
	return this.datas
}

func (this *EnumStructString) Keys() []string {
	result := []string{}
	for _, singleEnum := range this.datas {
		result = append(
			result,
			singleEnum.Id,
		)
	}
	return result
}

func (this *EnumStructString) Values() []string {
	result := []string{}
	for _, singleEnum := range this.datas {
		result = append(
			result,
			singleEnum.Name,
		)
	}
	return result
}
