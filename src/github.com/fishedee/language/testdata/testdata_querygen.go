package testdata

import (
	"github.com/fishedee/language"
	"time"
)

func queryColumnMap_1ece722a7d33673f2c49b8a06f99351e3dcce19f(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make(map[string]ContentType, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single.Name] = single
	}
	return result
}

func queryColumnMap_1f846e6386c993be8df7226e0aaeabe7a044ff66(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[string]int, len(dataIn))
	result := make(map[string][]ContentType, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Name
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		result[bufferData[kbegin].Name] = bufferData[kbegin:k]

	}
	return result
}

func queryColumnMap_22eaf4916410913f9870a2ed4e08c0d1838820a5(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make(map[bool]ContentType, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single.Ok] = single
	}
	return result
}

func queryColumnMap_3923b792e276005e09637544ecb3aec8be870f41(data interface{}, column string) interface{} {
	dataIn := data.([]string)
	result := make(map[string]string, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single] = single
	}
	return result
}

func queryColumnMap_3f5ce95fa1ae6f14be71a24adfceab807a606a7e(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[float64]int, len(dataIn))
	result := make(map[float64][]ContentType, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].CardMoney
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		result[bufferData[kbegin].CardMoney] = bufferData[kbegin:k]

	}
	return result
}

func queryColumnMap_4e2633c1c32af12dd2f1925b0eefbd4306b79b81(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[bool]int, len(dataIn))
	result := make(map[bool][]ContentType, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Ok
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		result[bufferData[kbegin].Ok] = bufferData[kbegin:k]

	}
	return result
}

func queryColumnMap_51b6a74d7266264feb8e264faefee78e0f7f7d84(data interface{}, column string) interface{} {
	dataIn := data.([]QueryInnerStruct2)
	bufferData := make([]QueryInnerStruct2, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	result := make(map[int][]QueryInnerStruct2, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].QueryInnerStruct.MM
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		result[bufferData[kbegin].QueryInnerStruct.MM] = bufferData[kbegin:k]

	}
	return result
}

func queryColumnMap_72e81d947328a78aeacab845c368290499b452b5(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make(map[string]ContentType, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single.Name] = single
	}
	return result
}

func queryColumnMap_91dacd60e87431951940b4b4c51428e7c1e5c1f2(data interface{}, column string) interface{} {
	dataIn := data.([]int)
	result := make(map[int]int, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single] = single
	}
	return result
}

func queryColumnMap_969c70c98e62ab1531c6ef65fb02a3ab0a60c63b(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make(map[float64]ContentType, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single.CardMoney] = single
	}
	return result
}

func queryColumnMap_b0ca9534ec19e54917839f5297f0d44469464db0(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[string]int, len(dataIn))
	result := make(map[string][]ContentType, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Name
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		result[bufferData[kbegin].Name] = bufferData[kbegin:k]

	}
	return result
}

func queryColumnMap_b67feeb79111bb008d14b9e3dc759d95b235a46d(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make(map[float32]ContentType, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single.Money] = single
	}
	return result
}

func queryColumnMap_b6c3dfa21cdfe36f75edf4ec12edc2572b8ab114(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[float32]int, len(dataIn))
	result := make(map[float32][]ContentType, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Money
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		result[bufferData[kbegin].Money] = bufferData[kbegin:k]

	}
	return result
}

func queryColumnMap_cfb8dbb56fd1a519299c7fc07b22cd07c8d8509d(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	result := make(map[int][]ContentType, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Age
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		result[bufferData[kbegin].Age] = bufferData[kbegin:k]

	}
	return result
}

func queryColumnMap_e4c08191085eb833a6da38ddffdcd489b977d915(data interface{}, column string) interface{} {
	dataIn := data.([]QueryInnerStruct2)
	result := make(map[int]QueryInnerStruct2, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single.QueryInnerStruct.MM] = single
	}
	return result
}

func queryColumnMap_f0c452a889049efb4433251e366aa1e97e11c451(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make(map[int]ContentType, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single.Age] = single
	}
	return result
}

func queryColumn_1897b26b0527e15e3bd79e174b6d77d289e1a65c(data interface{}, column string) interface{} {
	dataIn := data.([]QueryInnerStruct2)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.MM
	}
	return result
}

func queryColumn_1ece722a7d33673f2c49b8a06f99351e3dcce19f(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]string, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.Name
	}
	return result
}

func queryColumn_22eaf4916410913f9870a2ed4e08c0d1838820a5(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]bool, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.Ok
	}
	return result
}

func queryColumn_2db77a0314db25815cc49612c3dfc4e4dc7a7df5(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.Age
	}
	return result
}

func queryColumn_363289941cdcae601d6acca61118f49d8b7bb461(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]float32, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.Money
	}
	return result
}

func queryColumn_3923b792e276005e09637544ecb3aec8be870f41(data interface{}, column string) interface{} {
	dataIn := data.([]string)
	result := make([]string, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single
	}
	return result
}

func queryColumn_72e81d947328a78aeacab845c368290499b452b5(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]string, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.Name
	}
	return result
}

func queryColumn_796b51b89a4bdecef2a8ec13c12f83b7a8d99550(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]float64, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.CardMoney
	}
	return result
}

