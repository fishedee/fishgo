package modules

import (
	"io/ioutil"
)

func IsExistDir(dirName string) bool {
	_, err := ioutil.ReadDir(dirName)
	return err == nil
}