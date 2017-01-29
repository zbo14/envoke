package util

import (
	"strconv"
	"strings"
)

func ToLower(s string) string {
	return strings.ToLower(s)
}

func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

func SplitN(s, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

func ParseUint16(s string) (int, error) {
	x, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return int(x), nil
}

func ParseUint32(s string) (int, error) {
	x, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(x), nil
}

func ParseUint64(s string) (int, error) {
	x, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(x), nil
}

func MustParseUint16(s string) int {
	x, err := ParseUint16(s)
	Check(err)
	return x
}

func MustParseUint32(s string) int {
	x, err := ParseUint32(s)
	Check(err)
	return x
}

func MustParseUint64(s string) int {
	x, err := ParseUint64(s)
	Check(err)
	return x
}
