package util

import (
	"strconv"
	"strings"
)

func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

func SplitN(s, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

func ParseUint16(s string) (uint16, error) {
	x, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(x), nil
}

func ParseUint32(s string) (uint32, error) {
	x, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(x), nil
}

func ParseUint64(s string) (uint64, error) {
	x, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return x, nil
}

func MustParseUint16(s string) uint16 {
	x, err := ParseUint16(s)
	Check(err)
	return x
}

func MustParseUint32(s string) uint32 {
	x, err := ParseUint32(s)
	Check(err)
	return x
}

func MustParseUint64(s string) uint64 {
	x, err := ParseUint64(s)
	Check(err)
	return x
}
