package modules

func TestPackage(packageName string, args string) error {
	_, err := runCmdSyncAndStdOutput("go", "test", "-v", "-p", "1", "-args", args, packageName)
	return err
}
