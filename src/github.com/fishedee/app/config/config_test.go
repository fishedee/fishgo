package config

import (
	. "github.com/fishedee/assert"
	"os"
	"testing"
	"time"
)

type testNormal struct {
	A int
	B string
	C bool
	D time.Duration
	E []string
	F float64
}

func TestConfigNormal(t *testing.T) {
	config, _ := NewConfig("ini", "testdata/normal.ini")
	result := testNormal{}
	config.GetStruct("", &result)
	AssertEqual(t, result, testNormal{
		A: 2,
		B: "asdf",
		C: true,
		D: time.Second * 10,
		E: []string{"123", "323", "sdf"},
		F: 2.4,
	})
}

func TestConfigProd(t *testing.T) {
	os.Setenv("RUNMODE", "prod")
	config, _ := NewConfig("ini", "testdata/prod.ini")
	AssertEqual(t, config.GetInt("a"), 8)
}
