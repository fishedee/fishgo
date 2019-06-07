package testdata

import (
	"github.com/fishedee/language"
	"github.com/fishedee/language/querygen/testdata/subtest"
	"time"
)

func queryColumnMap_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0(data interface{}, column string) interface{} {
	dataIn := data.([]User)
	result := make(map[int]User, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single.UserId] = single
	}
	return result
}

func queryColumnMap_904b262f8e2329ec73c320ca0e5ca82f14165586(data interface{}, column string) interface{} {
	dataIn := data.([]int)
	result := make(map[int]int, len(dataIn))

	for i := len(dataIn) - 1; i >= 0; i-- {
		single := dataIn[i]
		result[single] = single
	}
	return result
}

func queryColumnMap_99ac937255d05542023b7c6c1eba945d1b199327(data interface{}, column string) interface{} {
	dataIn := data.([]User)
	bufferData := make([]User, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	result := make(map[int][]User, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].UserId
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

		result[bufferData[kbegin].UserId] = bufferData[kbegin:k]

	}
	return result
}

func queryColumn_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0(data interface{}, column string) interface{} {
	dataIn := data.([]User)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.UserId
	}
	return result
}

func queryColumn_7ab1c12a87c03fba395eb892ababc4a58ad5dec9(data interface{}, column string) interface{} {
	dataIn := data.([]User)
	result := make([]User, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single
	}
	return result
}

func queryColumn_904b262f8e2329ec73c320ca0e5ca82f14165586(data interface{}, column string) interface{} {
	dataIn := data.([]int)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single
	}
	return result
}

func queryColumn_b60cd8d06e3e435a78322ac375157c99ea3ee15e(data interface{}, column string) interface{} {
	dataIn := data.([]subtest.Address)
	result := make([]string, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.City
	}
	return result
}

func queryCombine_91ed07637405c43f5a89f7fc9cbd76e4e62acc92(leftData interface{}, rightData interface{}, combineFunctor interface{}) interface{} {
	leftDataIn := leftData.([]int)
	rightDataIn := rightData.([]User)
	combineFunctorIn := combineFunctor.(func(int, User) User)
	newData := make([]User, len(leftDataIn), len(leftDataIn))

	for i := 0; i != len(leftDataIn); i++ {
		newData[i] = combineFunctorIn(leftDataIn[i], rightDataIn[i])
	}
	return newData
}

func queryCombine_d4934c4702a916f609947de20b2197cb0f02712e(leftData interface{}, rightData interface{}, combineFunctor interface{}) interface{} {
	leftDataIn := leftData.([]Admin)
	rightDataIn := rightData.([]User)
	combineFunctorIn := combineFunctor.(func(Admin, User) AdminUser)
	newData := make([]AdminUser, len(leftDataIn), len(leftDataIn))

	for i := 0; i != len(leftDataIn); i++ {
		newData[i] = combineFunctorIn(leftDataIn[i], rightDataIn[i])
	}
	return newData
}

func queryGroup_0b68844c3917648685a142b599f750918372f47c(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]int)
	bufferData := make([]int, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]int) Department)
	result := make([]Department, 0, len(dataIn))

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

func queryGroup_2a871c94537a18ef0ccfd971807ac9a23c6f1b77(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]User)
	bufferData := make([]User, len(dataIn), len(dataIn))
	mapData := make(map[time.Time]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]User) Department)
	result := make([]Department, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].CreateTime
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

func queryGroup_37e53ff8d9e8cce0d72071f5eacc22898cd03373(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]User)
	bufferData := make([]User, len(dataIn), len(dataIn))
	mapData := make(map[int]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]User) Department)
	result := make([]Department, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].UserId
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

func queryGroup_b08ca1fc16af4d4e7b65b696e8ba35882c0a4961(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]User)
	bufferData := make([]User, len(dataIn), len(dataIn))
	mapData := make(map[time.Time]int, len(dataIn))
	groupFunctorIn := groupFunctor.(func([]User) []Department)
	result := make([]Department, 0, len(dataIn))

	length := len(dataIn)
	nextData := make([]int, length, length)
	for i := 0; i != length; i++ {
		single := dataIn[i].CreateTime
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

func queryJoin_18a90660a498dc8a2c84eae90b4a430815a5d594(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]Admin)
	rightDataIn := rightData.([]User)
	joinFunctorIn := joinFunctor.(func(Admin, User) AdminUser)
	result := make([]AdminUser, 0, len(leftDataIn))

	emptyLeftData := Admin{}
	emptyRightData := User{}
	joinPlace = "inner"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[int]int, len(rightDataIn))
	mapDataFirst := make(map[int]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].UserId
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
		fieldValue := leftValue.AdminId
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

