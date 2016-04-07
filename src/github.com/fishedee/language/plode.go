package language

import (
	"strings"
)

func Explode(input string, separator string, data interface{}) {
	dataResult := strings.Split(input, separator)
	dataResultNew := []string{}
	for _, singleResult := range dataResult {
		singleResult = strings.Trim(singleResult, " ")
		if len(singleResult) == 0 {
			continue
		}
		dataResultNew = append(dataResultNew, singleResult)
	}
	err := MapToArray(dataResultNew, data, "json")
	if err != nil {
		panic(err)
	}
}

func Implode(data interface{}, separator string) string {
	dataResult := []string{}
	err := MapToArray(data, &dataResult, "json")
	if err != nil {
		panic(err)
	}
	return strings.Join(dataResult, separator)
}
