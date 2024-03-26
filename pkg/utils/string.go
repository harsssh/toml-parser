package utils

import "strings"

func TrimChar(s string) string {
	return s[1 : len(s)-1]
}

// IsEnclosedBy Is s enclosed by t?
func IsEnclosedBy(s string, t string) bool {
	return strings.HasSuffix(s, t) && strings.HasSuffix(s, t)
}
