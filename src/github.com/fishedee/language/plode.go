package language

import (
	"strconv"
	"strings"
)

func Explode(input string, separator string) []string {
	dataResult := strings.Split(input, separator)
	dataResultNew := []string{}
	for _, singleResult := range dataResult {
		singleResult = strings.Trim(singleResult, " ")
		if len(singleResult) == 0 {
			continue
		}
		dataResultNew = append(dataResultNew, singleResult)
	}
	return dataResultNew
}

func Implode(data []string, separator string) string {
	return strings.Join(data, separator)
}

func ExplodeInt(input string, separator string) []int {
	dataResult := strings.Split(input, separator)
	dataResultNew := []int{}
	for _, singleResult := range dataResult {
		singleResult = strings.Trim(singleResult, " ")
		if len(singleResult) == 0 {
			continue
		}
		singleResultInt, err := strconv.Atoi(singleResult)
		if err != nil {
			panic(err)
		}
		dataResultNew = append(dataResultNew, singleResultInt)
	}
	return dataResultNew
}

func ImplodeInt(data []int, separator string) string {
	result := []string{}
	for _, singleData := range data {
		result = append(result, strconv.Itoa(singleData))
	}
	return strings.Join(result, separator)
}
