package main

import (
	"errors"
	"os"
)

var Config struct {
	fileregex string
	force     bool
}

func ReadConfig() error {
	argv := os.Args
	argv = argv[1:]
	for _, singleArgv := range argv {
		if singleArgv[0] == '-' {
			if singleArgv == "-force" {
				Config.force = true
			}
		} else {
			Config.fileregex = singleArgv
		}
	}
	if Config.fileregex == "" {
		return errors.New("need a file regex argument")
	}
	return nil
}
