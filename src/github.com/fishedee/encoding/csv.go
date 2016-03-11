package encoding

import (
	"bytes"
	"encoding/csv"
)

func EncodeCsv(data [][]string) ([]byte, error) {
	var writer bytes.Buffer
	data1, err := EncodeCsvHeader()
	if err != nil {
		return nil, err
	}
	data2, err := EncodeCsvBody(data)
	if err != nil {
		return nil, err
	}
	writer.Write(data1)
	writer.Write(data2)
	return writer.Bytes(), nil
}

func EncodeCsvHeader() ([]byte, error) {
	return []byte("\xEF\xBB\xBF"), nil
}

func EncodeCsvBody(data [][]string) ([]byte, error) {
	var writer bytes.Buffer
	csvWriter := csv.NewWriter(&writer)
	for i := 0; i != len(data); i++ {
		err := csvWriter.Write(data[i])
		if err != nil {
			return nil, err
		}
	}
	csvWriter.Flush()
	return writer.Bytes(), nil
}
