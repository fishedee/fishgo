package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func CompressGzip(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
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

func DecompressGzip(data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(dataReader)
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
