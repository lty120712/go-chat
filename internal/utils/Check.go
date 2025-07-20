package utils

import "reflect"

// IsZeroValue 检查传入的对象是否为零值
func IsZero(obj interface{}) bool {
	// 获取对象的反射值
	val := reflect.ValueOf(obj)
	// 如果对象是 nil，直接返回 true
	if val.IsValid() == false {
		return true
	}
	// 判断对象是否为零值
	return val.IsZero()
}
