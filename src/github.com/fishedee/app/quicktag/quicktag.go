package quicktag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type QuickTag struct {
	tag       string
	cache     sync.Map
	mutex     sync.Mutex
	basicType map[reflect.Kind]bool
}

func NewQuickTag(tag string) *QuickTag {
	basicTypeMap := map[reflect.Kind]bool{}

	basicType := []reflect.Kind{
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String,
		reflect.Chan,
		reflect.Func,
	}

	for _, kind := range basicType {
		basicTypeMap[kind] = true
	}
	return &QuickTag{
		tag:       tag,
		cache:     sync.Map{},
		mutex:     sync.Mutex{},
		basicType: basicTypeMap,
	}
}

func (this *QuickTag) GetTagType(t reflect.Type) reflect.Type {
	result, isExist := this.cache.Load(t)
	if isExist {
		return *(result.(*reflect.Type))
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	hasVisit := map[reflect.Type]bool{}
	return this.getTagTypeInner(hasVisit, t)
}

func (this *QuickTag) getTagTypeInner(hasVisit map[reflect.Type]bool, t reflect.Type) reflect.Type {
	timeType := reflect.TypeOf(time.Time{})
	rawMessageType := reflect.TypeOf(json.RawMessage{})

	cacheType, isExist := this.cache.Load(t)
	if isExist {
		return *(cacheType.(*reflect.Type))
	}

	hasVisit[t] = true
	var resultType reflect.Type

	kind := t.Kind()
	if this.basicType[kind] == true ||
		t == rawMessageType {
		resultType = t
	} else if kind == reflect.Ptr {
		tempType := this.getTagTypeInner(hasVisit, t.Elem())
		resultType = reflect.PtrTo(tempType)
	} else if kind == reflect.Array {
		tempType := this.getTagTypeInner(hasVisit, t.Elem())
		resultType = reflect.ArrayOf(t.Len(), tempType)
	} else if kind == reflect.Slice {
		tempType := this.getTagTypeInner(hasVisit, t.Elem())
		resultType = reflect.SliceOf(tempType)
	} else if kind == reflect.Map {
		tempType := this.getTagTypeInner(hasVisit, t.Elem())
		resultType = reflect.MapOf(t.Key(), tempType)
	} else if kind == reflect.Interface {
		panic(fmt.Sprintf("quick tag dosnot support interface %v", t))
	} else if t == timeType {
		resultType = reflect.TypeOf(myTime{})
	} else if kind == reflect.Struct {
		resultType = this.getStructType(hasVisit, t)
	} else {
		panic(fmt.Sprintf("unknown kind %v", kind))
	}

	hasVisit[t] = false
	this.cache.Store(t, &resultType)

	return resultType
}

func (this *QuickTag) getStructType(hasVisit map[reflect.Type]bool, t reflect.Type) reflect.Type {
	numField := t.NumField()
	newStructFields := []reflect.StructField{}

	for i := 0; i != numField; i++ {
		field := t.Field(i)
		if hasVisit[field.Type] == true {
			panic(fmt.Sprintf("quick tag can not support circle type %v->%v", t, field.Type))
		}
		fieldType := this.getTagTypeInner(hasVisit, field.Type)
		fieldName := field.Name
		fieldAnonymous := field.Anonymous
		fieldTag := this.getTag(field)
		newStructFields = append(newStructFields, reflect.StructField{
			Name:      fieldName,
			Type:      fieldType,
			Tag:       fieldTag,
			Anonymous: fieldAnonymous,
		})
	}
	return reflect.StructOf(newStructFields)
}

func (this *QuickTag) getTag(field reflect.StructField) reflect.StructTag {
	tag := field.Tag
	firstName := strings.ToLower(field.Name[0:1]) + field.Name[1:]
	secondSet := ""
	originInfo, hasOriginTag := tag.Lookup(this.tag)
	if hasOriginTag == true {
		originInfoList := strings.Split(originInfo, ",")
		if originInfoList[0] != "" {
			firstName = originInfoList[0]
		}
		if len(originInfoList) >= 2 && originInfoList[1] != "" {
			secondSet = originInfoList[1]
		}
	}
	var result = ""
	if secondSet == "" {
		result = fmt.Sprintf("%s:\"%s\"", this.tag, firstName)
	} else {
		result = fmt.Sprintf("%s:\"%s,%s\"", this.tag, firstName, secondSet)
	}
	return reflect.StructTag(result)
}

type emptyInterface struct {
	pt unsafe.Pointer
	pv unsafe.Pointer
}

func (this *QuickTag) pointerOfType(t reflect.Type) unsafe.Pointer {
	p := *(*emptyInterface)(unsafe.Pointer(&t))
	return p.pv
}

func (this *QuickTag) GetTagInstance(src interface{}) interface{} {
	if src == nil {
		return nil
	}
	srcType := reflect.TypeOf(src)
	eface := *(*emptyInterface)(unsafe.Pointer(&src))
	eface.pt = this.pointerOfType(this.GetTagType(srcType))
	dst := *(*interface{})(unsafe.Pointer(&eface))
	return dst
}
