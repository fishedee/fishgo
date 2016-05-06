package modules

import (
	"io/ioutil"
	"os"
	"time"
)

func IsExistDir(dirName string) bool {
	_, err := ioutil.ReadDir(dirName)
	return err == nil
}

func GetFileModifyTime(fileName string) (time.Time, error) {
	fileInfo, err := os.Lstat(fileName)
	if err != nil {
		return time.Time{}, nil
	}
	return fileInfo.ModTime(), nil
}
