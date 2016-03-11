package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strings"
)

func randFileName() (string, error) {
	randFileName := make([]byte, 20)
	_, err := rand.Read(randFileName)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randFileName), nil
}

func createTempDir(dir string) (string, error) {
	var tempDir string
	if dir != "" {
		tempDir = strings.TrimRight(os.TempDir(), "/") + "/" + dir
		tempDirInfo, err := os.Lstat(tempDir)
		if tempDirInfo != nil && err == nil {
			if tempDirInfo.IsDir() == false {
				return "", errors.New("has same name in temp dir and not dir " + dir)
			}
		} else {
			err := os.Mkdir(tempDir, os.ModePerm)
			if err != nil {
				return "", err
			}
		}
		tempDir += "/"
	} else {
		tempDir = strings.TrimRight(os.TempDir(), "/") + "/"
	}
	return tempDir, nil
}

func CreateTempFile(dir string, suffix string) (string, error) {
	var err error
	dir = strings.Trim(dir, " ")
	dir = strings.Trim(dir, "/")
	suffix = strings.Trim(suffix, " ")
	//创建文件夹
	dir, err = createTempDir(dir)
	if err != nil {
		return "", err
	}
	//创建文件
	fileName, err := randFileName()
	if err != nil {
		return "", err
	}
	return dir + fileName + suffix, nil
}
