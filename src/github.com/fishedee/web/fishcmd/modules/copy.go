package modules

func CopyFile(source string, target string) error {
	_, err := runCmdSync("cp", "-rf", source, target)
	return err
}
