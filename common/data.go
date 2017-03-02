package common

import (
	"time"
)

type Data map[string]interface{}

func (d Data) Get(key string) interface{}        { return d[key] }
func (d Data) Set(key string, value interface{}) { d[key] = value }
func (d Data) Clear(key string)                  { d[key] = nil }

func (d Data) GetBool(key string) bool                    { return MustAssertBool(d.Get(key)) }
func (d Data) GetData(key string) Data                    { return AssertData(d.Get(key)) }
func (d Data) GetDataSlice(key string) []Data             { return AssertDataSlice(d.Get(key)) }
func (d Data) GetFloat64(key string) float64              { return AssertFloat64(d.Get(key)) }
func (d Data) GetInt(key string) int                      { return AssertInt(d.Get(key)) }
func (d Data) GetInt32(key string) int32                  { return AssertInt32(d.Get(key)) }
func (d Data) GetInt32Slice(key string) []int32           { return AssertInt32Slice(d.Get(key)) }
func (d Data) GetInt64(key string) int64                  { return AssertInt64(d.Get(key)) }
func (d Data) GetInterfaceSlice(key string) []interface{} { return AssertInterfaceSlice(d.Get(key)) }
func (d Data) GetMap(key string) map[string]interface{}   { return AssertMap(d.Get(key)) }
func (d Data) GetMapData(key string) Data                 { return AssertMapData(d.Get(key)) }
func (d Data) GetStr(key string) string                   { return AssertStr(d.Get(key)) }

func (d Data) GetStrInt(key string) int {
	x, err := Atoi(d.GetStr(key))
	if err != nil {
		return 0
	}
	return x
}

func (d Data) GetStrSlice(key string) []string { return AssertStrSlice(d.Get(key)) }
func (d Data) GetTime(key string) time.Time    { return AssertTime(d.Get(key)) }

func (d Data) GetInnerValue(keys ...string) (v interface{}) {
	inner := d
	for i, k := range keys {
		if i == len(keys)-1 {
			v = inner.Get(k)
			break
		}
		inner = inner.GetMap(k)
	}
	return
}

func (d Data) SetInnerValue(v interface{}, keys ...string) {
	inner := d
	for i, k := range keys {
		if i == len(keys)-1 {
			inner.Set(k, v)
			return
		}
		inner = inner.GetMap(k)
	}
}

func (d Data) GetInnerStr(keys ...string) string {
	return AssertStr(d.GetInnerValue(keys...))
}

func (d Data) GetInnerData(keys ...string) Data {
	return AssertMapData(d.GetInnerValue(keys...))
}
