package subtest

import (
	. "github.com/fishedee/language"
)

type Address struct {
	AddressId int
	City      string
}

func logic() {
	QueryColumn([]Address{}, "City")
}
