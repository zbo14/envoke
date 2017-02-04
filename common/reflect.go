package util

import (
	"reflect"
)

func DeepEqual(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}
