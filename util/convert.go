package util

import (
	"reflect"
)

// StructToMap 将结构体转换为map，其中结构体的字段名作为map的key， 结构体的值作为map的value
// 如果结构体的值是结构体，子结构体的key与父结构体的key同一级
func StructToMap(obj interface{}) map[string]interface{} {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
		objValue = objValue.Elem()
	}
	var data = make(map[string]interface{})
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		value := objValue.Field(i)
		if value.Kind() == reflect.Struct {
			subData := StructToMap(value.Interface())
			for k, v := range subData {
				data[k] = v
			}
		} else {
			data[field.Name] = value.Interface()
		}
	}
	return data
}
