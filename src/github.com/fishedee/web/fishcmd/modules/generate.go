package modules

import (
	"fmt"
)

func GeneratePackage(packageName string) error {
	cmdLog, err := runCmdSync("go", "generate", packageName)
	fmt.Printf("%v", string(cmdLog))
	return err
}
