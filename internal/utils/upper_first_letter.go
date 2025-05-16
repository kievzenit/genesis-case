package utils

import "strings"

func UpperFirstLetter(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}