func queryJoin_1ba84e33b88ae2e0926f0db2690423b5df5992fc(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]Admin)
	rightDataIn := rightData.([]User)
	joinFunctorIn := joinFunctor.(func(Admin, User) AdminUser)
	result := make([]AdminUser, 0, len(leftDataIn))

	emptyLeftData := Admin{}
	emptyRightData := User{}
	joinPlace = "left"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[int]int, len(rightDataIn))
	mapDataFirst := make(map[int]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i].UserId
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
		fieldValue := leftValue.AdminId
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

func queryJoin_69248e7728f46e4a00d6cea63f2eabca6513fee8(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]User)
	rightDataIn := rightData.([]int)
	joinFunctorIn := joinFunctor.(func(User, int) User)
	result := make([]User, 0, len(leftDataIn))

	emptyLeftData := User{}
	emptyRightData := 0
	joinPlace = "right"

	nextData := make([]int, len(rightDataIn), len(rightDataIn))
	mapDataNext := make(map[int]int, len(rightDataIn))
	mapDataFirst := make(map[int]int, len(rightDataIn))

	for i := 0; i != len(rightDataIn); i++ {
		fieldValue := rightDataIn[i]
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
		fieldValue := leftValue.UserId
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

func querySelect_330d97a8f08ab419926dd507be00ec1c6a1de660(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]User)
	selectFunctorIn := selectFunctor.(func(User) Sex)
	result := make([]Sex, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySelect_7678e734f6e917b21b7a5d60d20beda1349d563e(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]int)
	selectFunctorIn := selectFunctor.(func(int) User)
	result := make([]User, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
}

