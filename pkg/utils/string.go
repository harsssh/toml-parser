package utils

import "strings"

func TrimChar(s string) string {
	return TrimNChar(s, 1)
}

func TrimNChar(s string, n int) string {
	return s[n : len(s)-n]
}

// IsEnclosedBy Is s enclosed by t?
func IsEnclosedBy(s string, t string) bool {
	return IsBracketedBy(s, t, t)
}

func IsBracketedBy(s, prefix, suffix string) bool {
	return strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix)
}
