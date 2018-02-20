package cors

import (
	"github.com/rs/cors"
	"net/http"
)

type Cors interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type corsImplement struct {
	cors *cors.Cors
}

type CorsConfig struct {
	AllowedOrigins     []string `config:"allowedorigins"`
	AllowedMethods     []string `config:"allowedmethods"`
	AllowedHeaders     []string `config:"allowedheaders"`
	ExposedHeaders     []string `config:"exposedheaders"`
	MaxAge             int      `config:"maxage"`
	AllowCredentials   bool     `config:"allowcredentials"`
	OptionsPassthrough bool     `config:"optionspassthrough"`
	Debug              bool     `config:"debug"`
}

func NewCors(config CorsConfig) (Cors, error) {
	option := cors.Options{
		AllowedOrigins:     config.AllowedOrigins,
		AllowedMethods:     config.AllowedMethods,
		AllowedHeaders:     config.AllowedHeaders,
		ExposedHeaders:     config.ExposedHeaders,
		MaxAge:             config.MaxAge,
		AllowCredentials:   config.AllowCredentials,
		OptionsPassthrough: config.OptionsPassthrough,
		Debug:              config.Debug,
	}
	cors := cors.New(option)
	return &corsImplement{
		cors: cors,
	}, nil
}

func (this *corsImplement) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	this.cors.ServeHTTP(w, r, next)
}
