package common

import "time"

func MustAssertBool(v interface{}) bool {
	return v.(bool)
}

func AssertData(v interface{}) Data {
	if d, ok := v.(Data); ok {
		return d
	}
	if m, ok := v.(map[string]interface{}); ok {
		return Data(m)
	}
	return nil
}

func AssertDataSlice(v interface{}) []Data {
	if slice, ok := v.([]Data); ok {
		return slice
	}
	Printf("%T", v)
	return nil
}

func AssertInt(v interface{}) int {
	if n, ok := v.(int); ok {
		return n
	}
	return 0
}

func AssertInt32(v interface{}) int32 {
	if n, ok := v.(int32); ok {
		return n
	}
	return 0
}

func AssertInt32Slice(v interface{}) []int32 {
	if slice, ok := v.([]int32); ok {
		return slice
	}
	return nil
}

func AssertInt64(v interface{}) int64 {
	if n, ok := v.(int64); ok {
		return n
	}
	return 0
}

func AssertStr(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func AssertStrSlice(v interface{}) []string {
	if slice, ok := v.([]string); ok {
		return slice
	}
	return nil
}

func AssertTime(v interface{}) time.Time {
	if time, ok := v.(time.Time); ok {
		return time
	}
	return time.Time{}
}
