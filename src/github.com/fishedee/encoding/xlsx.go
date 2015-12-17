package encoding

import (
	"github.com/tealeg/xlsx"
	"net/http"
	"io"
)

func EncodeXlsx(data [][]string,writer io.Writer)(error){
	file := xlsx.NewFile()
 	sheet, err := file.AddSheet("Sheet1")
	if err != nil{
		return err
	}
	for i := 0 ; i != len(data[0]) ; i++{
		sheet.Col(i).Width = 25
		sheet.Col(i).SetType(xlsx.CellTypeString)
	}
	for i := 0 ; i != len(data) ; i++{
		row := sheet.AddRow()
		for j := 0 ; j != len(data[i]) ; j++{
			cell := row.AddCell()
			cell.SetString( data[i][j] )
		}
	}
	file.Write(writer)
	return nil
}

func EncodeXlsxToRespnseWriter(data [][]string,writer http.ResponseWriter,title string)(error){
	writerHeader := writer.Header()
	writerHeader.Set("Content-Type","application/vnd.openxmlformats-officedocument; charset=UTF-8"); 
	writerHeader.Set("Pragma","public"); 
	writerHeader.Set("Expires","0"); 
	writerHeader.Set("Cache-Control","must-revalidate, post-check=0, pre-check=0"); 
	writerHeader.Set("Content-Type","application/force-download"); 
	writerHeader.Set("Content-Type","application/octet-stream"); 
	writerHeader.Set("Content-Type","application/download"); 
	writerHeader.Set("Content-Disposition","attachment;filename="+title+".xlsx"); 
	writerHeader.Set("Content-Transfer-Encoding","binary");
	return EncodeXlsx(data,writer)
}