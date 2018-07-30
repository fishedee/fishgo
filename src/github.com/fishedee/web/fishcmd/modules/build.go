package modules

func BuildPackage(appName string) error {
	_, err := runCmdSync("go", "build", "-o", appName)
	return err
}
