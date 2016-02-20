package main

import (
	"fmt"
)

func Generator(data map[string][]ParserInfo) error {
	fmt.Println(fmt.Sprintf("%+v", data))
	return nil
}
