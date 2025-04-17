package utils

import (
	"strconv"
)

func ConvertGiga[T int | string](size T) string {
	switch v := any(size).(type) {
	case int:
		return strconv.Itoa(v) + "Gi"
	case string:
		return v + "Gi"
	default:
		return ""
	}
}
