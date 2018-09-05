package render

import (
	"errors"
	"net/http"
	"strings"
)

type RenderFormatter interface {
	Name() string
	Format(w http.ResponseWriter, r *http.Request, data interface{}) error
}

type Render interface {
	Header(key string, value string)
	Status(code int)
	Format(name string, data interface{}) error
	ResponseWriter() http.ResponseWriter
	Request() *http.Request
}

type RenderFactory interface {
	RegisterFormatter(formatter RenderFormatter)
	Create(w http.ResponseWriter, r *http.Request) Render
}

type RenderConfig struct {
	TemplateDir string `config:"templatedir"`
}

type renderFactoryImplement struct {
	formatter map[string]RenderFormatter
}

func NewRenderFactory(config RenderConfig) (RenderFactory, error) {
	impl := &renderFactoryImplement{
		formatter: map[string]RenderFormatter{},
	}

	preFormatter := []func() (RenderFormatter, error){
		func() (RenderFormatter, error) {
			return NewRawFormatter()
		},
		func() (RenderFormatter, error) {
			return NewTextFormatter()
		},
		func() (RenderFormatter, error) {
			return NewRedirectFormatter()
		},
		func() (RenderFormatter, error) {
			return NewJsonFormatter()
		},
		func() (RenderFormatter, error) {
			return NewFileFormatter()
		},
	}
	if config.TemplateDir != "" {
		preFormatter = append(preFormatter, func() (RenderFormatter, error) {
			return NewHtmlFormatter(config.TemplateDir)
		})
	}

	for _, singlePreFormatter := range preFormatter {
		formatter, err := singlePreFormatter()
		if err != nil {
			return nil, err
		}
		impl.RegisterFormatter(formatter)
	}
	return impl, nil
}

func (this *renderFactoryImplement) RegisterFormatter(formatter RenderFormatter) {
	this.formatter[formatter.Name()] = formatter
}

func (this *renderFactoryImplement) Create(w http.ResponseWriter, r *http.Request) Render {
	return newRender(w, r, this.formatter)
}

type renderImplement struct {
	formatter map[string]RenderFormatter
	w         http.ResponseWriter
	r         *http.Request
}

func newRender(w http.ResponseWriter, r *http.Request, formatter map[string]RenderFormatter) Render {
	render := &renderImplement{}
	render.w = w
	render.r = r
	render.formatter = formatter
	return render
}

func (this *renderImplement) Header(key string, value string) {
	this.w.Header().Set(key, value)
}

func (this *renderImplement) Status(code int) {
	this.w.WriteHeader(code)
}

func (this *renderImplement) Format(name string, data interface{}) error {
	formatter, isExist := this.formatter[strings.ToLower(name)]
	if isExist == false {
		return errors.New("dos not exist formatter " + name)
	}
	return formatter.Format(this.w, this.r, data)
}

func (this *renderImplement) ResponseWriter() http.ResponseWriter {
	return this.w
}

func (this *renderImplement) Request() *http.Request {
	return this.r
}
