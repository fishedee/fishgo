package main

import (
	"errors"
	"go/format"
	"io/ioutil"
	"os"
)

func generateSingleFileFormat(filename string, data string) (string, error) {
	result, err := format.Source([]byte(data))
	if err != nil {
		return "", errors.New(err.Error() + "," + data)
	}
	return string(result), nil
}

func generateSingleFileWrite(filename string, data string) error {
	oldData, err := ioutil.ReadFile(filename)
	if err == nil && string(oldData) == data {
		return nil
	}
	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func generateSingleFileTest(dirname string, data []ParserInfo) error {
	filename := dirname + "/" + GetGenerateTestFileName(dirname)
	result, err := generateSingleTestFileContent(data)
	if err != nil {
		return err
	}
	if result == "" {
		return nil
	}

	result, err = generateSingleFileFormat(filename, result)
	if err != nil {
		return err
	}

	err = generateSingleFileWrite(filename, result)
	if err != nil {
		return err
	}

	return nil
}

func generateSingleFileNormal(dirname string, data []ParserInfo) error {
	filename := dirname + "/" + GetGenerateFileName(dirname)

	result, err := generateSingleFileContent(data)
	if err != nil {
		return err
	}

	result, err = generateSingleFileFormat(filename, result)
	if err != nil {
		return err
	}

	err = generateSingleFileWrite(filename, result)
	if err != nil {
		return err
	}

	return nil
}

func generateSingleFile(dirname string, data []ParserInfo) error {
	err := generateSingleFileNormal(dirname, data)
	if err != nil {
		return err
	}

	err = generateSingleFileTest(dirname, data)
	if err != nil {
		return err
	}

	return nil
}

func Generator(data map[string][]os.FileInfo) error {
	for singleKey, singleDir := range data {
		singleResult := []ParserInfo{}
		for _, singleFile := range singleDir {
			singleFileResult, err := ParserSingleFile(singleKey + "/" + singleFile.Name())
			if err != nil {
				return err
			}
			singleResult = append(singleResult, singleFileResult)
		}
		err := generateSingleFile(singleKey, singleResult)
		if err != nil {
			return errors.New(singleKey + ":" + err.Error())
		}
	}
	return nil
}

func init() {
	dirDeclTypeCache = map[string]map[string]bool{}
}
