package a

import (
	"os"
)

var workingDir string

func GetWorkingDir() string {
	return workingDir
}

func init() {
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}