func queryColumn_91dacd60e87431951940b4b4c51428e7c1e5c1f2(data interface{}, column string) interface{} {
	dataIn := data.([]int)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single
	}
	return result
}

func queryColumn_969c70c98e62ab1531c6ef65fb02a3ab0a60c63b(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]float64, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.CardMoney
	}
	return result
}

func queryColumn_b67feeb79111bb008d14b9e3dc759d95b235a46d(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]float32, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.Money
	}
	return result
}

func queryColumn_cb6e7dd22244a2efe95bc35135515b5dfe67e699(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]float64, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.CardMoney
	}
	return result
}

func queryColumn_e4c08191085eb833a6da38ddffdcd489b977d915(data interface{}, column string) interface{} {
	dataIn := data.([]QueryInnerStruct2)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.QueryInnerStruct.MM
	}
	return result
}

func queryColumn_f0c452a889049efb4433251e366aa1e97e11c451(data interface{}, column string) interface{} {
	dataIn := data.([]ContentType)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.Age
	}
	return result
}

func queryCombine_498f0875f7fd905de41f9a7faf01928fedf1d791(leftData interface{}, rightData interface{}, combineFunctor interface{}) interface{} {
	leftDataIn := leftData.([]ContentType)
	rightDataIn := rightData.([]ContentType)
	combineFunctorIn := combineFunctor.(func(ContentType, ContentType) ContentType)
	newData := make([]ContentType, len(leftDataIn), len(leftDataIn))

	for i := 0; i != len(leftDataIn); i++ {
		newData[i] = combineFunctorIn(leftDataIn[i], rightDataIn[i])
	}
	return newData
}

func queryCombine_e5205a418275994376962bfefa27e0c47b33c158(leftData interface{}, rightData interface{}, combineFunctor interface{}) interface{} {
	leftDataIn := leftData.([]ContentType)
	rightDataIn := rightData.([]int)
	combineFunctorIn := combineFunctor.(func(ContentType, int) ContentType)
	newData := make([]ContentType, len(leftDataIn), len(leftDataIn))

	for i := 0; i != len(leftDataIn); i++ {
		newData[i] = combineFunctorIn(leftDataIn[i], rightDataIn[i])
	}
	return newData
}

func queryGroup_241b1bcf7b38898b817a5fa0b9b1481b9de8bf64(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]ContentType) []float64)
	result := make([]float64, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Age
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single...)

	}
	return result
}

func queryGroup_28702cb820985189e203e1e182f040a09a15b748(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[string]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]ContentType) float32)
	result := make([]float32, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Name
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single)

	}
	return result
}

func queryGroup_404222111e316051b0635679da2db87e058dc0a7(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[time.Time]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]ContentType) []ContentType)
	result := make([]ContentType, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Register
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single...)

	}
	return result
}

func queryGroup_4d113904af3cb457527f8538140bf878898b11b0(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]ContentType) float64)
	result := make([]float64, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Age
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single)

	}
	return result
}

func queryGroup_7959aac2ba701c92b02938af82c21599cbf58c3d(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]int)
	bufferData := make([]int, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]int) int)
	result := make([]int, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i]
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single)

	}
	return result
}

func queryGroup_8b187885bc010254ab50a591ab9b1075b8c9a748(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]string)
	bufferData := make([]string, len(dataIn), len(dataIn))
	mapData := make(map[string]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]string) ContentType)
	result := make([]ContentType, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i]
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single)

	}
	return result
}

func queryGroup_a797b7d91361679a9f983b26327d7d2c1df03a26(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[string]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]ContentType) []ContentType)
	result := make([]ContentType, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Name
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single...)

	}
	return result
}

func queryGroup_ac1e1e41612f45db132d8e78b5eb4487fa5e320c(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[bool]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]ContentType) []ContentType)
	result := make([]ContentType, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Ok
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single...)

	}
	return result
}

func queryGroup_b8a77a0f693628d67f0b0bffcd4447f08e082e71(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]QueryInnerStruct2)
	bufferData := make([]QueryInnerStruct2, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]QueryInnerStruct2) []QueryInnerStruct2)
	result := make([]QueryInnerStruct2, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].QueryInnerStruct.MM
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single...)

	}
	return result
}

func queryGroup_bfbe2ed6a187c11ce8bd7eeaa5d27a22c9d3688a(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[time.Time]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]ContentType) int)
	result := make([]int, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Register
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single)

	}
	return result
}

func queryGroup_ec4ec5fc86ab566b92e3d11f3f681c9e496ce42d(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[bool]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]ContentType) []ContentType)
	result := make([]ContentType, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Ok
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single...)

	}
	return result
}

func queryGroup_f092b4ddbe0a6336f3fd4da0c7f8bbd534ee5373(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	bufferData := make([]ContentType, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]ContentType) []ContentType)
	result := make([]ContentType, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].Age
		lastIndex, isExist := mapData[single]
		if isExist == true {
			nextData[lastIndex] = i
		}
		nextData[i] = -1
		mapData[single] = i
	}
	k := 0
	for i := 0; i != length; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData[k] = dataIn[j]
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData[k] = dataIn[j]
		k++
		nextData[j] = 0

		single := groupFunctorIn(bufferData[kbegin:k])
		result = append(result, single...)

	}
	return result
}

