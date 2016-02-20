package main

import (
	"strings"
)

func GetGenerateFileName(name string) string {
	nameByte := []byte(name)
	nameByte = nameByte[0 : len(nameByte)-3]
	return string(nameByte) + "_ext.go"
}

func IsGenerateFileName(name string) bool {
	return strings.HasSuffix(name, "_ext.go")
}

func GetOriginFileName(name string) string {
	nameByte := []byte(name)
	nameByte = nameByte[0 : len(nameByte)-7]
	return string(nameByte) + ".go"
}
