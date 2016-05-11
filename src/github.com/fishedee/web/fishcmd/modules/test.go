package modules

func TestPackage(packageName string) error {
	_, err := runCmdSyncAndStdOutput("go", "test", "-v", "-p", "1", packageName)
	return err
}
