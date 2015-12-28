package encoding

import (
	"bytes"
	"github.com/tealeg/xlsx"
)

func EncodeXlsx(data [][]string) ([]byte, error) {
	var writer bytes.Buffer
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}
	for i := 0; i != len(data[0]); i++ {
		sheet.Col(i).Width = 25
		sheet.Col(i).SetType(xlsx.CellTypeString)
	}
	for i := 0; i != len(data); i++ {
		row := sheet.AddRow()
		for j := 0; j != len(data[i]); j++ {
			cell := row.AddCell()
			cell.SetString(data[i][j])
		}
	}
	file.Write(&writer)
	return writer.Bytes(), nil
}
