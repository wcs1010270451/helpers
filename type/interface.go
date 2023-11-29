package _type

import (
	"errors"
	"fmt"
	"reflect"
)

// Empty 判空函数
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// ToString 转字符串
func ToString(obj interface{}) string {
	if str, ok := obj.(string); ok {
		return str
	}
	return ""
}

// 转浮点数
func ToFloat(obj interface{}) (float64, error) {
	switch v := obj.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, errors.New(fmt.Sprintf("Expected a float value but got %T", v))
	}
}

// 转int
func ToInt(obj interface{}) (int, error) {
	switch v := obj.(type) {
	case int:
		return v, nil
	default:
		return 0, errors.New(fmt.Sprintf("Expected an integer value but got %T", v))
	}
}

// 转map
func ToMap(obj interface{}) (map[string]interface{}, error) {
	if reflect.TypeOf(obj).Kind() != reflect.Map {
		return nil, errors.New("obj is not a type of map")
	}

	//执行类型断言以确保输入确实是一个映射
	result := make(map[string]interface{})
	valueOf := reflect.ValueOf(obj)
	for _, key := range valueOf.MapKeys() {
		result[key.String()] = valueOf.MapIndex(key).Interface()
	}

	return result, nil
}

// 转布尔值
func ToBool(obj interface{}) (bool, error) {
	switch v := obj.(type) {
	case bool:
		return v, nil
	default:
		return false, errors.New(fmt.Sprintf("Expected a boolean value but got  %T", v))
	}
}

// 转数组
func ToSlice(obj interface{}) ([]interface{}, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Slice {
		var result []interface{}
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			result = append(result, elem.Interface())
		}
		return result, nil
	} else {
		return nil, errors.New("Not a slice")
	}
}