func querySort_20065ed6b2e88afe33c0f58931e75db892366b73(data interface{}, sortType string) interface{} {
	dataIn := data.([]Admin)
	newData := make([]Admin, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].IsMale == false && newData[j].IsMale == true {
			return -1
		} else if newData[i].IsMale == true && newData[j].IsMale == false {
			return 1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_8e0b118cde44520b4231889be9e1bb2d83505d2f(data interface{}, sortType string) interface{} {
	dataIn := data.([]User)
	newData := make([]User, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].UserId < newData[j].UserId {
			return -1
		} else if newData[i].UserId > newData[j].UserId {
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

func querySort_c0c7516f15f736e69120d675686e3649b43feff4(data interface{}, sortType string) interface{} {
	dataIn := data.([]User)
	newData := make([]User, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(len(newData), func(i int, j int) int {
		if newData[i].UserId < newData[j].UserId {
			return 1
		} else if newData[i].UserId > newData[j].UserId {
			return -1
		}

		if newData[i].Name < newData[j].Name {
			return -1
		} else if newData[i].Name > newData[j].Name {
			return 1
		}

		if newData[i].CreateTime.Before(newData[j].CreateTime) {
			return -1
		} else if newData[i].CreateTime.After(newData[j].CreateTime) {
			return 1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func queryWhere_73d00d0e091a8cd964916be9d13848bedc08c8bb(data interface{}, whereFunctor interface{}) interface{} {
	dataIn := data.([]User)
	whereFunctorIn := whereFunctor.(func(User) bool)
	result := make([]User, 0, len(dataIn))

	for _, single := range dataIn {
		shouldStay := whereFunctorIn(single)
		if shouldStay == true {
			result = append(result, single)
		}
	}
	return result
}

func queryWhere_df6742e675632943121cefdb3ad29ba75c08eaac(data interface{}, whereFunctor interface{}) interface{} {
	dataIn := data.([]int)
	whereFunctorIn := whereFunctor.(func(int) bool)
	result := make([]int, 0, len(dataIn))

	for _, single := range dataIn {
		shouldStay := whereFunctorIn(single)
		if shouldStay == true {
			result = append(result, single)
		}
	}
	return result
}

func init() {

	language.QueryColumnMapMacroRegister([]User{}, "UserId", queryColumnMap_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0)

	language.QueryColumnMapMacroRegister([]int{}, ".", queryColumnMap_904b262f8e2329ec73c320ca0e5ca82f14165586)

	language.QueryColumnMapMacroRegister([]User{}, "[]UserId", queryColumnMap_99ac937255d05542023b7c6c1eba945d1b199327)

	language.QueryColumnMacroRegister([]User{}, "UserId", queryColumn_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0)

	language.QueryColumnMacroRegister([]User{}, ".", queryColumn_7ab1c12a87c03fba395eb892ababc4a58ad5dec9)

	language.QueryColumnMacroRegister([]int{}, ".", queryColumn_904b262f8e2329ec73c320ca0e5ca82f14165586)

	language.QueryColumnMacroRegister([]subtest.Address{}, "City", queryColumn_b60cd8d06e3e435a78322ac375157c99ea3ee15e)

	language.QueryCombineMacroRegister([]int{}, []User{}, (func(int, User) User)(nil), queryCombine_91ed07637405c43f5a89f7fc9cbd76e4e62acc92)

	language.QueryCombineMacroRegister([]Admin{}, []User{}, (func(Admin, User) AdminUser)(nil), queryCombine_d4934c4702a916f609947de20b2197cb0f02712e)

	language.QueryGroupMacroRegister([]int{}, ".", (func([]int) Department)(nil), queryGroup_0b68844c3917648685a142b599f750918372f47c)

	language.QueryGroupMacroRegister([]User{}, "CreateTime", (func([]User) Department)(nil), queryGroup_2a871c94537a18ef0ccfd971807ac9a23c6f1b77)

	language.QueryGroupMacroRegister([]User{}, "UserId", (func([]User) Department)(nil), queryGroup_37e53ff8d9e8cce0d72071f5eacc22898cd03373)

	language.QueryGroupMacroRegister([]User{}, "CreateTime", (func([]User) []Department)(nil), queryGroup_b08ca1fc16af4d4e7b65b696e8ba35882c0a4961)

	language.QueryJoinMacroRegister([]Admin{}, []User{}, "inner", "AdminId = UserId", (func(Admin, User) AdminUser)(nil), queryJoin_18a90660a498dc8a2c84eae90b4a430815a5d594)

	language.QueryJoinMacroRegister([]Admin{}, []User{}, "left", "AdminId = UserId", (func(Admin, User) AdminUser)(nil), queryJoin_1ba84e33b88ae2e0926f0db2690423b5df5992fc)

	language.QueryJoinMacroRegister([]User{}, []int{}, "right", "UserId = .", (func(User, int) User)(nil), queryJoin_69248e7728f46e4a00d6cea63f2eabca6513fee8)

	language.QuerySelectMacroRegister([]User{}, (func(User) Sex)(nil), querySelect_330d97a8f08ab419926dd507be00ec1c6a1de660)

	language.QuerySelectMacroRegister([]int{}, (func(int) User)(nil), querySelect_7678e734f6e917b21b7a5d60d20beda1349d563e)

	language.QuerySortMacroRegister([]Admin{}, "IsMale asc", querySort_20065ed6b2e88afe33c0f58931e75db892366b73)

	language.QuerySortMacroRegister([]User{}, "UserId asc", querySort_8e0b118cde44520b4231889be9e1bb2d83505d2f)

	language.QuerySortMacroRegister([]int{}, ". desc", querySort_af891d058d5a2e0a3ac4b4b291ae9bb959364795)

	language.QuerySortMacroRegister([]User{}, "UserId desc,Name asc,CreateTime asc", querySort_c0c7516f15f736e69120d675686e3649b43feff4)

	language.QueryWhereMacroRegister([]User{}, (func(User) bool)(nil), queryWhere_73d00d0e091a8cd964916be9d13848bedc08c8bb)

	language.QueryWhereMacroRegister([]int{}, (func(int) bool)(nil), queryWhere_df6742e675632943121cefdb3ad29ba75c08eaac)

}
