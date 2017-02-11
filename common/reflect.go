package common

import (
	"github.com/fatih/structs"
	"reflect"
)

func DeepEqual(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

func TypeOf(i interface{}) string {
	return reflect.TypeOf(i).String()
}

func StructToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

func FillStruct(s interface{}, data map[string]interface{}) (err error) {
	for k, v := range data {
		if err = SetField(s, k, v); err != nil {
			return err
		}
	}
	return nil
}

func SetField(s interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(s).Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return Errorf("No such field: %s in struct", name)
	}
	if !structFieldValue.CanSet() {
		return Errorf("Cannot set '%s' field", name)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return Errorf("Value type: %v doesn't match field type: %v", val.Type(), structFieldType)
	}
	structFieldValue.Set(val)
	return nil
}

/*
func StructToMap(s interface{}, tag string) map[string]interface{} {
	mp := make(map[string]interface{})
	val := reflect.ValueOf(s).Elem()
	_type := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := _type.Field(i)
		if t := field.Tag.Get(tag); t != "" {
			mp[t] = val.Field(i).Interface()
		}
	}
	return mp
}
*/
