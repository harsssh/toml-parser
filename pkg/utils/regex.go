package utils

import "regexp"

func MatchesPattern(s, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(s)
}
