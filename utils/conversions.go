package utils

import (
	"strings"
)

func ConvertLocalizedString(s string) string {
	if strings.ContainsAny(s, ",") {
		return strings.Replace(s, ",", ".", 1)
	}
	return s
}