func queryJoin_105ea13485eb95132e906e4902f07d85a10d7f9a(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]UserType)
	joinFunctorIn := joinFunctor.(func(UserType, UserType) UserType)
	result := make([]UserType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := UserType{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[bool]int, len(rightDataIn))
	mapDataFirst := make(map[bool]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].Ok
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.Ok
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_188f1f3cb965deb4519198627cc5d9044947832e(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]string)
	rightDataIn := rightData.([]ContentType2)
	joinFunctorIn := joinFunctor.(func(string, ContentType2) ContentType2)
	result := make([]ContentType2, 0, len(leftDataIn))

	emptyLeftData := ""
	emptyRightData := ContentType2{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[string]int, len(rightDataIn))
	mapDataFirst := make(map[string]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].UserName
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_1f0a245beb92a9c9cfaae1eba65bb739221d3473(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]string)
	rightDataIn := rightData.([]UserType)
	joinFunctorIn := joinFunctor.(func(string, UserType) UserType)
	result := make([]UserType, 0, len(leftDataIn))

	emptyLeftData := ""
	emptyRightData := UserType{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[string]int, len(rightDataIn))
	mapDataFirst := make(map[string]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].Name
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_2b411fd124b38f38bbf505ae955b6a638a5fbdee(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]ContentType2)
	joinFunctorIn := joinFunctor.(func(UserType, ContentType2) resultType)
	result := make([]resultType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := ContentType2{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[string]int, len(rightDataIn))
	mapDataFirst := make(map[string]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].UserName
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.Name
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_4636c90d2e369aff71bf0d17243006ecc61cc1d6(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]UserType)
	joinFunctorIn := joinFunctor.(func(UserType, UserType) UserType)
	result := make([]UserType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := UserType{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[float64]int, len(rightDataIn))
	mapDataFirst := make(map[float64]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].Money
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.CardMoney
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_5c416f1d7f32d01c9845fce598463dd17d8ed325(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]UserType)
	joinFunctorIn := joinFunctor.(func(UserType, UserType) UserType)
	result := make([]UserType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := UserType{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[time.Time]int, len(rightDataIn))
	mapDataFirst := make(map[time.Time]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].Register
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.Register
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_5f1db64a7741403db18e58f5098e2d5545e9a483(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]ContentType2)
	joinFunctorIn := joinFunctor.(func(UserType, ContentType2) resultType)
	result := make([]resultType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := ContentType2{}
	joinPlace = "inner"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[string]int, len(rightDataIn))
	mapDataFirst := make(map[string]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].UserName
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.Name
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_702e3375a4aa64fdc9ae63bf9aba116949c59030(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]UserType)
	joinFunctorIn := joinFunctor.(func(UserType, UserType) UserType)
	result := make([]UserType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := UserType{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[float64]int, len(rightDataIn))
	mapDataFirst := make(map[float64]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].Money
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.Money
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_767dcfaf94351c32b545f391f13af41d570c4309(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]ExtendType)
	rightDataIn := rightData.([]ExtendType)
	joinFunctorIn := joinFunctor.(func(ExtendType, ExtendType) ExtendType)
	result := make([]ExtendType, 0, len(leftDataIn))

	emptyLeftData := ExtendType{}
	emptyRightData := ExtendType{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[int]int, len(rightDataIn))
	mapDataFirst := make(map[int]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].ContentID
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.ContentID
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_86cd12dcff69e784b920db7e61a83b4633a89290(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]ContentType2)
	joinFunctorIn := joinFunctor.(func(UserType, ContentType2) resultType)
	result := make([]resultType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := ContentType2{}
	joinPlace = "outer"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[string]int, len(rightDataIn))
	mapDataFirst := make(map[string]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].UserName
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.Name
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_936894a3c4a35ea014b97cb193c02b5b39b196e7(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]int)
	rightDataIn := rightData.([]ExtendType)
	joinFunctorIn := joinFunctor.(func(int, ExtendType) ExtendType)
	result := make([]ExtendType, 0, len(leftDataIn))

	emptyLeftData := 0
	emptyRightData := ExtendType{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[int]int, len(rightDataIn))
	mapDataFirst := make(map[int]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].ContentID
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_9d76a75216f27a6e35512b0a2658192b4ecfccdf(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]ContentType2)
	joinFunctorIn := joinFunctor.(func(UserType, ContentType2) resultType)
	result := make([]resultType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := ContentType2{}
	joinPlace = "right"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[string]int, len(rightDataIn))
	mapDataFirst := make(map[string]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].UserName
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.Name
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_9ec6b004b7ba349caa4deb7e4381212eb23da4ce(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]UserType)
	joinFunctorIn := joinFunctor.(func(UserType, UserType) UserType)
	result := make([]UserType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := UserType{}
	joinPlace = "right"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[int]int, len(rightDataIn))
	mapDataFirst := make(map[int]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].Age
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.Age
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_ac0bbb35db9996020984ef030135eee9ed90099f(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]UserType)
	rightDataIn := rightData.([]UserType)
	joinFunctorIn := joinFunctor.(func(UserType, UserType) UserType)
	result := make([]UserType, 0, len(leftDataIn))

	emptyLeftData := UserType{}
	emptyRightData := UserType{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[string]int, len(rightDataIn))
	mapDataFirst := make(map[string]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].Name
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.Name
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func queryJoin_e0c990030119f0a636c5de134b81ff3dcba6a4ef(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]QueryInnerStruct2)
	rightDataIn := rightData.([]QueryInnerStruct2)
	joinFunctorIn := joinFunctor.(func(QueryInnerStruct2, QueryInnerStruct2) QueryInnerStruct2)
	result := make([]QueryInnerStruct2, 0, len(leftDataIn))

	emptyLeftData := QueryInnerStruct2{}
	emptyRightData := QueryInnerStruct2{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[int]int, len(rightDataIn))
	mapDataFirst := make(map[int]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].QueryInnerStruct.MM
		lastIndex, isExist := mapDataNext[fieldValue]
		if isExist {
			nextData[lastIndex] = i
		} else {
			mapDataFirst[fieldValue] = i
		}
		nextData[i] = -1
		mapDataNext[fieldValue] = i
	}

	rightHaveJoin := make([]bool, len(rightDataIn), len(rightDataIn))
	for i := 0; i != len(leftDataIn); i++ {
		leftValue := leftDataIn[i]
		fieldValue := leftValue.QueryInnerStruct.MM
		rightIndex, isExist := mapDataFirst[fieldValue]
		if isExist {
			//找到右值
			j := rightIndex
			for nextData[j] != -1 {
				singleResult := joinFunctorIn(leftValue, rightDataIn[j])
				result = append(result, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFunctorIn(leftValue, rightDataIn[j])
			result = append(result, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFunctorIn(leftValue, emptyRightData)
				result = append(result, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != len(rightDataIn); j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFunctorIn(emptyLeftData, rightDataIn[j])
			result = append(result, singleResult)
		}
	}
	return result
}

func querySelect_349dab3622f6aba4c841265d75047d4846bd50f1(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	selectFunctorIn := selectFunctor.(func(ContentType) time.Time)
	result := make([]time.Time, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySelect_751a48272555f411640654a568234591de2c989d(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	selectFunctorIn := selectFunctor.(func(ContentType) map[string]int)
	result := make([]map[string]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySelect_7e8cdb175849e0310239d97a0d85c76f82097851(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	selectFunctorIn := selectFunctor.(func(ContentType) bool)
	result := make([]bool, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySelect_83677541a15193737a6e60ef6b644b717c3213c6(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	selectFunctorIn := selectFunctor.(func(ContentType) ContentType)
	result := make([]ContentType, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySelect_a52280403b46e09c6036e28f196e54d5028d6698(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	selectFunctorIn := selectFunctor.(func(ContentType) string)
	result := make([]string, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySelect_b2b356919faffda6b9c4fefdf38cee4e918f99cd(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	selectFunctorIn := selectFunctor.(func(ContentType) int)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySelect_b4985d2137d7ebcc21846adc496bf8b366e536ed(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	selectFunctorIn := selectFunctor.(func(ContentType) float32)
	result := make([]float32, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySelect_da391223926715b44206be57f539447211d3353f(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	selectFunctorIn := selectFunctor.(func(ContentType) float64)
	result := make([]float64, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySort_01e20a2d98037c6228fe46a9f25b6f3afa94cc6f(data interface{}, sortType string) interface{} {
	dataIn := data.([]ContentType)
	newData := make([]ContentType, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].CardMoney < newData[j].CardMoney {
			return -1
		} else if newData[i].CardMoney > newData[j].CardMoney {
			return 1
		}

		if newData[i].Register.Before(newData[j].Register) {
			return 1
		} else if newData[i].Register.After(newData[j].Register) {
			return -1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_21a89d848f4663606bb436706b2468766a4bd198(data interface{}, sortType string) interface{} {
	dataIn := data.([]language.Decimal)
	newData := make([]language.Decimal, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		{
			tempDecimalCmp := newData[i].Cmp(newData[j])
			if tempDecimalCmp < 0 {
				return 1
			} else if tempDecimalCmp > 0 {
				return -1
			}
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_3984d355d2e723bca314b4298b8199ea48d02fa4(data interface{}, sortType string) interface{} {
	dataIn := data.([]Student)
	newData := make([]Student, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].Name < newData[j].Name {
			return -1
		} else if newData[i].Name > newData[j].Name {
			return 1
		}

		{
			tempDecimalCmp := newData[i].Score.Cmp(newData[j].Score)
			if tempDecimalCmp < 0 {
				return 1
			} else if tempDecimalCmp > 0 {
				return -1
			}
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_3e5610abb58143b11a227c0338ad79745f949d87(data interface{}, sortType string) interface{} {
	dataIn := data.([]Student)
	newData := make([]Student, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		{
			tempDecimalCmp := newData[i].Score.Cmp(newData[j].Score)
			if tempDecimalCmp < 0 {
				return -1
			} else if tempDecimalCmp > 0 {
				return 1
			}
		}

		{
			tempDecimalCmp := newData[i].Score2.Cmp(newData[j].Score2)
			if tempDecimalCmp < 0 {
				return 1
			} else if tempDecimalCmp > 0 {
				return -1
			}
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_69aeea64e4ceb06cec2898980c86869d76d5ded1(data interface{}, sortType string) interface{} {
	dataIn := data.([]ContentType)
	newData := make([]ContentType, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].Name < newData[j].Name {
			return 1
		} else if newData[i].Name > newData[j].Name {
			return -1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_6a05af9df6bce94634838e76503ac326a38fafb9(data interface{}, sortType string) interface{} {
	dataIn := data.([]ContentType)
	newData := make([]ContentType, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].Ok == false && newData[j].Ok == true {
			return 1
		} else if newData[i].Ok == true && newData[j].Ok == false {
			return -1
		}

		if newData[i].Name < newData[j].Name {
			return -1
		} else if newData[i].Name > newData[j].Name {
			return 1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_74654e8b45593005ef783b89255269f7c6ecc39b(data interface{}, sortType string) interface{} {
	dataIn := data.([]int)
	newData := make([]int, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i] < newData[j] {
			return -1
		} else if newData[i] > newData[j] {
			return 1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_ad528a611518b9a3daaf165009306aeabc073462(data interface{}, sortType string) interface{} {
	dataIn := data.([]ContentType)
	newData := make([]ContentType, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].Money < newData[j].Money {
			return 1
		} else if newData[i].Money > newData[j].Money {
			return -1
		}

		if newData[i].Age < newData[j].Age {
			return -1
		} else if newData[i].Age > newData[j].Age {
			return 1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_af891d058d5a2e0a3ac4b4b291ae9bb959364795(data interface{}, sortType string) interface{} {
	dataIn := data.([]int)
	newData := make([]int, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i] < newData[j] {
			return 1
		} else if newData[i] > newData[j] {
			return -1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_b62d6f59f1fd62be17dc174c0ae49226de3490b8(data interface{}, sortType string) interface{} {
	dataIn := data.([]ContentType)
	newData := make([]ContentType, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].Age < newData[j].Age {
			return 1
		} else if newData[i].Age > newData[j].Age {
			return -1
		}

		if newData[i].Ok == false && newData[j].Ok == true {
			return 1
		} else if newData[i].Ok == true && newData[j].Ok == false {
			return -1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_c27d7e1f10e5b181192befeb719a29c23de75a68(data interface{}, sortType string) interface{} {
	dataIn := data.([]ContentType)
	newData := make([]ContentType, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].Money < newData[j].Money {
			return -1
		} else if newData[i].Money > newData[j].Money {
			return 1
		}

		if newData[i].Register.Before(newData[j].Register) {
			return 1
		} else if newData[i].Register.After(newData[j].Register) {
			return -1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_d2e998dae6686adadd86175e119e0ab490e8959f(data interface{}, sortType string) interface{} {
	dataIn := data.([]ContentType)
	newData := make([]ContentType, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].Name < newData[j].Name {
			return -1
		} else if newData[i].Name > newData[j].Name {
			return 1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_db2d6656b4568f16581fa7db7f72b2a0f4b36411(data interface{}, sortType string) interface{} {
	dataIn := data.([]QueryInnerStruct2)
	newData := make([]QueryInnerStruct2, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].MM < newData[j].MM {
			return 1
		} else if newData[i].MM > newData[j].MM {
			return -1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_e72c6d15cf7f64d6195886f7666e574b033ef607(data interface{}, sortType string) interface{} {
	dataIn := data.([]QueryInnerStruct2)
	newData := make([]QueryInnerStruct2, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].QueryInnerStruct.MM < newData[j].QueryInnerStruct.MM {
			return -1
		} else if newData[i].QueryInnerStruct.MM > newData[j].QueryInnerStruct.MM {
			return 1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_fcc91fba4bd2e004c8bcd771f04272f9b64df696(data interface{}, sortType string) interface{} {
	dataIn := data.([]ContentType)
	newData := make([]ContentType, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].Money < newData[j].Money {
			return 1
		} else if newData[i].Money > newData[j].Money {
			return -1
		}

		if newData[i].Age < newData[j].Age {
			return -1
		} else if newData[i].Age > newData[j].Age {
			return 1
		}

		if newData[i].Name < newData[j].Name {
			return 1
		} else if newData[i].Name > newData[j].Name {
			return -1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func queryWhere_7e8cdb175849e0310239d97a0d85c76f82097851(data interface{}, whereFunctor interface{}) interface{} {
	dataIn := data.([]ContentType)
	whereFunctorIn := whereFunctor.(func(ContentType) bool)
	result := make([]ContentType, 0, len(dataIn))

	for _, single := range dataIn {
		shouldStay := whereFunctorIn(single)
		if shouldStay == true {
			result = append(result, single)
		}
	}
	return result
}

func init() {

	language.QueryColumnMapMacroRegister([]ContentType{}, " Name ", queryColumnMap_1ece722a7d33673f2c49b8a06f99351e3dcce19f)

	language.QueryColumnMapMacroRegister([]ContentType{}, " []Name ", queryColumnMap_1f846e6386c993be8df7226e0aaeabe7a044ff66)

	language.QueryColumnMapMacroRegister([]ContentType{}, "Ok        ", queryColumnMap_22eaf4916410913f9870a2ed4e08c0d1838820a5)

	language.QueryColumnMapMacroRegister([]string{}, " . ", queryColumnMap_3923b792e276005e09637544ecb3aec8be870f41)

	language.QueryColumnMapMacroRegister([]ContentType{}, "    []CardMoney", queryColumnMap_3f5ce95fa1ae6f14be71a24adfceab807a606a7e)

	language.QueryColumnMapMacroRegister([]ContentType{}, "[]Ok        ", queryColumnMap_4e2633c1c32af12dd2f1925b0eefbd4306b79b81)

	language.QueryColumnMapMacroRegister([]QueryInnerStruct2{}, "[]QueryInnerStruct.MM", queryColumnMap_51b6a74d7266264feb8e264faefee78e0f7f7d84)

	language.QueryColumnMapMacroRegister([]ContentType{}, "     Name         ", queryColumnMap_72e81d947328a78aeacab845c368290499b452b5)

	language.QueryColumnMapMacroRegister([]int{}, " . ", queryColumnMap_91dacd60e87431951940b4b4c51428e7c1e5c1f2)

	language.QueryColumnMapMacroRegister([]ContentType{}, "    CardMoney", queryColumnMap_969c70c98e62ab1531c6ef65fb02a3ab0a60c63b)

	language.QueryColumnMapMacroRegister([]ContentType{}, "     [] Name         ", queryColumnMap_b0ca9534ec19e54917839f5297f0d44469464db0)

	language.QueryColumnMapMacroRegister([]ContentType{}, "    Money  ", queryColumnMap_b67feeb79111bb008d14b9e3dc759d95b235a46d)

	language.QueryColumnMapMacroRegister([]ContentType{}, "    []Money  ", queryColumnMap_b6c3dfa21cdfe36f75edf4ec12edc2572b8ab114)

	language.QueryColumnMapMacroRegister([]ContentType{}, "[]Age        ", queryColumnMap_cfb8dbb56fd1a519299c7fc07b22cd07c8d8509d)

	language.QueryColumnMapMacroRegister([]QueryInnerStruct2{}, "QueryInnerStruct.MM", queryColumnMap_e4c08191085eb833a6da38ddffdcd489b977d915)

	language.QueryColumnMapMacroRegister([]ContentType{}, "Age        ", queryColumnMap_f0c452a889049efb4433251e366aa1e97e11c451)

	language.QueryColumnMacroRegister([]QueryInnerStruct2{}, "  MM  ", queryColumn_1897b26b0527e15e3bd79e174b6d77d289e1a65c)

	language.QueryColumnMacroRegister([]ContentType{}, " Name ", queryColumn_1ece722a7d33673f2c49b8a06f99351e3dcce19f)

	language.QueryColumnMacroRegister([]ContentType{}, "Ok        ", queryColumn_22eaf4916410913f9870a2ed4e08c0d1838820a5)

	language.QueryColumnMacroRegister([]ContentType{}, "  Age  ", queryColumn_2db77a0314db25815cc49612c3dfc4e4dc7a7df5)

	language.QueryColumnMacroRegister([]ContentType{}, "  Money  ", queryColumn_363289941cdcae601d6acca61118f49d8b7bb461)

	language.QueryColumnMacroRegister([]string{}, " . ", queryColumn_3923b792e276005e09637544ecb3aec8be870f41)

	language.QueryColumnMacroRegister([]ContentType{}, "     Name         ", queryColumn_72e81d947328a78aeacab845c368290499b452b5)

	language.QueryColumnMacroRegister([]ContentType{}, "CardMoney  ", queryColumn_796b51b89a4bdecef2a8ec13c12f83b7a8d99550)

	language.QueryColumnMacroRegister([]int{}, " . ", queryColumn_91dacd60e87431951940b4b4c51428e7c1e5c1f2)

	language.QueryColumnMacroRegister([]ContentType{}, "    CardMoney", queryColumn_969c70c98e62ab1531c6ef65fb02a3ab0a60c63b)

	language.QueryColumnMacroRegister([]ContentType{}, "    Money  ", queryColumn_b67feeb79111bb008d14b9e3dc759d95b235a46d)

	language.QueryColumnMacroRegister([]ContentType{}, "  CardMoney  ", queryColumn_cb6e7dd22244a2efe95bc35135515b5dfe67e699)

	language.QueryColumnMacroRegister([]QueryInnerStruct2{}, "QueryInnerStruct.MM", queryColumn_e4c08191085eb833a6da38ddffdcd489b977d915)

	language.QueryColumnMacroRegister([]ContentType{}, "Age        ", queryColumn_f0c452a889049efb4433251e366aa1e97e11c451)

	language.QueryCombineMacroRegister([]ContentType{}, []ContentType{}, (func(ContentType, ContentType) ContentType)(nil), queryCombine_498f0875f7fd905de41f9a7faf01928fedf1d791)

	language.QueryCombineMacroRegister([]ContentType{}, []int{}, (func(ContentType, int) ContentType)(nil), queryCombine_e5205a418275994376962bfefa27e0c47b33c158)

	language.QueryGroupMacroRegister([]ContentType{}, " Age ", (func([]ContentType) []float64)(nil), queryGroup_241b1bcf7b38898b817a5fa0b9b1481b9de8bf64)

	language.QueryGroupMacroRegister([]ContentType{}, "Name", (func([]ContentType) float32)(nil), queryGroup_28702cb820985189e203e1e182f040a09a15b748)

	language.QueryGroupMacroRegister([]ContentType{}, "Register ", (func([]ContentType) []ContentType)(nil), queryGroup_404222111e316051b0635679da2db87e058dc0a7)

	language.QueryGroupMacroRegister([]ContentType{}, " Age ", (func([]ContentType) float64)(nil), queryGroup_4d113904af3cb457527f8538140bf878898b11b0)

	language.QueryGroupMacroRegister([]int{}, ".", (func([]int) int)(nil), queryGroup_7959aac2ba701c92b02938af82c21599cbf58c3d)

	language.QueryGroupMacroRegister([]string{}, ".", (func([]string) ContentType)(nil), queryGroup_8b187885bc010254ab50a591ab9b1075b8c9a748)

	language.QueryGroupMacroRegister([]ContentType{}, "Name", (func([]ContentType) []ContentType)(nil), queryGroup_a797b7d91361679a9f983b26327d7d2c1df03a26)

	language.QueryGroupMacroRegister([]ContentType{}, " Ok ", (func([]ContentType) []ContentType)(nil), queryGroup_ac1e1e41612f45db132d8e78b5eb4487fa5e320c)

	language.QueryGroupMacroRegister([]QueryInnerStruct2{}, "QueryInnerStruct.MM", (func([]QueryInnerStruct2) []QueryInnerStruct2)(nil), queryGroup_b8a77a0f693628d67f0b0bffcd4447f08e082e71)

	language.QueryGroupMacroRegister([]ContentType{}, "Register ", (func([]ContentType) int)(nil), queryGroup_bfbe2ed6a187c11ce8bd7eeaa5d27a22c9d3688a)

	language.QueryGroupMacroRegister([]ContentType{}, "Ok", (func([]ContentType) []ContentType)(nil), queryGroup_ec4ec5fc86ab566b92e3d11f3f681c9e496ce42d)

	language.QueryGroupMacroRegister([]ContentType{}, " Age ", (func([]ContentType) []ContentType)(nil), queryGroup_f092b4ddbe0a6336f3fd4da0c7f8bbd534ee5373)

	language.QueryJoinMacroRegister([]UserType{}, []UserType{}, "left", "Ok  =  Ok", (func(UserType, UserType) UserType)(nil), queryJoin_105ea13485eb95132e906e4902f07d85a10d7f9a)

	language.QueryJoinMacroRegister([]string{}, []ContentType2{}, "left", "  .  =  UserName ", (func(string, ContentType2) ContentType2)(nil), queryJoin_188f1f3cb965deb4519198627cc5d9044947832e)

	language.QueryJoinMacroRegister([]string{}, []UserType{}, "left", " . = Name", (func(string, UserType) UserType)(nil), queryJoin_1f0a245beb92a9c9cfaae1eba65bb739221d3473)

	language.QueryJoinMacroRegister([]UserType{}, []ContentType2{}, "left", "  Name  =  UserName ", (func(UserType, ContentType2) resultType)(nil), queryJoin_2b411fd124b38f38bbf505ae955b6a638a5fbdee)

	language.QueryJoinMacroRegister([]UserType{}, []UserType{}, "left", " CardMoney = Money ", (func(UserType, UserType) UserType)(nil), queryJoin_4636c90d2e369aff71bf0d17243006ecc61cc1d6)

	language.QueryJoinMacroRegister([]UserType{}, []UserType{}, "left", " Register = Register ", (func(UserType, UserType) UserType)(nil), queryJoin_5c416f1d7f32d01c9845fce598463dd17d8ed325)

	language.QueryJoinMacroRegister([]UserType{}, []ContentType2{}, "inner", "  Name  =  UserName ", (func(UserType, ContentType2) resultType)(nil), queryJoin_5f1db64a7741403db18e58f5098e2d5545e9a483)

	language.QueryJoinMacroRegister([]UserType{}, []UserType{}, "left", " Money=Money ", (func(UserType, UserType) UserType)(nil), queryJoin_702e3375a4aa64fdc9ae63bf9aba116949c59030)

	language.QueryJoinMacroRegister([]ExtendType{}, []ExtendType{}, " left ", "  ContentID  =  ContentID ", (func(ExtendType, ExtendType) ExtendType)(nil), queryJoin_767dcfaf94351c32b545f391f13af41d570c4309)

	language.QueryJoinMacroRegister([]UserType{}, []ContentType2{}, "outer", "  Name  =  UserName ", (func(UserType, ContentType2) resultType)(nil), queryJoin_86cd12dcff69e784b920db7e61a83b4633a89290)

	language.QueryJoinMacroRegister([]int{}, []ExtendType{}, " left ", "  .  =  ContentID ", (func(int, ExtendType) ExtendType)(nil), queryJoin_936894a3c4a35ea014b97cb193c02b5b39b196e7)

	language.QueryJoinMacroRegister([]UserType{}, []ContentType2{}, "right", "  Name  =  UserName ", (func(UserType, ContentType2) resultType)(nil), queryJoin_9d76a75216f27a6e35512b0a2658192b4ecfccdf)

	language.QueryJoinMacroRegister([]UserType{}, []UserType{}, "right", "Age=Age", (func(UserType, UserType) UserType)(nil), queryJoin_9ec6b004b7ba349caa4deb7e4381212eb23da4ce)

	language.QueryJoinMacroRegister([]UserType{}, []UserType{}, " left ", "  Name  =  Name ", (func(UserType, UserType) UserType)(nil), queryJoin_ac0bbb35db9996020984ef030135eee9ed90099f)

	language.QueryJoinMacroRegister([]QueryInnerStruct2{}, []QueryInnerStruct2{}, "left", "QueryInnerStruct.MM = QueryInnerStruct.MM", (func(QueryInnerStruct2, QueryInnerStruct2) QueryInnerStruct2)(nil), queryJoin_e0c990030119f0a636c5de134b81ff3dcba6a4ef)

	language.QuerySelectMacroRegister([]ContentType{}, (func(ContentType) time.Time)(nil), querySelect_349dab3622f6aba4c841265d75047d4846bd50f1)

	language.QuerySelectMacroRegister([]ContentType{}, (func(ContentType) map[string]int)(nil), querySelect_751a48272555f411640654a568234591de2c989d)

	language.QuerySelectMacroRegister([]ContentType{}, (func(ContentType) bool)(nil), querySelect_7e8cdb175849e0310239d97a0d85c76f82097851)

	language.QuerySelectMacroRegister([]ContentType{}, (func(ContentType) ContentType)(nil), querySelect_83677541a15193737a6e60ef6b644b717c3213c6)

	language.QuerySelectMacroRegister([]ContentType{}, (func(ContentType) string)(nil), querySelect_a52280403b46e09c6036e28f196e54d5028d6698)

	language.QuerySelectMacroRegister([]ContentType{}, (func(ContentType) int)(nil), querySelect_b2b356919faffda6b9c4fefdf38cee4e918f99cd)

	language.QuerySelectMacroRegister([]ContentType{}, (func(ContentType) float32)(nil), querySelect_b4985d2137d7ebcc21846adc496bf8b366e536ed)

	language.QuerySelectMacroRegister([]ContentType{}, (func(ContentType) float64)(nil), querySelect_da391223926715b44206be57f539447211d3353f)

	language.QuerySortMacroRegister([]ContentType{}, "CardMoney,Register desc", querySort_01e20a2d98037c6228fe46a9f25b6f3afa94cc6f)

	language.QuerySortMacroRegister([]language.Decimal{}, ". desc", querySort_21a89d848f4663606bb436706b2468766a4bd198)

	language.QuerySortMacroRegister([]Student{}, "Name asc,Score desc", querySort_3984d355d2e723bca314b4298b8199ea48d02fa4)

	language.QuerySortMacroRegister([]Student{}, "Score asc,Score2 desc", querySort_3e5610abb58143b11a227c0338ad79745f949d87)

	language.QuerySortMacroRegister([]ContentType{}, "Name desc", querySort_69aeea64e4ceb06cec2898980c86869d76d5ded1)

	language.QuerySortMacroRegister([]ContentType{}, "Ok desc,Name", querySort_6a05af9df6bce94634838e76503ac326a38fafb9)

	language.QuerySortMacroRegister([]int{}, ". asc", querySort_74654e8b45593005ef783b89255269f7c6ecc39b)

	language.QuerySortMacroRegister([]ContentType{}, " Money desc,Age asc", querySort_ad528a611518b9a3daaf165009306aeabc073462)

	language.QuerySortMacroRegister([]int{}, ". desc", querySort_af891d058d5a2e0a3ac4b4b291ae9bb959364795)

	language.QuerySortMacroRegister([]ContentType{}, "Age desc,Ok desc", querySort_b62d6f59f1fd62be17dc174c0ae49226de3490b8)

	language.QuerySortMacroRegister([]ContentType{}, "Money,Register desc", querySort_c27d7e1f10e5b181192befeb719a29c23de75a68)

	language.QuerySortMacroRegister([]ContentType{}, "Name asc", querySort_d2e998dae6686adadd86175e119e0ab490e8959f)

	language.QuerySortMacroRegister([]QueryInnerStruct2{}, "MM desc", querySort_db2d6656b4568f16581fa7db7f72b2a0f4b36411)

	language.QuerySortMacroRegister([]QueryInnerStruct2{}, "QueryInnerStruct.MM asc", querySort_e72c6d15cf7f64d6195886f7666e574b033ef607)

	language.QuerySortMacroRegister([]ContentType{}, " Money desc,Age asc,Name desc", querySort_fcc91fba4bd2e004c8bcd771f04272f9b64df696)

	language.QueryWhereMacroRegister([]ContentType{}, (func(ContentType) bool)(nil), queryWhere_7e8cdb175849e0310239d97a0d85c76f82097851)

}
