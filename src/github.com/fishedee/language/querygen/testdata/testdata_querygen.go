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

func init() {

	language.QueryColumnMacroRegister([]User{}, "UserId", queryColumn_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0)

	language.QueryColumnMacroRegister([]subtest.Address{}, "City", queryColumn_b60cd8d06e3e435a78322ac375157c99ea3ee15e)

	language.QuerySelectMacroRegister([]User{}, (func(User) Sex)(nil), querySelect_330d97a8f08ab419926dd507be00ec1c6a1de660)

}
