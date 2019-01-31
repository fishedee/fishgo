package testdata

import (
	"github.com/fishedee/language"
	"github.com/fishedee/language/querygen/testdata/subtest"
)

func queryColumn_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0(data interface{}, column string) interface{} {
	dataIn := data.([]User)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.UserId
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

func querySelect_330d97a8f08ab419926dd507be00ec1c6a1de660(data interface{}, selectFunctor interface{}) interface{} {
	dataIn := data.([]User)
	selectFunctorIn := selectFunctor.(func(User) Sex)
	result := make([]Sex, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = selectFunctorIn(single)
	}
	return result
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

func querySort_8e0b118cde44520b4231889be9e1bb2d83505d2f(data interface{}, sortType string) interface{} {
	dataIn := data.([]User)
	newData := make([]User, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(func() int {
		return len(newData)
	}, func(i int, j int) bool {
		if newData[i].UserId < newData[j].UserId {
			return true
		} else if newData[i].UserId > newData[j].UserId {
			return false
		}

		return false
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func querySort_c0c7516f15f736e69120d675686e3649b43feff4(data interface{}, sortType string) interface{} {
	dataIn := data.([]User)
	newData := make([]User, len(dataIn), len(dataIn))
	copy(newData, dataIn)

	language.QuerySortInternal(func() int {
		return len(newData)
	}, func(i int, j int) bool {
		if newData[i].UserId < newData[j].UserId {
			return false
		} else if newData[i].UserId > newData[j].UserId {
			return true
		}

		if newData[i].Name < newData[j].Name {
			return true
		} else if newData[i].Name > newData[j].Name {
			return false
		}

		if newData[i].CreateTime.Before(newData[j].CreateTime) {
			return true
		} else if newData[i].CreateTime.After(newData[j].CreateTime) {
			return false
		}

		return false
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	})
	return newData
}

func init() {

	language.QueryColumnMacroRegister([]User{}, "UserId", queryColumn_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0)

	language.QueryColumnMacroRegister([]subtest.Address{}, "City", queryColumn_b60cd8d06e3e435a78322ac375157c99ea3ee15e)

	language.QuerySelectMacroRegister([]User{}, (func(User) Sex)(nil), querySelect_330d97a8f08ab419926dd507be00ec1c6a1de660)

	language.QueryWhereMacroRegister([]User{}, (func(User) bool)(nil), queryWhere_73d00d0e091a8cd964916be9d13848bedc08c8bb)

	language.QuerySortMacroRegister([]User{}, "UserId asc", querySort_8e0b118cde44520b4231889be9e1bb2d83505d2f)

	language.QuerySortMacroRegister([]User{}, "UserId desc,Name asc,CreateTime asc", querySort_c0c7516f15f736e69120d675686e3649b43feff4)

}
