package main

import (
	"path"
)

func GetGenerateFileName(name string) string {
	basepath := path.Base(name)
	return basepath + "_fishgen.go"
}
