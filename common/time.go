package common

import (
	"encoding/binary"
	"strconv"
	"time"
)

func Time() time.Time {
	return time.Now()
}

func Timestr() string {
	return Time().String()
}

func Datestr() string {
	return ToTheDay(Timestr())
}

func Timestamp() int64 {
	return Time().Unix()
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

func ParseTimestr(timestr string) time.Time {
	yr, _ := strconv.Atoi(timestr[:4])
	mo, _ := strconv.Atoi(timestr[5:7])
	d, _ := strconv.Atoi(timestr[8:10])
	hr, _ := strconv.Atoi(timestr[11:13])
	min, _ := strconv.Atoi(timestr[14:16])
	sec, _ := strconv.Atoi(timestr[17:19])
	return time.Date(yr, time.Month(mo), d, hr, min, sec, 0, time.Local)
}

func ParseDatestr(datestr string) time.Time {
	yr, _ := strconv.Atoi(datestr[:4])
	mo, _ := strconv.Atoi(datestr[5:7])
	d, _ := strconv.Atoi(datestr[8:10])
	return time.Date(yr, time.Month(mo), d, 0, 0, 0, 0, time.Local)
}

func ToTheDay(timestr string) string {
	return timestr[:10]
}
