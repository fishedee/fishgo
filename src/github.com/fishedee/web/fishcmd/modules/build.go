package modules

func BuildPackage(packageName string, appName string) error {
	_, err := runCmdSync("go", "build", "-linkshared", "-o", appName, packageName)
	return err
}
