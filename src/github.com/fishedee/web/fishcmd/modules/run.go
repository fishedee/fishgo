package modules

import (
	"os"
)

var (
	process *os.Process
)

func RunPackage(packageName string) error {
	var err error
	if process != nil {
		process.Kill()
	}
	process, err = runCmdAsync("./" + packageName)
	return err
}
