package utils

import "strings"

func MapToString(m map[string]string) string {
	var sb strings.Builder
	first := true
	for key, value := range m {
		if !first {
			sb.WriteString(",")
		}
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(value)
		first = false
	}
	return sb.String()
}
