package main

import (
	"fmt"
	"os"
)

func main() {
	err := ReadConfig()
	if err != nil {
		fmt.Println("read config error " + err.Error())
		os.Exit(1)
		return
	}

	data, err := ReadDir(".")
	if err != nil {
		fmt.Println("read dir error " + err.Error())
		os.Exit(1)
		return
	}

	data, err = FilterDir(data)
	if err != nil {
		fmt.Println("filter dir error " + err.Error())
		os.Exit(1)
		return
	}

	err = Generator(data)
	if err != nil {
		fmt.Println("generate dir error " + err.Error())
		os.Exit(1)
		return
	}
}
