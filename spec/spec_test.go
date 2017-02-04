package spec

import (
	. "github.com/zbo14/envoke/common"
	"testing"
)

var data = Data{
	"hello world": Data{"/": "ip4/127.0.0.1/udp/1234"},
	"1":           Data{"2": Data{"/": "ip4/10.20.30.40/tcp/443"}},
	"3":           Data{"4": Data{"5": Data{"6": Data{"7": Data{"8": Data{"/": "ip4/127.0.0.1/udp/5678"}}}}}},
	"person":      Data{"name": "zach"},
}

func TestImpl(t *testing.T) {
	copy := CopyData(data)
	if !reflect.DeepEqual(copy, data) {
		t.Fatal("Expected data models to be the same")
	}
	TransformJSON(copy)
	TransformIPLD(copy)
	if !DeepEqual(data, copy) {
		t.Fatal("Expected data models to be the same")
	}
}

func BenchmarkTransform(b *testing.B) {
	for i := 0; i < b.N; i++ {
		copy := CopyData(data)
		TransformJSON(copy)
		TransformIPLD(copy)
	}
}

func BenchmarkIterTransform(b *testing.B) {
	for i := 0; i < b.N; i++ {
		copy := CopyData(data)
		IterTransformJSON(copy)
		IterTransformIPLD(copy)
	}
}
