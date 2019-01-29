package testdata

import (
	"github.com/fishedee/language"
)

func queryColumn_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0(data interface{}, column string) interface{} {
	dataIn := data.([]User)
	result := make([]int, len(dataIn), len(dataIn))

	for i, single := range dataIn {
		result[i] = single.UserId
	}
	return result
}

func init() {

	language.QueryColumnMacroRegister([]User{}, "UserId", queryColumn_4cb77d7ba8d1eeb02c714a053eefbaa439c736f0)

}
