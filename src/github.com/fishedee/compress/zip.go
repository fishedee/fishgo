package compress

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path"
)

func CompressZipFile(inputFileName string, outputFileName string) error {
	//建立读取文件
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	//建立写入文件
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	//zip压缩数据
	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()
	zipFileWriter, err := zipWriter.Create(path.Base(inputFileName))
	if err != nil {
		return err
	}
	_, err = io.Copy(zipFileWriter, inputFile)
	if err != nil {
		return err
	}
	return nil
}

func DecompressZipFile(inputFileName string, outputFileName string) error {
	//建立写入文件
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	//zip解压数据
	zipReader, err := zip.OpenReader(inputFileName)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	if len(zipReader.File) == 0 {
		return errors.New("zip is empty!")
	}
	fileOpen, err := zipReader.File[0].Open()
	if err != nil {
		return err
	}
	defer fileOpen.Close()
	_, err = io.Copy(outputFile, fileOpen)
	if err != nil {
		return err
	}
	return nil
}
