package common

import re "regexp"

func Match(pattern string, p []byte) bool {
	match, err := re.Match(pattern, p)
	Check(err)
	return match
}

func Submatch(pattern string, p []byte) [][]byte {
	regex := re.MustCompile(pattern)
	return regex.FindSubmatch(p)
}

func MatchStr(pattern, s string) bool {
	match, err := re.MatchString(pattern, s)
	Check(err)
	return match
}

func SubmatchStr(pattern, s string) []string {
	regex := re.MustCompile(pattern)
	return regex.FindStringSubmatch(s)
}
