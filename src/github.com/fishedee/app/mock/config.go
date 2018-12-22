package main

import (
	"errors"
	"os"
)

var Config struct {
	fileregex string
	typeregex string
}

func ReadConfig() error {
	argv := os.Args
	argv = argv[1:]
	if len(argv) < 2 {
		return errors.New("need a file regex argument and a type name regex argument")
	}
	Config.fileregex = argv[0]
	Config.typeregex = argv[1]
	return nil
}
