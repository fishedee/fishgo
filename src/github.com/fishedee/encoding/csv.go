package encoding

import (
	"bytes"
	"encoding/csv"
)

func EncodeCsv(data [][]string) ([]byte, error) {
	var writer bytes.Buffer
	csvWriter := csv.NewWriter(&writer)
	writer.Write([]byte("\xEF\xBB\xBF"))
	for i := 0; i != len(data); i++ {
		err := csvWriter.Write(data[i])
		if err != nil {
			return nil, err
		}
	}
	csvWriter.Flush()
	return writer.Bytes(), nil
}
