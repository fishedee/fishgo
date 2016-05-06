package modules

func InstallPackage(packageName string) error {
	_, err := runCmdSync("go", "install", packageName)
	return err
}
