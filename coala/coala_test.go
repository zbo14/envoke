package coala

import (
	"testing"
)

func TestImpl(t *testing.T) {
	/*
		ipld := Data{
			"hello world": Data{"/": "ip4/127.0.0.1/udp/1234"},
			"1":           Data{"2": Data{"/": "ip4/10.20.30.40/tcp/443"}},
			"name":        "zach",
		}
		t.Log(ipld)
		json := TransformJSON(ipld)
		t.Log(json)
	*/
	data := Data{
		"hello world": Data{"/": "ip4/127.0.0.1/udp/1234"},
		"1":           Data{"2": Data{"/": "ip4/10.20.30.40/tcp/443"}},
		"name":        "zach",
	}
	t.Log(data)
	IterTransformJSON(data)
	t.Log(data)
}
