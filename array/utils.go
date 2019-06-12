package array

import "strings"

func parseArray(array string) []string {
	results := make([]string, 0)
	array = strings.Replace(array, "{", "", -1)
	array = strings.Replace(array, "}", "", -1)
	if array == "" {
		return results
	}
	matches := strings.Split(array, ",")
	for _, match := range matches {
		match := strings.Trim(match, " ")
		results = append(results, match)
	}
	return results
}
