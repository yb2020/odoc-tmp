package utils

import (
	"fmt"
	"reflect"
	"time"
)

// ConvertUint64SliceToInt64Slice 将 uint64 切片转换为 int64 切片
func ConvertUint64SliceToInt64Slice(s []uint64) []int64 {
	ids := make([]int64, len(s))
	for i, id := range s {
		ids[i] = int64(id)
	}
	return ids
}

// ConvertInt64SliceToUint64Slice 将 int64 切片转换为 uint64 切片
func ConvertInt64SliceToUint64Slice(s []int64) []uint64 {
	ids := make([]uint64, len(s))
	for i, id := range s {
		ids[i] = uint64(id)
	}
	return ids
}

// GetInt64FromUint64Ptr 安全地从 uint64 指针获取 int64 值，如果指针为 nil 则返回默认值
func GetInt64FromUint64Ptr(ptr *uint64, defaultValue int64) int64 {
	if ptr == nil {
		return defaultValue
	}
	return int64(*ptr)
}

// GetUint64FromInt64Ptr 安全地从 int64 指针获取 uint64 值，如果指针为 nil 则返回默认值
func GetUint64FromInt64Ptr(ptr *int64, defaultValue uint64) uint64 {
	if ptr == nil {
		return defaultValue
	}
	return uint64(*ptr)
}

// GetUint64PtrFromInt64Ptr 安全地从 int64 指针获取 uint64 指针，如果指针为 nil 则返回默认值的指针
func GetUint64PtrFromInt64Ptr(ptr *int64, defaultValue uint64) *uint64 {
	if ptr == nil {
		return &defaultValue
	}
	value := uint64(*ptr)
	return &value
}

// GetIntFromIntPtr 安全地从 int 指针获取 int 值，如果指针为 nil 则返回默认值
func GetIntFromIntPtr(ptr *int, defaultValue int) int {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetInt64FromInt64Ptr 安全地从 int64 指针获取 int64 值，如果指针为 nil 则返回默认值
func GetInt64FromInt64Ptr(ptr *int64, defaultValue int64) int64 {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetInt32FromInt32Ptr 安全地从 int32 指针获取 int32 值，如果指针为 nil 则返回默认值
func GetInt32FromInt32Ptr(ptr *int32, defaultValue int32) int32 {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetStringPtrValue 安全地获取字符串指针的值，如果指针为 nil 则返回默认值
func GetStringPtrValue(ptr *string, defaultValue string) string {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetInt64PtrValue 安全地获取 int64 指针的值，如果指针为 nil 则返回默认值
func GetInt64PtrValue(ptr *int64, defaultValue int64) int64 {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetIntPtrValue 安全地获取 int 指针的值，如果指针为 nil 则返回默认值
func GetIntPtrValue(ptr *int, defaultValue int) int {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetUint64PtrValue 安全地获取 uint64 指针的值，如果指针为 nil 则返回默认值
func GetUint64PtrValue(ptr *uint64, defaultValue uint64) uint64 {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetBoolPtrValue 安全地获取布尔指针的值，如果指针为 nil 则返回默认值
func GetBoolPtrValue(ptr *bool, defaultValue bool) bool {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetEnumStringValue 安全地获取枚举类型的字符串表示，如果指针为 nil 则返回默认值
func GetEnumStringValue(ptr interface{}, defaultValue string) string {
	if ptr == nil {
		return defaultValue
	}

	// 使用类型断言获取 String() 方法
	if stringer, ok := ptr.(fmt.Stringer); ok {
		return stringer.String()
	}

	return defaultValue
}

// GetEnumFromPtr 安全地从枚举类型指针获取枚举值，如果指针为 nil 则返回默认值
// T 可以是任何枚举类型
func GetEnumFromPtr[T comparable](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetStringFromPtr 安全地从任意类型指针获取字符串表示，如果指针为 nil 则返回默认值
func GetStringFromPtr(ptr interface{}, defaultValue string) string {
	if ptr == nil {
		return defaultValue
	}

	// 使用反射获取指针指向的值
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return defaultValue
	}

	// 获取指针指向的值
	val := v.Elem().Interface()

	// 如果值实现了 fmt.Stringer 接口，调用 String() 方法
	if stringer, ok := val.(fmt.Stringer); ok {
		return stringer.String()
	}

	// 否则使用 fmt.Sprint 转换为字符串
	return fmt.Sprint(val)
}

// GetDurationFromPtr 安全地从时间指针获取时间值，如果指针为 nil 则返回默认值
func GetDurationFromPtr(ptr *time.Duration, defaultValue time.Duration) time.Duration {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GetMillisecondDuration 将整数值转换为毫秒单位的 time.Duration
// 这在从配置文件读取时间值时特别有用，因为配置文件通常以毫秒为单位
func GetMillisecondDuration(milliseconds int) time.Duration {
	return time.Duration(milliseconds) * time.Millisecond
}
