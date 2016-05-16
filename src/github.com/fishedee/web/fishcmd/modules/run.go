package modules

import (
	"os/exec"
)

var (
	cmd *exec.Cmd
)

func RunPackage(packageName string, isAsync bool) error {
	var err error
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
	}
	runCmdSync("pkill", "-9", packageName)
	if isAsync {
		cmd, err = runCmdAsync("./" + packageName)
		return err
	} else {
		cmd, err = runCmdSyncAndStdOutput("./" + packageName)
		return err
	}
}
