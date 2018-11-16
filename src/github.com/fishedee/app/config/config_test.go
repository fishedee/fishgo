package config

import (
	"fmt"
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
	config.MustBind("", &result)
	AssertEqual(t, result, testNormal{
		A: 2,
		B: "asdf",
		C: true,
		D: time.Second * 10,
		E: []string{"123", "323", "sdf"},
		F: 2.4,
	})
}

func TestConfigDefaultDev(t *testing.T) {
	config, _ := NewConfig("ini", "testdata/prod.ini")
	AssertEqual(t, config.MustInt("a"), 7)
}

func TestConfigEnvProd(t *testing.T) {
	os.Setenv("RUNMODE", "prod")
	config, _ := NewConfig("ini", "testdata/prod.ini")
	AssertEqual(t, config.MustInt("a"), 8)
}

func TestConfigConfigTest(t *testing.T) {
	os.Setenv("RUNMODE", "")
	config, _ := NewConfig("ini", "testdata/test.ini")
	fmt.Println(config.String("runmode"))
	AssertEqual(t, config.MustInt("a"), 9)
}

func TestConfigError(t *testing.T) {
	config, _ := NewConfig("ini", "testdata/error.ini")

	var err error
	//格式错误
	_, err = config.Int("a")
	AssertEqual(t, err != nil, true)
	_, err = config.Bool("c")
	AssertEqual(t, err != nil, true)

	//空值略过
	_, err = config.String("b")
	AssertEqual(t, err == nil, true)
	_, err = config.Duration("d")
	AssertEqual(t, err == nil, true)
	_, err = config.StringList("e")
	AssertEqual(t, err == nil, true)
	_, err = config.Float("f")
	AssertEqual(t, err == nil, true)

	data := testNormal{}
	err = config.Bind("", &data)
	t.Logf("%v", err)
	AssertEqual(t, err != nil, true)
}
