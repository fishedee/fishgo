package render

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type HtmlFormatter struct {
	tmpl *template.Template
}

func (this *HtmlFormatter) load(dir string) error {
	fileList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() == true {
			return nil
		}
		if strings.HasSuffix(f.Name(), ".html") == false {
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		return err
	}
	this.tmpl, err = template.ParseFiles(fileList...)
	if err != nil {
		return err
	}
	return nil
}

func (this *HtmlFormatter) Name() string {
	return "html"
}

func (this *HtmlFormatter) Format(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if dataArray, isOk := data.([]interface{}); isOk == true {
		fileName := dataArray[0].(string)
		fileData := dataArray[1]
		err := this.tmpl.ExecuteTemplate(w, fileName, fileData)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid data type for template formatter")
	}
	return nil
}

func NewHtmlFormatter(dir string) (*HtmlFormatter, error) {
	htmlFormatter := &HtmlFormatter{}
	err := htmlFormatter.load(dir)
	if err != nil {
		return nil, err
	}
	return htmlFormatter, nil
}
