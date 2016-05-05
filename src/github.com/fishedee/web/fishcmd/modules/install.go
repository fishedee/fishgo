package modules

func InstallPackage(packageName string) error {
	_, err := runCmdSync("go", "install", "-buildmode=shared", "-linkshared", packageName)
	return err
}
