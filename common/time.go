package common

import (
	"encoding/binary"
	"strconv"
	"time"
)

func SleepMilli(d time.Duration) {
	time.Sleep(d * time.Millisecond)
}

func SleepSeconds(d time.Duration) {
	time.Sleep(d * time.Second)
}

func Now() time.Time {
	return time.Now()
}

func NowStr() string {
	return Now().String()
}

func Time(yr int, mo time.Month, day, hr, min, sec int) time.Time {
	return time.Date(yr, mo, day, hr, min, sec, 0, time.Local)
}

func Date(yr int, mo time.Month, day int) time.Time {
	return Time(yr, mo, day, 0, 0, 0)
}

func TimeStr(yr int, mo time.Month, day, hr, min, sec int) string {
	time := Time(yr, mo, day, hr, min, sec).String()
	return time[:10]
}

func DateStr(yr int, mo time.Month, day int) string {
	date := Date(yr, mo, day).String()
	return date[:19]
}

func Timestamp() int64 {
	return Now().Unix()
}

func TimestampBytes(x int64) []byte {
	p := make([]byte, 10)
	n := binary.PutVarint(p, x)
	return p[:n]
}

func TimestampFromBytes(p []byte) int64 {
	x, _ := binary.Varint(p)
	return x
}

func ParseTimeStr(timestr string) (time.Time, error) {
	if len(timestr) < 19 {
		return time.Time{}, ErrInvalidSize
	}
	yr, err := strconv.Atoi(timestr[:4])
	if err != nil {
		return time.Time{}, err
	}
	mo, err := strconv.Atoi(timestr[5:7])
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(timestr[8:10])
	if err != nil {
		return time.Time{}, err
	}
	hr, err := strconv.Atoi(timestr[11:13])
	if err != nil {
		return time.Time{}, err
	}
	min, err := strconv.Atoi(timestr[14:16])
	if err != nil {
		return time.Time{}, err
	}
	sec, err := strconv.Atoi(timestr[17:19])
	if err != nil {
		return time.Time{}, err
	}
	return Time(yr, time.Month(mo), day, hr, min, sec), nil
}

func ParseDateStr(datestr string) (time.Time, error) {
	if len(datestr) < 10 {
		return time.Time{}, ErrInvalidSize
	}
	yr, err := strconv.Atoi(datestr[:4])
	if err != nil {
		return time.Time{}, err
	}
	mo, err := strconv.Atoi(datestr[5:7])
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(datestr[8:10])
	if err != nil {
		return time.Time{}, err
	}
	return Date(yr, time.Month(mo), day), nil
}

func MustParseTimeStr(timestr string) time.Time {
	time, err := ParseTimeStr(timestr)
	Check(err)
	return time
}

func MustParseDateStr(datestr string) time.Time {
	date, err := ParseDateStr(datestr)
	Check(err)
	return date
}
