package main

import (
	"errors"
	"os"
)

var Config struct {
	fileregex string
}

func ReadConfig() error {
	argv := os.Args
	argv = argv[1:]
	if len(argv) < 1 {
		return errors.New("need a file regex argument")
	}
	Config.fileregex = argv[0]
	return nil
}
