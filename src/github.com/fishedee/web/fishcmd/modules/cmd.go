package modules

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func runCmd(isSync bool, stdOutput bool, name string, args ...string) (*exec.Cmd, []byte, error) {
	var buf = bytes.NewBuffer([]byte(""))
	cmd := exec.Command(name)
	if stdOutput == false {
		cmd.Stdout = buf
		cmd.Stderr = buf
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
	}
	cmd.Args = append([]string{name}, args...)
	if name == "go" {
		cmd.Env = append(os.Environ(), "GOGC=off")
	} else {
		cmd.Env = os.Environ()
	}
	if isSync {
		err := cmd.Run()
		return cmd, buf.Bytes(), err
	} else {
		go cmd.Run()
		return cmd, nil, nil
	}
}

func runCmdSync(name string, args ...string) ([]byte, error) {
	_, data, err := runCmd(true, false, name, args...)
	if err != nil {
		return nil, errors.New(string(data))
	}
	return data, nil
}

func runCmdSyncAndStdOutput(name string, args ...string) (*exec.Cmd, error) {
	cmd, _, err := runCmd(true, true, name, args...)
	return cmd, err
}

func runCmdAsync(name string, args ...string) (*exec.Cmd, error) {
	cmd, _, err := runCmd(false, true, name, args...)
	return cmd, err
}
