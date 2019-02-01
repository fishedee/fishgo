package testdata

import (
	"github.com/fishedee/language"
	"github.com/fishedee/language/querygen/testdata/subtest"
)

func queryColumnMap_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0(data interface{}, column string) interface{} {
	dataIn := data.([]User)
	result := make(map[int]User, len(dataIn))

	for _, single := range dataIn {
		result[single.UserId] = single
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

func queryColumn_b60cd8d06e3e435a78322ac375157c99ea3ee15e(data interface{}, column string) interface{} {
	dataIn := data.([]subtest.Address)
	result := make([]string, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.City
	}
	return result
}

func queryGroup_37e53ff8d9e8cce0d72071f5eacc22898cd03373(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	dataIn := data.([]User)
	groupFunctorIn := groupFunctor.(func([]User) Department)
	newData := make([]User, len(dataIn), len(dataIn))
	copy(newData, dataIn)
	newData2 := make([]Department, 0, len(dataIn))

	language.QueryGroupInternal(len(newData), func(i int, j int) int {
		if newData[i].UserId < newData[j].UserId {
			return -1
		} else if newData[i].UserId > newData[j].UserId {
			return 1
		}

		return 0
	}, func(i int, j int) {
		newData[j], newData[i] = newData[i], newData[j]
	}, func(i int, j int) {
		single := groupFunctorIn(newData[i:j])
		newData2 = append(newData2, single)
	})
	return newData2
}

func queryJoin_18a90660a498dc8a2c84eae90b4a430815a5d594(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]Admin)
	rightDataIn := rightData.([]User)
	joinFunctorIn := joinFunctor.(func(Admin, User) AdminUser)
	newRightData := make([]User, len(rightDataIn), len(rightDataIn))
	copy(newRightData, rightDataIn)
	newData2 := make([]AdminUser, 0, len(leftDataIn))

	emptyLeftData := Admin{}
	emptyRightData := User{}
	language.QueryJoinInternal(
		"inner",
		len(leftDataIn),
		len(rightDataIn),
		func(i int, j int) int {
			if newRightData[i].UserId < newRightData[j].UserId {
				return -1
			} else if newRightData[i].UserId > newRightData[j].UserId {
				return 1
			}

			return 0
		},
		func(i int, j int) {
			newRightData[j], newRightData[i] = newRightData[i], newRightData[j]
		},
		func(i int, j int) int {
			if leftDataIn[i].AdminId < newRightData[j].UserId {
				return -1
			} else if leftDataIn[i].AdminId > newRightData[j].UserId {
				return 1
			}

			return 0
		},
		func(i int, j int) {
			left := emptyLeftData
			if i != -1 {
				left = leftDataIn[i]
			}
			right := emptyRightData
			if j != -1 {
				right = newRightData[j]
			}
			single := joinFunctorIn(left, right)
			newData2 = append(newData2, single)
		},
	)
	return newData2
}

func queryJoin_1ba84e33b88ae2e0926f0db2690423b5df5992fc(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFunctor interface{}) interface{} {
	leftDataIn := leftData.([]Admin)
	rightDataIn := rightData.([]User)
	joinFunctorIn := joinFunctor.(func(Admin, User) AdminUser)
	newRightData := make([]User, len(rightDataIn), len(rightDataIn))
	copy(newRightData, rightDataIn)
	newData2 := make([]AdminUser, 0, len(leftDataIn))

	emptyLeftData := Admin{}
	emptyRightData := User{}
	language.QueryJoinInternal(
		"left",
		len(leftDataIn),
		len(rightDataIn),
		func(i int, j int) int {
			if newRightData[i].UserId < newRightData[j].UserId {
				return -1
			} else if newRightData[i].UserId > newRightData[j].UserId {
				return 1
			}

			return 0
		},
		func(i int, j int) {
			newRightData[j], newRightData[i] = newRightData[i], newRightData[j]
		},
		func(i int, j int) int {
			if leftDataIn[i].AdminId < newRightData[j].UserId {
				return -1
			} else if leftDataIn[i].AdminId > newRightData[j].UserId {
				return 1
			}

			return 0
		},
		func(i int, j int) {
			left := emptyLeftData
			if i != -1 {
				left = leftDataIn[i]
			}
			right := emptyRightData
			if j != -1 {
				right = newRightData[j]
			}
			single := joinFunctorIn(left, right)
			newData2 = append(newData2, single)
		},
	)
	return newData2
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

func init() {

	language.QueryColumnMapMacroRegister([]User{}, "UserId", queryColumnMap_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0)

	language.QueryColumnMacroRegister([]User{}, "UserId", queryColumn_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0)

	language.QueryColumnMacroRegister([]subtest.Address{}, "City", queryColumn_b60cd8d06e3e435a78322ac375157c99ea3ee15e)

	language.QueryGroupMacroRegister([]User{}, "UserId", (func([]User) Department)(nil), queryGroup_37e53ff8d9e8cce0d72071f5eacc22898cd03373)

	language.QueryJoinMacroRegister([]Admin{}, []User{}, "inner", "AdminId = UserId", (func(Admin, User) AdminUser)(nil), queryJoin_18a90660a498dc8a2c84eae90b4a430815a5d594)

	language.QueryJoinMacroRegister([]Admin{}, []User{}, "left", "AdminId = UserId", (func(Admin, User) AdminUser)(nil), queryJoin_1ba84e33b88ae2e0926f0db2690423b5df5992fc)

	language.QuerySelectMacroRegister([]User{}, (func(User) Sex)(nil), querySelect_330d97a8f08ab419926dd507be00ec1c6a1de660)

	language.QuerySortMacroRegister([]User{}, "UserId asc", querySort_8e0b118cde44520b4231889be9e1bb2d83505d2f)

	language.QuerySortMacroRegister([]User{}, "UserId desc,Name asc,CreateTime asc", querySort_c0c7516f15f736e69120d675686e3649b43feff4)

	language.QueryWhereMacroRegister([]User{}, (func(User) bool)(nil), queryWhere_73d00d0e091a8cd964916be9d13848bedc08c8bb)

}
