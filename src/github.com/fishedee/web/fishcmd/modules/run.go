package modules

import (
	"os/exec"
)

var (
	cmd *exec.Cmd
)

func RunPackage(packageName string) error {
	var err error
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
	}
	cmd, err = runCmdAsync("./" + packageName)
	return err
}
