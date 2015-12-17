package encoding

import (
	"encoding/csv"
	"net/http"
	"io"
)

func EncodeCsv(data [][]string,writer io.Writer)(error){
	csvWriter := csv.NewWriter(writer)
	for i := 0 ; i != len(data) ; i++{
		err := csvWriter.Write(data[i])
		if err != nil{
			return err
		}
	}
	csvWriter.Flush()
	return nil
}

func EncodeCsvToRespnseWriter(data [][]string,writer http.ResponseWriter,title string)(error){
	writerHeader := writer.Header()
	writerHeader.Set("Content-Type","application/vnd.ms-excel; charset=UTF-8"); 
	writerHeader.Set("Pragma","public"); 
	writerHeader.Set("Expires","0"); 
	writerHeader.Set("Cache-Control","must-revalidate, post-check=0, pre-check=0"); 
	writerHeader.Set("Content-Type","application/force-download"); 
	writerHeader.Set("Content-Type","application/octet-stream"); 
	writerHeader.Set("Content-Type","application/download"); 
	writerHeader.Set("Content-Disposition","attachment;filename="+title+".csv"); 
	writerHeader.Set("Content-Transfer-Encoding","binary");
	writer.Write([]byte("\xEF\xBB\xBF"))
	return EncodeCsv(data,writer)
}