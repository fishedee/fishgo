package render

import (
	"errors"
	"net/http"
	"os"
	"path"
)

type FileFormatter struct {
}

func (this *FileFormatter) Name() string {
	return "file"
}

func (this *FileFormatter) Format(w http.ResponseWriter, r *http.Request, data interface{}) error {
	var result string
	if fileName, isOk := data.(string); isOk == true {
		result = fileName
	} else {
		return errors.New("invalid data type for file formatter")
	}
	file, err := os.Open(result)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	http.ServeContent(w, r, path.Ext(result), fileInfo.ModTime(), file)
	return nil
}

func NewFileFormatter() (*FileFormatter, error) {
	return &FileFormatter{}, nil
}
