package slice

import "strings"

func ParseFromString(str string) []string {
	str = strings.Trim(str, "{}")
	return strings.Split(str, ",")
}
