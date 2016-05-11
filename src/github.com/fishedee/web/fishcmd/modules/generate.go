package modules

func GeneratePackage(packageName string) error {
	_, err := runCmdSync("go", "generate", packageName)
	return err
}
