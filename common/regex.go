package common

import (
	re "regexp"
)

func MatchString(pattern, s string) bool {
	match, err := re.MatchString(pattern, s)
	Check(err)
	return match
}

func MatchBytes(pattern string, p []byte) bool {
	match, err := re.Match(pattern, p)
	Check(err)
	return match
}
