package compress

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

func CompressZlib(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DecompressZlib(data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)
	gzipReader, err := zlib.NewReader(dataReader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	result, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}
	return result, nil
}
