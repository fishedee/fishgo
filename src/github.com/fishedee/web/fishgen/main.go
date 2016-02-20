package main

import (
	"fmt"
)

func main() {
	err := ReadConfig()
	if err != nil {
		fmt.Println("read config error " + err.Error())
		return
	}

	data, err := ReadDir(".")
	if err != nil {
		fmt.Println("read dir error " + err.Error())
		return
	}

	data, err = FilterDir(data)
	if err != nil {
		fmt.Println("filter dir error " + err.Error())
		return
	}
	fmt.Println(data)

	parserData, err := Parser(data)
	if err != nil {
		fmt.Println("parser dir error " + err.Error())
		return
	}

	err = Generator(parserData)
	if err != nil {
		fmt.Println("generate dir error " + err.Error())
		return
	}
}
