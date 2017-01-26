package util

import (
	re "regexp"
)

func MatchString(pattern, s string) bool {
	match, err := re.MatchString(pattern, s)
	Check(err)
	return match
}
